/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/git-qfzhang/hello-golang/grpc-go/helloworld/helloworld"
	"github.com/git-qfzhang/hello-golang/grpc-go/helloworld/utils"

	"github.com/cenkalti/backoff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	err = utils.RetryWithCondition(ctx, backoff.NewConstantBackOff(time.Second), func() (bool, error) {
		reply, serverErr := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
		if serverErr == nil {
			log.Printf("Greeting: %s", reply.GetMessage())
			return false, nil
		}
		switch grpc.Code(serverErr) {
		case codes.Unavailable:
			log.Printf("server error: %v, retry", serverErr)
			return true, serverErr
		case codes.Canceled:
			log.Printf("server error: %v, retry", serverErr)
			return true, serverErr
		case codes.Internal:
			log.Printf("server error: %v, no need to retry", serverErr)
			return false, serverErr
		default:
			log.Printf("server error: %v, no need to retry", serverErr)
			return false, serverErr
		}
		log.Printf("Greeting: %s", reply.GetMessage())
		return false, nil
	})
	if err != nil {
		log.Printf("retry done, err: %v", err)
	}
}
