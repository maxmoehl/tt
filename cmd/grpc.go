package cmd

import (
	"fmt"
	"net"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "[alpha] Launch a gRPC server",
	RunE: func(_ *cobra.Command, _ []string) error {
		err := runGrpc()
		if err != nil {
			return fmt.Errorf("grpc: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
}

func runGrpc() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	tt.RegisterTtServer(grpcServer, &tt.GrpcServer{})
	return grpcServer.Serve(lis)
}
