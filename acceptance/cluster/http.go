// Copyright 2015 The Cockroach Authors.
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
// Author: Peter Mattis (peter@cockroachlabs.com)

package cluster

import (
	"crypto/tls"
	"net/http"

	"github.com/cockroachdb/cockroach/base"
	"github.com/cockroachdb/cockroach/util"
)

// HTTPClient is an http.Client configured for querying a cluster. We need to
// run with "InsecureSkipVerify" (at least on Docker) due to the fact that we
// cannot use a fixed hostname to reach the cluster. This in turn means that we
// do not have a verified server name in the certs.
var HTTPClient = http.Client{
	Timeout: base.NetworkTimeout,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}

// getJSON is a convenience wrapper around cockroach/util.GetJSON(), which retrieves
// an URL specified by the parameters and unmarshals the result into the supplied
// interface.
func getJSON(tls bool, hostport, path string, v interface{}) error {
	scheme := "https"
	if !tls {
		scheme = "http"
	}

	return util.GetJSON(&HTTPClient, scheme, hostport, path, v)
}
