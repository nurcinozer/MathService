package service

import (
	"context"

	"github.com/go-kit/kit/log"
)

// Struct representing our service
type service struct {
	logger log.Logger
}

// Service Interface for our service
type Service interface {
	// Add The add method takes in two 32bit floats and returns
	// a float + an error (if there is one)
	Add(ctx context.Context, numA, numB float32) (float32, error)
	Sub(ctx context.Context, numA, numB float32) (float32, error)
	Div(ctx context.Context, numA, numB float32) (float32, error)
	Mul(ctx context.Context, numA, numB float32) (float32, error)
}

// Add Implement add method in Service interface
func (s service) Add(ctx context.Context, numA, numB float32) (float32, error) {
	s.logger.Log("INFO", "Adding up request...")
	return numA + numB, nil
}

func (s service) Sub(ctx context.Context, numA, numB float32) (float32, error) {
	s.logger.Log("INFO", "Subtracting request...")
	return numA - numB, nil
}

func (s service) Div(ctx context.Context, numA, numB float32) (float32, error) {
	s.logger.Log("INFO", "Dividing request...")
	return numA / numB, nil
}

func (s service) Mul(ctx context.Context, numA, numB float32) (float32, error) {
	s.logger.Log("INFO", "Multiplying request...")
	return numA * numB, nil
}

// NewService Returns a service instance (struct) with our logger
func NewService(logger log.Logger) Service {
	return &service{
		logger: logger,
	}
}
