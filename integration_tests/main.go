package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	origName := flag.Bool("orig", false, "tell server to use origin name in jsonpb")
	flag.Parse()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	l, err := net.Listen("tcp4", "localhost:9000")
	if err != nil {
		panic(err)
	}

	rcs := &RealCounterService{}
	s := grpc.NewServer()
	RegisterCounterServiceServer(s, rcs)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	go func() {
		err := s.Serve(l)
		if err != nil {
			panic(err)
		}
	}()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		OrigName: *origName,
	}))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = RegisterCounterServiceHandlerFromEndpoint(ctx, mux, "localhost:9000", opts)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8081", allowCORS(mux))
	if err != nil {
		panic(err)
	}

}
