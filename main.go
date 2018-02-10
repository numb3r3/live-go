package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/numb3r3/h5-rtms-server/config"
)

var (
	version        = "master"
	configfilename = flag.String("config", "default_config.yaml", "configure filename")
	loglevel       = flag.String("loglevel", "info", "log level")
	logfile        = flag.String("logfile", "h5-rtms-server.log", "log file path")
	argHelp        = flag.Bool("help", false, "Shows the help and usage instead of running the broker.")
)

func init() {
	flag.Parse()
	if *argHelp {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.SetOutputByName(*logfile)
	log.SetRotateByDay()
	log.SetLevelByString(*loglevel)
}

func main() {
	log.Info("start h5-rtms-server: ", version)

	cfg, err := config.readConfig("default_config.yaml", map[string]interface{}{
		"port":     9090,
		"hostname": "localhost",
		"auth": map[string]string{
			"username": "titpetric",
			"password": "12fa",
		},
	})
	if err != nil {
		log.Fatal(err)
		panic(fmt.Errorf("Error when reading config: %v", err))
	}

}
