package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/fikrirnurhidayat/codeot-golang-executor/domain/service"
	"github.com/fikrirnurhidayat/codeot-golang-executor/server"
	"github.com/fikrirnurhidayat/codeotapis/proto"
)

var (
	err                error                                                                  // NOTE: Define error variable
	tcp                net.Listener                                                           // NOTE: Create TCP Listener
	srv                *grpc.Server                                                           // NOTE: Create grpc server variable
	mux                *runtime.ServeMux                                                      // NOTE: Set runtime mux
	dialOptions        []grpc.DialOption                                                      // NOTE: Create grpc dial options
	srvImpl            *server.Server                                                         // NOTE: Create example server variable
	grpcServerEndpoint = flag.String("grpc-server-endpoint", ":8080", "gRPC server endpoint") // NOTE: grpc server endpoint options
	logger             = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)       // NOTE: Initialize Logger
)

func init() {
	// TODO: Connect to Postgres

	// TODO: Connect to NSQ

	// TODO: Connect to Redis

	// NOTE: Initialize gRPC Dial Option
	dialOptions = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// NOTE: Initialize TCP Connection
	tcp, err = net.Listen("tcp", *grpcServerEndpoint)
	if err != nil {
		logger.Fatalf("net.Listen: cannot initialize tcp connection")
	}

	// NOTE: Create gRPC Server
	srv = grpc.NewServer()

	// NOTE: Create Mux Handler
	mux = runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	// TODO: Initialize Manager

	// TODO: Initialize Repositories

	// NOTE: Initialize Services
	executionService := service.NewExecutionService(logger)

	// NOTE: Initialize Server
	srvImpl = server.New(server.WithExecutionService(executionService))
}

func runGRPCServer() error {
	// NOTE: Register internal servers
	proto.RegisterGoexeServer(srv, srvImpl)

	return srv.Serve(tcp)
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func runGatewayServer(ctx context.Context) (err error) {
	// NOTE: Regsiter request handlers
	err = proto.RegisterGoexeHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, dialOptions)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    ":8081",
		Handler: cors(mux),
	}

	return srv.ListenAndServe()
}

func run() error {
	// NOTE: Setup context, so the requets can be cancelled
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// NOTE: Run grpc server as go routine
	go runGRPCServer()

	// NOTE: Start HTTP server (and proxy calls to gRPC server endpoint)
	return runGatewayServer(ctx)
}

func main() {
	flag.Parse()

	grpclog.SetLoggerV2(logger)

	if err := run(); err != nil {
		logger.Fatal(err)
	}
}
