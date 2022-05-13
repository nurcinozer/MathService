package main

import (
	"fmt"
	"math_service/endpoints"
	"math_service/pb"
	"math_service/service"
	"math_service/transports"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
)

// Main isn't what it would be in most cases. In Go-kit, we use main to "wire up"
// our services to our endpoints to our transports. Init loggers/tracing, etc.
func main() {

	// Instantiate logger that will output JSON to console
	var logger log.Logger = log.NewJSONLogger(os.Stdout)
	// Add context
	logger = log.With(logger, "time_stamp", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// Instantiate our service with the above logger
	mathservice := service.NewService(logger)
	// Instantiate endpoints with the above service
	mathendpoints := endpoints.MakeEndpoints(mathservice)
	// Instantiate a new gRPC transport with the above endpoints
	gRPCTransport := transports.NewGRPCTransport(mathendpoints, logger)
	// Instantiate a new HTTP (REST) transport with the above endpoints
	httpTransport := transports.NewHTTPTransport(mathendpoints, logger)
	// Create channel that will be used to communicate errors
	errs := make(chan error)

	// Track term/alarm signals in new goroutine
	go func() {
		// Create a channel
		c := make(chan os.Signal)
		// Start tracking signals, and notify on channel c
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		// Listen to channel c and push error into errs
		errs <- fmt.Errorf("%s", <-c)
	}()

	// Instantiate tcp listener on port 12345 for gRPC
	grpcListener, err := net.Listen("tcp", ":12345")
	if err != nil {
		logger.Log("during", "Listen for gRPC", "ERROR", err)
		os.Exit(1)
	}
	// Instantiate tcp listener on port 8080 for HTTP (REST)
	httpListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Log("during", "Listen for HTTP", "ERROR", err)
	}

	// Launch gRPC server in a new goroutine
	go func() {
		// Close listener on exit
		defer grpcListener.Close()
		// Instantiate our server
		srv := grpc.NewServer()
		// Register our server with service
		// server uses gRPC transport
		pb.RegisterMathServiceServer(srv, gRPCTransport)
		level.Info(logger).Log("msg", "gRPC Server started...")
		// Start serving using gRPC listener
		err := srv.Serve(grpcListener)
		if err != nil {
			logger.Log("during", "Launch gRPC", "ERROR", err)
		}
	}()

	// Launch HTTP server in a new goroutine
	go func() {
		// Close listener on exit
		defer httpListener.Close()
		level.Info(logger).Log("msg", "HTTP Server started...")
		// Start serving via HTTP using our mux (group of handlers)
		err := http.Serve(httpListener, httpTransport)
		if err != nil {
			logger.Log("during", "Launch HTTP", "ERROR", err)
		}
	}()
	// Wait here until there's an error or we quit
	level.Error(logger).Log("INFO", <-errs)
}
