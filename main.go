package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/numb3r3/live-go/broker"
	"github.com/numb3r3/live-go/config"
	"github.com/numb3r3/live-go/log"
)

var (
	version        = "master"
	configFileName = flag.String("config", "default_config", "configure filename")
	logLevel       = flag.String("loglevel", "info", "log level")
	logFile        = flag.String("logfile", "console.log", "log file path")
	argHelp        = flag.Bool("help", false, "Shows the help and usage instead of running the broker.")
)

func init() {
	flag.Parse()
	if *argHelp {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// logging.SetOutputByName(*logFile)
	// logging.SetRotateByDay()
	// logging.SetLevelByString(*logLevel)

}

func main() {
	logging.Info("start live-go: ", version)

	cfg, err := config.ReadConfig(*configFileName, map[string]interface{}{
		"listen_addr": "0.0.0.0:9090",
		"auth": map[string]string{
			"username": "numb3r3",
			"password": "314159",
		},
	})
	if err != nil {
		logging.Fatal(err)
		panic(fmt.Errorf("Error when reading config: %v", err))
	}

	// Setup the new service
	svc, err := broker.NewService(cfg)
	if err != nil {
		panic(err.Error())
	}

	// Listen and serve
	svc.Listen()

}
