//go:generate protoc -I=./ --go_out=../ --go-grpc_out=../ ./simulator_controller.proto
package controller

import (
	context "context"
	"net"
	"net/url"

	"github.com/obicons/rmck/sim"
	"google.golang.org/grpc"
)

type SimulatorController struct {
	url        *url.URL
	grpcServer *grpc.Server
	simulator  sim.Sim
	listener   net.Listener
}

func New(addrStr string, simulator sim.Sim) (*SimulatorController, error) {
	var err error
	server := SimulatorController{}
	server.url, err = url.Parse(addrStr)
	if err != nil {
		return nil, err
	}
	server.simulator = simulator
	return &server, nil

}

// Starts the SimulatorController.
// It is an error to call this method if server has already been started.
func (server *SimulatorController) Start() error {
	var err error
	// TODO -- clean this up
	server.listener, err = net.Listen(server.url.Scheme, server.url.Path)
	if err != nil {
		return err
	}
	server.grpcServer = grpc.NewServer()
	service := NewSimulatorControllerService(server)
	RegisterSimulatorControllerService(server.grpcServer, service)
	return server.grpcServer.Serve(server.listener)
}

// Stops the SimulatorController.
// It is an error to call this method if server has not been started.
func (server *SimulatorController) Stop() {
	server.grpcServer.GracefulStop()
	server.listener.Close()
}

// Implements RPC
func (s *SimulatorController) Step(ctx context.Context, req *StepRequest) (*StepResponse, error) {
	err := s.simulator.Step(ctx)
	return &StepResponse{}, err
}

// Implements RPC
func (s *SimulatorController) Position(ctx context.Context, req *PositionRequest) (*PositionResponse, error) {
	return &PositionResponse{}, nil
}