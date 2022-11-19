/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mvcris/biziu/internal/tcp"
	"github.com/spf13/cobra"
)

var requests uint32
var concurrency uint32
var nodes uint32
var port uint16

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example`,
	Run: func(cmd *cobra.Command, args []string) {
		server := tcp.NewTcpServer(requests, concurrency, nodes, port)
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().Uint32VarP(&requests, "requests", "r", 0, "total number of requests")
	serverCmd.PersistentFlags().Uint32VarP(&concurrency, "concurrency", "c", 0, "total number of concurrency requests")
	serverCmd.PersistentFlags().Uint32VarP(&nodes, "nodes", "n", 0, "total number of nodes")
	serverCmd.PersistentFlags().Uint16VarP(&port, "port", "p", 3000, "port to run server, default is 3000")
	serverCmd.MarkFlagRequired("r")
	serverCmd.MarkFlagRequired("concurrency")
	serverCmd.MarkFlagRequired("nodes")
}
