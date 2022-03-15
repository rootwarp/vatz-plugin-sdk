package sdk

import (
	"context"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/rootwarp/vatz-plugin-sdk/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PluginServer is...
type PluginServer struct {
	pb.UnimplementedManagerPluginServer

	callbacks []func(map[string]interface{}) error
}

// Init initializes plugin.
func (s *PluginServer) Init(context.Context, *emptypb.Empty) (*pb.PluginInfo, error) {
	// TODO: TBD
	return nil, nil
}

// Verify returns liveness.
func (s *PluginServer) Verify(context.Context, *emptypb.Empty) (*pb.VerifyInfo, error) {
	// TODO: TBD
	return nil, nil
}

// Execute runs plugin features.
func (s *PluginServer) Execute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	log.Println("PluginServer.Execute")

	resp := &pb.ExecuteResponse{
		State:   pb.ExecuteResponse_SUCCESS,
		Message: "OK",
	}

	// TODO: Dummy parameter.
	param := map[string]interface{}{
		"function": "IsBorUp",
	}

	for _, f := range s.callbacks {
		f(param)
	}

	return resp, nil
}

// Start starts gRPC service.
func (s *PluginServer) Start() error {
	log.Println("Start vatz-matic-plugin")

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := net.Listen("tcp", "0.0.0.0:9091")
	if err != nil {
		log.Println(err)
	}

	ch := make(chan os.Signal, 1)
	srv := grpc.NewServer()

	pb.RegisterManagerPluginServer(srv, s)

	reflection.Register(srv)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		_ = <-ch
		cancel()
		srv.GracefulStop()
		wg.Done()
	}()

	if err := srv.Serve(c); err != nil {
		log.Panic(err)
	}

	wg.Wait()

	return nil
}

// Register registers callback function.
func (s *PluginServer) Register(cb func(map[string]interface{}) error) error {
	log.Println("RegisterFeature function")

	s.callbacks = append(s.callbacks, cb)

	log.Println("After register", len(s.callbacks))
	return nil
}

// NewPlugin creates new plugin service instance.
func NewPlugin() *PluginServer {
	return &PluginServer{
		callbacks: make([]func(map[string]interface{}) error, 0),
	}
}
