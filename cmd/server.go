/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mvcris/biziu/internal/tcp"
	"github.com/spf13/cobra"
)

var file string
var port uint16

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example`,
	Run: func(cmd *cobra.Command, args []string) {
		server := tcp.NewTcpServer(file, port)
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "file with request template")
	serverCmd.PersistentFlags().Uint16VarP(&port, "port", "p", 0, "server port")
	serverCmd.MarkFlagRequired("r")
	serverCmd.MarkFlagRequired("concurrency")
	serverCmd.MarkFlagRequired("nodes")
}
