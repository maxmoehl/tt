//go:build grpc

package cmd

import (
	"fmt"
	"net"

	ttGrpc "github.com/maxmoehl/tt/grpc"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "[alpha] Launch a gRPC server",
	Long: `[alpha] Launch a gRPC server. So far it's not actually
doing anything and it's just here as a dummy/test functionality.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := getGrpcParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("grpc: %w", err)
		}
		err = runGrpc(port)
		if err != nil {
			return fmt.Errorf("grpc: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
}

func runGrpc(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	ttGrpc.RegisterTtServer(grpcServer, &ttGrpc.Server{})
	return grpcServer.Serve(lis)
}

func getGrpcParameters(cmd *cobra.Command, _ []string) (port int, err error) {
	flags, err := flags(cmd, flagPort)
	if err != nil {
		return
	}
	return flags[flagPort].(int), nil
}
