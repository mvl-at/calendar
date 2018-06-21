package calendar

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"
)

//Path for the config file.
const ConfigPath = "conf.json"

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var errLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
var conf = config()

//Reads the config from file and assigns it to the context.Conf
func config() (conf *Configuration) {
	conf = &Configuration{}
	fil, err := os.OpenFile(ConfigPath, 0, 0644)
	defer fil.Close()

	if err != nil {
		fil, err = os.Create(ConfigPath)
		defer fil.Close()
		rand.Seed(time.Now().UnixNano())
		jwtSecret := make([]byte, 8)
		rand.Read(jwtSecret)
		conf = &Configuration{
			Host:     "0.0.0.0",
			Port:     8081,
			RestHost: "127.0.0.1:8080",
			Threads:  0}
		enc := json.NewEncoder(fil)
		enc.SetIndent("", "  ")
		err = enc.Encode(conf)

	} else {
		err = json.NewDecoder(fil).Decode(conf)
	}

	if err != nil {
		errLogger.Fatalln(err.Error())
	}
	return
}

//Struct which holds the configuration.
type Configuration struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	RestHost string `json:"restHost"`
	Threads  int    `json:"threads"`
}
