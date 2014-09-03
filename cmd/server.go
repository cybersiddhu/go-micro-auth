package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cybersiddhu/go-micro-auth/service"
	"gopkg.in/codegangsta/cli.v0"
	"gopkg.in/yaml.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "microauth"
	app.Usage = "App to work with HTTP authentication mircoservice"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{"config", "c", "config.yaml", "Name of the config file"},
	}
	app.Commands = []cli.Command{
		{
			Name:   "server",
			Usage:  "Run the HTTP server",
			Action: runServer,
		},
		{
			Name:   "generate-key",
			Usage:  "Generate private key file",
			Action: genPrivateKeyFile,
		},
		{
			Name:   "create-table",
			Usage:  "Create user table in the database",
			Action: createUserTable,
		},
	}
	app.Run(os.Args)
}

func runServer(c *cli.Context) {
	config, err := readConfig(c)
	if err != nil {
		log.Fatal(err)
	}
	svc := &service.AuthService{}
	handler, err := svc.GetHttpHandler(config)
	if err != nil {
		log.Fatal(err)
	}
	http.ListenAndServe(config.ServiceHost, handler)
}

func createUserTable(c *cli.Context) {
	conf, err := readConfig(c)
	if err != nil {
		log.Fatal(err)
	}
	// figure out the schema file
	parentDir, _ := filepath.Split(currSrcDir())
	schema := filepath.Join(parentDir, "db", "user_"+conf.DbDriver+".sql")
	ct, err := ioutil.ReadFile(schema)
	if err != nil {
		log.Fatal(err)
	}
	//connect to database
	dbh, err := service.GetDBHandler(conf)
	if err != nil {
		log.Fatalf("unable to connect to database with error %s\n", err)
	}
	//load the schema
	tx := dbh.MustBegin()
	_ = tx.MustExec(string(ct))
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func currSrcDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("unable to retreive current src file path")
	}
	return filepath.Dir(filename)
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
