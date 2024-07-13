package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "datastore/protos/keyval"

	"google.golang.org/grpc"
)

// Server for the keyvalue service
type keyValServer struct {
	pb.UnimplementedKeyValueServer

	// map from string key to string value
	sets map[string]string

	// mutex for sets
	setsMutex sync.RWMutex
}

// Get value using key
// returns a pb.Value type
func (s *keyValServer) GetValue(ctx context.Context, key *pb.Key) (*pb.Value, error) {
	s.setsMutex.RLock()
	value, found := s.sets[key.Val]
	s.setsMutex.RUnlock()
	return &pb.Value{Val: value, Found: found}, nil
}

// Set key,value pair
// returns a pb.Value type
func (s *keyValServer) SetValue(ctx context.Context, kvpair *pb.KeyValuePair) (*pb.Value, error) {
	s.setsMutex.Lock()
	s.sets[kvpair.Key.Val] = kvpair.Value.Val
	s.setsMutex.Unlock()
	return kvpair.Value, nil
}

// main program entry
// starts the gRPC server
func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 6379))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterKeyValueServer(grpcServer, &keyValServer{sets: make(map[string]string)})
	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
