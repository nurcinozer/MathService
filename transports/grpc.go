package transports

import (
	"context"
	"github.com/go-kit/kit/log"
	kittransport "github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"math_service/endpoints"
	"math_service/pb"
)

// Defines a gRPC transport
type gRPCTransport struct {
	add kitgrpc.Handler
	sub kitgrpc.Handler
	div kitgrpc.Handler
	mul kitgrpc.Handler
	pb.UnimplementedMathServiceServer
}

// Add Defines the add handler
func (s *gRPCTransport) Add(ctx context.Context, req *pb.MathRequest) (*pb.MathResponse, error) {
	_, resp, err := s.add.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp.(*pb.MathResponse), nil
}

// Sub Defines the sub handler
func (s *gRPCTransport) Sub(ctx context.Context, req *pb.MathRequest) (*pb.MathResponse, error) {
	_, resp, err := s.sub.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp.(*pb.MathResponse), nil
}

// Div Defines the div handler
func (s *gRPCTransport) Div(ctx context.Context, req *pb.MathRequest) (*pb.MathResponse, error) {
	_, resp, err := s.div.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp.(*pb.MathResponse), nil
}

// Mul Defines the mul handler
func (s *gRPCTransport) Mul(ctx context.Context, req *pb.MathRequest) (*pb.MathResponse, error) {
	_, resp, err := s.mul.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp.(*pb.MathResponse), nil
}

// MathRequest decoder
func decodeMathRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.MathRequest)
	return endpoints.MathRequest{
		NumA: req.NumA,
		NumB: req.NumB,
	}, nil
}

// MathResponse encoder
func encodeMathResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.MathResponse)
	return &pb.MathResponse{Result: resp.Result}, nil
}

// NewGRPCTransport Inits gRPC transport
func NewGRPCTransport(endpoints endpoints.Endpoints, logger log.Logger) pb.MathServiceServer {
	// Add transport type to our logger
	logger = log.With(logger, "transport", "gRPC")

	options := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
	}

	return &gRPCTransport{
		add: kitgrpc.NewServer(
			endpoints.Add,
			decodeMathRequest,
			encodeMathResponse,
			// Spread options
			append(options)...,
		),
		sub: kitgrpc.NewServer(
			endpoints.Sub,
			decodeMathRequest,
			encodeMathResponse,
			// Spread options
			append(options)...,
		),
		div: kitgrpc.NewServer(
			endpoints.Div,
			decodeMathRequest,
			encodeMathResponse,
			append(options)...,
		),
		mul: kitgrpc.NewServer(
			endpoints.Mul,
			decodeMathRequest,
			encodeMathResponse,
			append(options)...),
	}
}
