package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Server configuration values
type Configuration struct {
	// Database connection values
	Host     string
	Port     string
	Dbname   string
	Sslmode  string
	User     string
	Password string

	// Session keys directory path.
	// RSA keys could be generated using following commands:
	// openssl genrsa -out app.rsa keysize
	// openssl rsa -in app.rsa -pubout > app.rsa.pub
	SessionKeysPath string
}

func GetConfigValues() (string, []byte, []byte) {
	file, err := os.Open("config/config.json")
	if err != nil {
		panic(err)
	}
	conf := Configuration{}
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		panic(err)
	}

	userTokenSignKey, err := ioutil.ReadFile(conf.SessionKeysPath + "/app.rsa")
	if err != nil {
		panic(err)
	}

	userTokenVerifyKey, err := ioutil.ReadFile(conf.SessionKeysPath + "/app.rsa.pub")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("host=%s port=%s dbname=%s "+
		"sslmode=%s user=%s password=%s ",
		conf.Host,
		conf.Port,
		conf.Dbname,
		conf.Sslmode,
		conf.User,
		conf.Password,
	), userTokenSignKey, userTokenVerifyKey
}
