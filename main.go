package main

import (
  "net"
  "fmt"
  "flag"
  "os"
  "io/ioutil"

  "github.com/fikrirnurhidayat/codeotapis/proto"
  "github.com/fikrirnurhidayat/codeot-golang-executor/app/server"

  "google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var port = flag.Int("port", 50051, "The server port")
var log grpclog.LoggerV2

func main() {
  // Get the flag
  flag.Parse()

  // Set logger
	log = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)

  // Setup tcp listener
	tcp, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

  // Instantiate grpc server
	s := grpc.NewServer()

  // Instantiate server
  server := server.New()

  // register server to grpc
	proto.RegisterGoexeServer(s, server)

	log.Infof("server listening at %v", tcp.Addr())
	if err := s.Serve(tcp); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
