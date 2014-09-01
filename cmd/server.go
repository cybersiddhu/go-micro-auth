package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/cybersiddhu/go-micro-auth/service"

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
				svc := &service.AuthService{}
				if err = svc.Run(config); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:   "genkey",
			Usage:  "Generate private key file",
			Action: genPrivateKeyFile,
		},
	}
	app.Run(os.Args)
}

// Generate a 2048 bit RSA private key and write to the file given
// in the config file
func genPrivateKeyFile(c *cli.Context) {
	config, err := readConfig(c)
	if err != nil {
		log.Fatal(err)
	}
	out, err := os.Create(config.KeyFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	prv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	err = pem.Encode(out, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(prv),
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Reads the yaml config file from the path given in command line
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
