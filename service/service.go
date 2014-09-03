package service

import (
	"io/ioutil"

	"gopkg.in/gin-gonic/gin.v0"
	"gopkg.in/jmoiron/sqlx.v0"
	_ "gopkg.in/lib/pq.v0"
	_ "gopkg.in/mattn/go-sqlite3.v0"
)

type Config struct {
	ServiceHost string `yaml:"host,flow"`
	DbDriver    string `yaml:"driver,flow"`
	DbSource    string `yaml:"datasource,flow"`
	KeyFile     string
}

type AuthService struct{}

func (s *AuthService) Run(conf Config) error {
	// database connection
	dbh, err := GetDBHandler(conf)
	if err != nil {
		return err
	}

	// reading the private key
	keyData, err := ioutil.ReadFile(conf.KeyFile)
	if err != nil {
		return err
	}

	// setup resource
	resource := &AuthResource{dbh, keyData}
	r := gin.Default()
	auth := r.Group("/auth")
	auth.POST("/login", resource.CreateSession)
	auth.POST("/signup", resource.CreateUser)
	r.Run(conf.ServiceHost)
	return nil
}

func GetDBHandler(conf Config) (*sqlx.DB, error) {
	dbh, err := sqlx.Connect(conf.DbDriver, conf.DbSource)
	return dbh, err
}
