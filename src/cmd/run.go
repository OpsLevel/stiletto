/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/opslevel/stiletto/pkg"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var stilettoFile string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		jobs, err := readStilettoInput()
		cobra.CheckErr(err)

		ctx := context.Background()
		//client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
		client, err := dagger.Connect(ctx)
		if err != nil {
			panic(err)
		}
		defer client.Close()

		for _, job := range jobs.Jobs {
			log.Info().Msgf("Running job: %s", job.Name)

			container := client.Container().From(job.Image)
			for _, mount := range job.Mounts {
				//log.Info().Msgf("With Volume Mount: %s:%s", mount.Host, mount.Container)
				container = container.WithDirectory(mount.Container, client.Host().Directory(mount.Host))
			}
			for _, cache := range job.Caches {
				//log.Info().Msgf("With Cache: %s:%s", cache.Name, cache.Path)
				container = container.WithMountedCache(cache.Path, client.CacheVolume(cache.Name))
			}
			for key, value := range job.Env {
				//log.Info().Msgf("With Env : %s:%s", key, value)
				container = container.WithEnvVariable(key, value)
			}
			container = container.WithWorkdir(job.Workdir)
			for _, command := range job.Commands {
				//log.Info().Msgf("Running command: %s", command)
				container = container.WithExec(strings.Split(command, " "))
			}
			out, err := container.Stdout(ctx)
			cobra.CheckErr(err)
			log.Info().Msg(out)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&stilettoFile, "file", "f", ".", "File to read data from. If '-' then reads from stdin. Defaults to read from './stiletto.yaml'")

}

func readStilettoInput() (*pkg.Stiletto, error) {
	if stilettoFile == "" {
		return nil, fmt.Errorf("please specify a stiletto.yaml file")
	}
	if stilettoFile == "-" {
		viper.SetConfigType("yaml")
		viper.ReadConfig(os.Stdin)
	} else if stilettoFile == "." {
		viper.SetConfigFile("./job.yaml")
	} else {
		viper.SetConfigFile(stilettoFile)
	}
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	job := &pkg.Stiletto{}
	viper.Unmarshal(&job)
	return job, nil
}
