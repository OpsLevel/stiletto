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
		config, err := readStilettoInput()
		cobra.CheckErr(err)

		ctx := context.Background()
		// TODO: pipe dagger output a logfile on disk at the execution directory
		//client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
		client, err := dagger.Connect(ctx)
		if err != nil {
			panic(err)
		}
		defer client.Close()

		secretEngines := map[string]any{}
		for _, secret := range config.SecretEngines {
			secretEngines[secret.Name] = secret.Type
		}

		services := map[string]*dagger.Container{}
		for _, service := range config.Services {
			container := client.Container().From(service.Image)
			for _, mount := range service.Mounts {
				container = container.WithDirectory(mount.Container, client.Host().Directory(mount.Host))
			}
			for _, env := range service.Env {
				container = container.WithEnvVariable(env.Key, env.Value)
			}
			for _, port := range service.Ports {
				container = container.WithExposedPort(port.Port, dagger.ContainerWithExposedPortOpts{Protocol: dagger.NetworkProtocol(port.Protocol), Description: port.Name})
			}
			if service.Command != "" {
				container = container.WithExec(strings.Split(service.Command, " "))
			}
			services[service.Name] = container
		}

		for _, job := range config.Pipeline {
			log.Info().Msgf("Running job: %s", job.Name)

			container := client.Container().From(job.Image)
			for host, key := range job.Services {
				container = container.WithServiceBinding(host, services[key])
			}
			for _, mount := range job.Mounts {
				//log.Info().Msgf("With Volume Mount: %s:%s", mount.Host, mount.Container)
				container = container.WithDirectory(mount.Container, client.Host().Directory(mount.Host))
			}
			for _, cache := range job.Caches {
				//log.Info().Msgf("With Cache: %s:%s", cache.Name, cache.Path)
				container = container.WithMountedCache(cache.Path, client.CacheVolume(cache.Name))
			}
			for _, env := range job.Env {
				//log.Info().Msgf("With Env : %s:%s", key, value)
				container = container.WithEnvVariable(env.Key, env.Value)
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
