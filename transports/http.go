package transports

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	kittransport "github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"math_service/endpoints"
	"net/http"
)

type wrappedError struct {
	Error string `json:"error"`
}

// We don't need (afaik) an error encoder for gRPC as errors are returned
// in "native" format
func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	// This is where we'd put the logic for encoding errors, but because this is
	// a toy example. We'll write a response containing a 500
	w.WriteHeader(http.StatusInternalServerError)
	// Encode error into JSON + write it to our writer (w)
	json.NewEncoder(w).Encode(wrappedError{Error: err.Error()})

}

// This takes a JSON request and decodes it into a MathRequest struct
// much like in gRPC we need to take a pointer to a MathRequest message + turn it into
// our MathRequest struct
func decodeHTTPMathRequest(_ context.Context, request *http.Request) (interface{}, error) {
	var decodedReq endpoints.MathRequest
	// Try decoding JSON into MathRequest (decoded_req)
	err := json.NewDecoder(request.Body).Decode(&decodedReq)
	return decodedReq, err
}

// NewHTTPTransport Unlike our gRPC transport, we don't need a struct
// to define individual handlers, we'll just
// use a function to generate a combined handler (router or mux)
// and return it
func NewHTTPTransport(endpoints endpoints.Endpoints, logger log.Logger) http.Handler {
	// Add transport type to our logger
	logger = log.With(logger, "transport", "HTTP")

	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
		kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
	}

	// Instantiate new mux (router)
	mux := http.NewServeMux()
	// Tell mux how to handle /add endpoint
	mux.Handle("/add", kithttp.NewServer(
		endpoints.Add,
		decodeHTTPMathRequest,
		kithttp.EncodeJSONResponse,
		// Spread options
		append(options)...,
	))
	// Tell mux how to handle /sub endpoint
	mux.Handle("/sub", kithttp.NewServer(
		endpoints.Sub,
		decodeHTTPMathRequest,
		kithttp.EncodeJSONResponse,
		append(options)...,
	))
	// Tell mux how to handle /div endpoint
	mux.Handle("/div", kithttp.NewServer(
		endpoints.Div,
		decodeHTTPMathRequest,
		kithttp.EncodeJSONResponse,
		append(options)...,
	))
	mux.Handle("/mul", kithttp.NewServer(
		endpoints.Mul,
		decodeHTTPMathRequest,
		kithttp.EncodeJSONResponse,
		append(options)...,
	))
	// Return our mux
	return mux
}
