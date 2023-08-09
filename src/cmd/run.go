/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// initialize Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		if err != nil {
			panic(err)
		}
		defer client.Close()

		// use a golang:1.19 container
		// get version
		// execute
		golang := client.Container().From("golang:1.19").WithExec([]string{"go", "version"})

		version, err := golang.Stdout(ctx)
		if err != nil {
			panic(err)
		}

		// print output
		fmt.Println("Hello from Dagger and " + version)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
