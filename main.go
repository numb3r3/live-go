package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/numb3r3/h5-rtms-server/config"
	"github.com/numb3r3/h5-rtms-server/log"
)

var (
	version        = "master"
	configfilename = flag.String("config", "default_config.yaml", "configure filename")
	loglevel       = flag.String("loglevel", "info", "log level")
	logfile        = flag.String("logfile", "console.log", "log file path")
	argHelp        = flag.Bool("help", false, "Shows the help and usage instead of running the broker.")
)

func init() {
	flag.Parse()
	if *argHelp {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// logging.SetOutputByName(*logfile)
	// logging.SetRotateByDay()
	// logging.SetLevelByString(*loglevel)

}

func main() {
	logging.Info("start h5-rtms-server: ", version)

	cfg, err := config.ReadConfig("default_config.yaml", map[string]interface{}{
		"port":     9090,
		"hostname": "localhost",
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

}
