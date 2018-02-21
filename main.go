package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/numb3r3/h5-rtms-server/config"
	"github.com/numb3r3/h5-rtms-server/log"
	"github.com/numb3r3/h5-rtms-server/broker"
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
	logging.Info("start h5-rtms-server: ", version)

	cfg, err := config.ReadConfig(*configFileName, map[string]interface{}{
		"port":     9090,
		"hostname": "localhost",
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

	logging.Infof("Host: %s", cfg.GetString("hostname"))

	// Setup the new service
	svc, err := broker.NewService(cfg)
	if err != nil {
		panic(err.Error())
	}

	// Listen and serve
	svc.Listen()

}
