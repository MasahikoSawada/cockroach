// Code generated by protoc-gen-grpc-gateway
// source: cockroach/server/admin.proto
// DO NOT EDIT!

/*
Package server is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gengo/grpc-gateway/runtime"
	"github.com/gengo/grpc-gateway/utilities"
	"github.com/gogo/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
)

var _ codes.Code
var _ io.Reader
var _ = runtime.String
var _ = json.Marshal
var _ = utilities.NewDoubleArray

func request_Admin_Users_0(ctx context.Context, client AdminClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq UsersRequest
	var metadata runtime.ServerMetadata

	msg, err := client.Users(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func request_Admin_Databases_0(ctx context.Context, client AdminClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq DatabasesRequest
	var metadata runtime.ServerMetadata

	msg, err := client.Databases(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func request_Admin_DatabaseDetails_0(ctx context.Context, client AdminClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq DatabaseDetailsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["database"]
	if !ok {
		return nil, metadata, grpc.Errorf(codes.InvalidArgument, "missing parameter %s", "database")
	}

	protoReq.Database, err = runtime.String(val)

	if err != nil {
		return nil, metadata, err
	}

	msg, err := client.DatabaseDetails(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func request_Admin_TableDetails_0(ctx context.Context, client AdminClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq TableDetailsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["database"]
	if !ok {
		return nil, metadata, grpc.Errorf(codes.InvalidArgument, "missing parameter %s", "database")
	}

	protoReq.Database, err = runtime.String(val)

	if err != nil {
		return nil, metadata, err
	}

	val, ok = pathParams["table"]
	if !ok {
		return nil, metadata, grpc.Errorf(codes.InvalidArgument, "missing parameter %s", "table")
	}

	protoReq.Table, err = runtime.String(val)

	if err != nil {
		return nil, metadata, err
	}

	msg, err := client.TableDetails(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

var (
	filter_Admin_Events_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_Admin_Events_0(ctx context.Context, client AdminClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq EventsRequest
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_Admin_Events_0); err != nil {
		return nil, metadata, grpc.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.Events(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func request_Admin_SetUIData_0(ctx context.Context, client AdminClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq SetUIDataRequest
	var metadata runtime.ServerMetadata

	if err := json.NewDecoder(req.Body).Decode(&protoReq); err != nil {
		return nil, metadata, grpc.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.SetUIData(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

var (
	filter_Admin_GetUIData_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_Admin_GetUIData_0(ctx context.Context, client AdminClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetUIDataRequest
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_Admin_GetUIData_0); err != nil {
		return nil, metadata, grpc.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.GetUIData(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// RegisterAdminHandlerFromEndpoint is same as RegisterAdminHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterAdminHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterAdminHandler(ctx, mux, conn)
}

// RegisterAdminHandler registers the http handlers for service Admin to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterAdminHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	client := NewAdminClient(conn)

	mux.Handle("GET", pattern_Admin_Users_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		resp, md, err := request_Admin_Users_0(runtime.AnnotateContext(ctx, req), client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, w, req, err)
			return
		}

		forward_Admin_Users_0(ctx, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Admin_Databases_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		resp, md, err := request_Admin_Databases_0(runtime.AnnotateContext(ctx, req), client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, w, req, err)
			return
		}

		forward_Admin_Databases_0(ctx, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Admin_DatabaseDetails_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		resp, md, err := request_Admin_DatabaseDetails_0(runtime.AnnotateContext(ctx, req), client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, w, req, err)
			return
		}

		forward_Admin_DatabaseDetails_0(ctx, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Admin_TableDetails_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		resp, md, err := request_Admin_TableDetails_0(runtime.AnnotateContext(ctx, req), client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, w, req, err)
			return
		}

		forward_Admin_TableDetails_0(ctx, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Admin_Events_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		resp, md, err := request_Admin_Events_0(runtime.AnnotateContext(ctx, req), client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, w, req, err)
			return
		}

		forward_Admin_Events_0(ctx, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_Admin_SetUIData_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		resp, md, err := request_Admin_SetUIData_0(runtime.AnnotateContext(ctx, req), client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, w, req, err)
			return
		}

		forward_Admin_SetUIData_0(ctx, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Admin_GetUIData_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		resp, md, err := request_Admin_GetUIData_0(runtime.AnnotateContext(ctx, req), client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, w, req, err)
			return
		}

		forward_Admin_GetUIData_0(ctx, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_Admin_Users_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"_admin", "v1", "users"}, ""))

	pattern_Admin_Databases_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"_admin", "v1", "databases"}, ""))

	pattern_Admin_DatabaseDetails_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 1, 0, 4, 1, 5, 3}, []string{"_admin", "v1", "databases", "database"}, ""))

	pattern_Admin_TableDetails_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 1, 0, 4, 1, 5, 3, 2, 4, 1, 0, 4, 1, 5, 5}, []string{"_admin", "v1", "databases", "database", "tables", "table"}, ""))

	pattern_Admin_Events_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"_admin", "v1", "events"}, ""))

	pattern_Admin_SetUIData_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"_admin", "v1", "uidata"}, ""))

	pattern_Admin_GetUIData_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"_admin", "v1", "uidata"}, ""))
)

var (
	forward_Admin_Users_0 = runtime.ForwardResponseMessage

	forward_Admin_Databases_0 = runtime.ForwardResponseMessage

	forward_Admin_DatabaseDetails_0 = runtime.ForwardResponseMessage

	forward_Admin_TableDetails_0 = runtime.ForwardResponseMessage

	forward_Admin_Events_0 = runtime.ForwardResponseMessage

	forward_Admin_SetUIData_0 = runtime.ForwardResponseMessage

	forward_Admin_GetUIData_0 = runtime.ForwardResponseMessage
)
