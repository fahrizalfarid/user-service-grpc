/*
Copyright © 2023 fahrizalfarid
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/fahrizalfarid/user-service-grpc/conf"
	api "github.com/fahrizalfarid/user-service-grpc/src/api/server"
	user "github.com/fahrizalfarid/user-service-grpc/src/user-service/server"
	validator "github.com/fahrizalfarid/user-service-grpc/src/validator-service/server"
	"github.com/spf13/cobra"
)

// allSvcCmd represents the allSvc command
var allSvcCmd = &cobra.Command{
	Use:   "allSvc",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runAllSrv,
}

var port int

func init() {
	rootCmd.AddCommand(allSvcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allSvcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allSvcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	allSvcCmd.Flags().IntVarP(&port, "port", "p", 8080, "")
}

func runAllSrv(cmd *cobra.Command, args []string) {
	err := conf.LoadEnv("./.env")
	if err != nil {
		panic(err)
	}

	go func() {
		log.Fatal(user.RunUserSrv(conf.GetUserSrv()))
	}()

	go func() {
		log.Fatal(validator.RunUserValidatorSrv(conf.GetValidatorSrv()))
	}()

	go func() {
		log.Fatal(api.RunServer().Start(fmt.Sprintf(":%d", port)))
	}()

	fmt.Println("all services called")
	select {}
}
