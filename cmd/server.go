package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/codegangsta/cli.v0"
	"gopkg.in/yaml.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "mauth"
	app.Usage = "App to work with HTTP authentication mircoservice"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{"config", "c", "config.yaml", "Name of the config file"},
	}
	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run the HTTP server",
			Action: func(c *cli.Context) {
				config, err := readConfig(c)
				if err != nil {
					log.Fatal(err)
				}
				svc := service.AuthService{}
				if err = svc.Run(config); err != nil {
					log.Fatal(err)
				}

			},
		},
	}
	app.Run(os.Args)
}

func readConfig(c *cli.Context) (service.Config, error) {
	ypath := c.GlobalString("config")
	config := service.Config{}
	if _, err := os.Stat(ypath); err != nil {
		return config, errors.New("config file path is not valid")
	}
	ydata, err := ioutil.ReadFile(ypath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal([]byte(ydata), &config)
	return config, err
}
