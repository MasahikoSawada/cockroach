// Copyright 2014 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.
//
// Author: Spencer Kimball (spencer.kimball@gmail.com)

package kv

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/cockroachdb/cockroach/roachpb"
	"github.com/cockroachdb/cockroach/rpc"
	"github.com/cockroachdb/cockroach/util"
	"github.com/cockroachdb/cockroach/util/log"
	"github.com/cockroachdb/cockroach/util/retry"
)

// orderingPolicy is an enum for ordering strategies when there
// are multiple endpoints available.
type orderingPolicy int

const (
	// orderStable uses endpoints in the order provided.
	orderStable = iota
	// orderRandom randomly orders available endpoints.
	orderRandom
)

// A SendOptions structure describes the algorithm for sending RPCs to one or
// more replicas, depending on error conditions and how many successful
// responses are required.
type SendOptions struct {
	// Ordering indicates how the available endpoints are ordered when
	// deciding which to send to (if there are more than one).
	Ordering orderingPolicy
	// SendNextTimeout is the duration after which RPCs are sent to
	// other replicas in a set.
	SendNextTimeout time.Duration
	// Timeout is the maximum duration of an RPC before failure.
	// 0 for no timeout.
	Timeout time.Duration
	// Information about the request is added to this trace. Must not be nil.
	Trace opentracing.Span
}

// An rpcError indicates a failure to send the RPC. rpcErrors are
// retryable.
type rpcError struct {
	error
}

func newRPCError(err error) rpcError {
	return rpcError{err}
}

// CanRetry implements the Retryable interface.
// TODO(tschottdorf): the way this is used by rpc/send suggests that it
// may be better if these weren't retriable - they are returned when the
// connection fails, i.e. for example when a node is down or the network
// fails. Retrying on such errors keeps the caller waiting for a long time
// and without a positive outlook.
func (r rpcError) CanRetry() bool { return true }

type batchClient struct {
	remoteAddr string
	conn       *grpc.ClientConn
	client     roachpb.InternalClient
	args       roachpb.BatchRequest
}

func shuffleClients(clients []batchClient) {
	for i, n := 0, len(clients); i < n-1; i++ {
		j := rand.Intn(n-i) + i
		clients[i], clients[j] = clients[j], clients[i]
	}
}

type batchCall struct {
	reply *roachpb.BatchResponse
	err   error
}

// Send sends one or more RPCs to clients specified by the slice of
// replicas. On success, Send returns the first successful reply. Otherwise,
// Send returns an error if and as soon as the number of failed RPCs exceeds
// the available endpoints less the number of required replies.
func send(opts SendOptions, replicas ReplicaSlice,
	args roachpb.BatchRequest, rpcContext *rpc.Context) (*roachpb.BatchResponse, error) {
	sp := opts.Trace // must not be nil

	if len(replicas) < 1 {
		return nil, roachpb.NewSendError(
			fmt.Sprintf("insufficient replicas (%d) to satisfy send request of %d",
				len(replicas), 1), false)
	}

	done := make(chan batchCall, len(replicas))

	clients := make([]batchClient, 0, len(replicas))
	for _, replica := range replicas {
		conn, err := rpcContext.GRPCDial(replica.NodeDesc.Address.String())
		if err != nil {
			return nil, err
		}
		argsCopy := args
		argsCopy.Replica = replica.ReplicaDescriptor
		clients = append(clients, batchClient{
			remoteAddr: replica.NodeDesc.Address.String(),
			conn:       conn,
			client:     roachpb.NewInternalClient(conn),
			args:       argsCopy,
		})
	}

	var orderedClients []batchClient
	switch opts.Ordering {
	case orderStable:
		orderedClients = clients
	case orderRandom:
		// Randomly permute order, but keep known-unhealthy clients last.
		var nHealthy int
		for i, client := range clients {
			clientState, err := client.conn.State()
			if err != nil {
				return nil, err
			}
			if clientState == grpc.Ready {
				clients[i], clients[nHealthy] = clients[nHealthy], clients[i]
				nHealthy++
			}
		}

		shuffleClients(clients[:nHealthy])
		shuffleClients(clients[nHealthy:])

		orderedClients = clients
	}
	// TODO(spencer): going to need to also sort by affinity; closest
	// ping time should win. Makes sense to have the rpc client/server
	// heartbeat measure ping times. With a bit of seasoning, each
	// node will be able to order the healthy replicas based on latency.

	// Send the first request.
	sendOneFn(orderedClients[0], opts.Timeout, rpcContext, sp, done)
	orderedClients = orderedClients[1:]

	var errors, retryableErrors int

	// Wait for completions.
	var sendNextTimer util.Timer
	defer sendNextTimer.Stop()
	for {
		sendNextTimer.Reset(opts.SendNextTimeout)
		select {
		case <-sendNextTimer.C:
			sendNextTimer.Read = true
			// On successive RPC timeouts, send to additional replicas if available.
			if len(orderedClients) > 0 {
				sp.LogEvent("timeout, trying next peer")
				sendOneFn(orderedClients[0], opts.Timeout, rpcContext, sp, done)
				orderedClients = orderedClients[1:]
			}

		case call := <-done:
			err := call.err
			if err == nil {
				if log.V(2) {
					log.Infof("successful reply: %+v", call.reply)
				}

				return call.reply, nil
			}

			// Error handling.
			if log.V(1) {
				log.Warningf("error reply: %s", err)
			}

			errors++

			// Since we have a reconnecting client here, disconnect errors are retryable.
			disconnected := err == io.ErrUnexpectedEOF
			if retryErr, ok := err.(retry.Retryable); disconnected || (ok && retryErr.CanRetry()) {
				retryableErrors++
			}

			if remainingNonErrorRPCs := len(replicas) - errors; remainingNonErrorRPCs < 1 {
				return nil, roachpb.NewSendError(
					fmt.Sprintf("too many errors encountered (%d of %d total): %v",
						errors, len(clients), err), remainingNonErrorRPCs+retryableErrors >= 1)
			}
			// Send to additional replicas if available.
			if len(orderedClients) > 0 {
				sp.LogEvent("error, trying next peer")
				sendOneFn(orderedClients[0], opts.Timeout, rpcContext, sp, done)
				orderedClients = orderedClients[1:]
			}
		}
	}
}

// Allow local calls to be dispatched directly to the local server without
// sending an RPC.
var enableLocalCalls = os.Getenv("ENABLE_LOCAL_CALLS") != "0"

// sendOneFn is overwritten in tests to mock sendOne.
var sendOneFn = sendOne

// sendOne invokes the specified RPC on the supplied client when the
// client is ready. On success, the reply is sent on the channel;
// otherwise an error is sent.
//
// Do not call directly, but instead use sendOneFn. Tests mock out this method
// via sendOneFn in order to test various error cases.
func sendOne(client batchClient, timeout time.Duration,
	rpcContext *rpc.Context, trace opentracing.Span, done chan batchCall) {
	addr := client.remoteAddr
	if log.V(2) {
		log.Infof("sending request to %s: %+v", addr, client.args)
	}
	trace.LogEvent(fmt.Sprintf("sending to %s", addr))

	// TODO(tamird/tschottdorf): pass this in from DistSender.
	ctx := context.TODO()
	if timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, timeout)
	}

	if localServer := rpcContext.LocalInternalServer; enableLocalCalls && localServer != nil && addr == rpcContext.LocalAddr {
		reply, err := localServer.Batch(ctx, &client.args)
		done <- batchCall{reply: reply, err: err}
		return
	}

	go func() {
		c := client.conn
		for state, err := c.State(); state != grpc.Ready; state, err = c.WaitForStateChange(ctx, state) {
			if err != nil {
				done <- batchCall{err: newRPCError(
					util.Errorf("rpc to %s failed: %s", addr, err))}
				return
			}
			if state == grpc.Shutdown {
				done <- batchCall{err: newRPCError(
					util.Errorf("rpc to %s failed as client connection was closed", addr))}
				return
			}
		}

		reply, err := client.client.Batch(ctx, &client.args)
		done <- batchCall{reply: reply, err: err}
	}()
}
