// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import "C"

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/openconfig/lemming/dataplane/forwarding"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

//export initialize
func initialize(port int) {
	fwdSrv := forwarding.New("engine")

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	fwdpb.RegisterForwardingServer(srv, fwdSrv)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve forwarding server: %v", err)
	}
}

func main() {
}