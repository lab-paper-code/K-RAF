package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/lab-paper-code/ksv/volume-service/commons"
	"github.com/lab-paper-code/ksv/volume-service/db"
	"github.com/lab-paper-code/ksv/volume-service/k8s"
	"github.com/lab-paper-code/ksv/volume-service/logic"
	"github.com/lab-paper-code/ksv/volume-service/rest"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var config *commons.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "volume-service [args..]",
	Short: "Provides controls of volumes via REST interface",
	Long:  `Volume-Serivce provides controls of volumes via REST interface.`,
	RunE:  processCommand,
}

func Execute() error {
	return rootCmd.Execute()
}

func processCommand(command *cobra.Command, args []string) error {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "processCommand",
	})

	cont, err := processFlags(command)
	if err != nil {
		logger.Error(err)
	}

	if !cont {
		return err
	}

	// start service
	logger.Info("Starting DB Adapter...")
	dbAdapter, err := db.Start(config)
	if err != nil {
		logger.Fatal(err)
	}
	defer dbAdapter.Stop()
	logger.Info("DB Adapter Started")

	logger.Info("Starting K8S Adapter...")
	var k8sAdapter *k8s.K8SAdapter
	if config.NoKubernetes {
		logger.Info("Disable K8S Adapter")
	} else {
		k8sAdapter, err = k8s.Start(config)
		if err != nil {
			logger.Fatal(err)
		}
		defer k8sAdapter.Stop()
		logger.Info("K8S Adapter Started")
	}

	logik, err := logic.Start(config, dbAdapter, k8sAdapter)
	if err != nil {
		logger.Fatal(err)
	}
	defer logik.Stop()

	logger.Info("Starting REST Adapter...")
	restAdapter, err := rest.Start(config, logik)
	if err != nil {
		logger.Fatal(err)
	}
	defer restAdapter.Stop()
	logger.Info("REST Adapter Started")

	// wait
	fmt.Println("Press Ctrl+C to stop...")
	waitForCtrlC()

	return nil
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})

	config = commons.GetDefaultConfig()
	log.SetLevel(config.GetLogLevel())

	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "main",
	})

	// attach common flags
	setCommonFlags(rootCmd)

	err := Execute()
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}

func waitForCtrlC() {
	var endWaiter sync.WaitGroup

	endWaiter.Add(1)
	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		<-signalChannel
		endWaiter.Done()
	}()

	endWaiter.Wait()
}
