package endpoints

import (
	"context"
	"math_service/service"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints This struct hold a list of all endpoint defs
type Endpoints struct {
	Add endpoint.Endpoint
	Sub endpoint.Endpoint
	Div endpoint.Endpoint
	Mul endpoint.Endpoint
}

// MathRequest MathReq type represents our MathRequest gRPC message
type MathRequest struct {
	NumA float32
	NumB float32
}

// MathResponse type represents our MathResponse gRPC message
type MathResponse struct {
	Result float32
}

// makeAddEndpoint closure that processes a request with the Add service method
// and returns the result, wraps Add service
func makeAddEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// Type assert request interface into MathReq
		req := request.(MathRequest)
		// Call Add service method
		result, _ := s.Add(ctx, req.NumA, req.NumB)
		// Return a Math
		return MathResponse{Result: result}, nil
	}
}

func makeSubEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(MathRequest)
		result, _ := s.Sub(ctx, req.NumA, req.NumB)
		return MathResponse{Result: result}, nil
	}
}

func makeDivEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(MathRequest)
		result, _ := s.Div(ctx, req.NumA, req.NumB)
		if result == 0 {
			return MathResponse{Result: 50}, nil
		}
		return MathResponse{Result: result}, nil
	}
}

func makeMulEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(MathRequest)
		result, err := s.Mul(ctx, req.NumA, req.NumB)
		if err != nil {
			return nil, err
		}
		return MathResponse{Result: result}, nil
	}
}

// MakeEndpoints Inits and pools all our endpoint instances
func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Add: makeAddEndpoint(s),
		Sub: makeSubEndpoint(s),
		Div: makeDivEndpoint(s),
		Mul: makeMulEndpoint(s),
	}
}
