/*
Copyright © 2023 fahrizalfarid
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fahrizalfarid/user-service-grpc/conf"
	"github.com/fahrizalfarid/user-service-grpc/src/api/server"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("api called")

		err := conf.LoadEnv("./.env")
		if err != nil {
			panic(err)
		}
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		go func() {
			log.Fatal(server.RunServer().Start(":8080"))
		}()

		<-quit
		fmt.Println("shutdown")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
