package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	core_grpcclient "github.com/cnc-csku/task-nexus/go-lib/grpcclient"
	"github.com/cnc-csku/task-nexus/go-lib/utils/network"
	"github.com/cnc-csku/task-nexus/task-management/internal/wire"
)

func main() {
	app := wire.InitializeApp()

	grpcServer := wire.InitializeGrpcServer()
	defer grpcServer.Server.Stop()

	// create a listener on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcServer.Config.GrpcServer.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	// create a context that will be can	celed when SIGINT or SIGTERM is received
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// get local IP address to display on startup
	localIP, _ := network.GetLocalIP()

	// start gRPC server
	grpcReady := make(chan bool)
	go func() {
		log.Printf("✅ gRPC server is running on %s:%s", localIP, grpcServer.Config.GrpcServer.Port)
		close(grpcReady)

		if err := grpcServer.Server.Serve(lis); err != nil {
			fmt.Println("error", err)
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// start the rest server
	go func() {
		<-grpcReady // wait for grpc server to be ready

		err = app.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// wait for SIGINT or SIGTERM
	<-ctx.Done()

	// gracefully shutdown gRPC server
	grpcServer.Server.GracefulStop()

	// cancel context after the server gracefully stopped
	stop()

	// close all gRPC client connections
	core_grpcclient.CloseAllGrpcConnections()
}
