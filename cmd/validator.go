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
	validator "github.com/fahrizalfarid/user-service-grpc/src/validator-service/server"
	"github.com/spf13/cobra"
)

// validatorCmd represents the validator command
var validatorCmd = &cobra.Command{
	Use:   "validator",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validator called")

		err := conf.LoadEnv("./.env")
		if err != nil {
			panic(err)
		}
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		go func() {
			log.Fatal(validator.RunUserValidatorSrv(conf.GetValidatorSrv()))
		}()

		<-quit
		fmt.Println("shutdown")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(validatorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validatorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
