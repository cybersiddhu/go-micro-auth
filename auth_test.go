package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/cybersiddhu/go-micro-auth/client"
	"github.com/cybersiddhu/go-micro-auth/service"

	"gopkg.in/jmoiron/sqlx.v0"
	_ "gopkg.in/mattn/go-sqlite3.v0"
)

func setUpSQLiteDB() string {
	file, err := ioutil.TempFile(os.TempDir(), "login")
	if err != nil {
		log.Fatal(err)
	}
	dbh := sqlx.MustConnect("sqlite3", file.Name())
	sfile := filepath.Join(currSrcDir(), "db", "user_sqlite3.sql")
	cnt, err := ioutil.ReadFile(sfile)
	if err != nil {
		log.Fatal(err)
	}
	_ = dbh.MustExec(string(cnt))
	return file.Name()
}

func currSrcDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("unable to retreive current src file path")
	}
	return filepath.Dir(filename)
}

func genKeyFile() string {
	file, err := ioutil.TempFile(os.TempDir(), "keyfile")
	if err != nil {
		log.Fatal(err)
	}
	prv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	err = pem.Encode(file, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(prv),
	})
	if err != nil {
		log.Fatal(err)
	}
	return file.Name()
}

func setUpAllWithSQLiteDB() (string, string) {
	return setUpSQLiteDB(), genKeyFile()
}

func tearDownAllWithSQLiteDB(ds string, key string) {
	os.Remove(ds)
	os.Remove(key)
}

func TestAuthenticationWithSQLiteDB(t *testing.T) {
	ds, key := setUpAllWithSQLiteDB()
	defer tearDownAllWithSQLiteDB(ds, key)
	conf := service.Config{
		DbDriver: "sqlite3",
		DbSource: ds,
		KeyFile:  key,
	}
	as := &service.AuthService{}
	handler, err := as.GetHttpHandler(conf)
	if err != nil {
		t.Errorf("Error in getting http handler %s\n", err)
	}
	server := httptest.NewServer(handler)
	client := &client.AuthClient{server.URL}
	email := "ryan@ryan.com"
	pass := "jogakhituchdi"
	//Sign up
	msg, err := client.SignUp(email, pass)
	if err != nil {
		t.Error(err)
	} else {
		if m, err := regexp.MatchString(email, msg); !m {
			t.Errorf("Did not match email %s error: %s", email, err)
		}
	}
	//Login
	token, err := client.Login(email, pass)
	if err != nil {
		t.Error(err)
	} else {
		el := strings.Split(token, ".")
		if len(el) != 3 {
			t.Errorf("Expected 3 elements got %d", len(el))
		} else {
			rgxp := regexp.MustCompile(`^\S+$`)
			for _, v := range el {
				if !rgxp.MatchString(v) {
					t.Errorf("Did not match %s part", v)
				}
			}
		}
	}
}
