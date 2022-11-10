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

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"

	"log"
	"net"
	"sync"

	pb "github.com/git-qfzhang/hello-golang/grpc-go/helloworld/helloworld"
	"github.com/git-qfzhang/hello-golang/grpc-go/helloworld/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	// Get metadata which dose not contain grpc-timeout from request ctx
	serverCtx, cancel := context.WithCancel(context.Background())
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		serverCtx, cancel = context.WithCancel(metadata.NewIncomingContext(context.Background(), md))
	}
	var reqWG sync.WaitGroup
	reqWG.Add(1)
	go func() {
		defer reqWG.Done()
		select {
		case <-ctx.Done():
			log.Printf("context from client is done. Err: %v", ctx.Err())
			cancel()
		case <-serverCtx.Done():
			log.Printf("context in server is done. Err: %v", serverCtx.Err())
		}
	}()
	err := utils.HandleRequest(serverCtx)
	cancel()
	reqWG.Wait()
	if err != nil {
		log.Printf("handle request done, err: %v", err)
		// 建议：处理请求出错后，应该先判断 ctx 是否超时，如果是，则应该保证此时返回的 error code 为 DeadlineExceeded
		// 如果 ctx 未超时，则可以基于请求处理错误信息返回对应的 error code
		if err.Error() == context.Canceled.Error() {
			return &pb.HelloReply{}, grpc.Errorf(codes.Canceled, err.Error())
		}
		return &pb.HelloReply{}, err
	}

	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
