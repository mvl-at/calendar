package calendar

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	ConfigPath     = "conf.json"
	StaticThreads  = "static"
	DynamicThreads = "dynamic"
	MainThread     = "main"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var errLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
var conf = config()
var threadType = threadControlChooser()

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
			Host:         "0.0.0.0",
			Port:         7303,
			RestHost:     "127.0.0.1:7301",
			ThreadType:   StaticThreads,
			Threads:      0,
			Obm:          Person{Name: "Max Mustermann", Address: "Musterweg 19", Telephone: "06647315794", Role: "obmann", Email: "obmann@mvl.at"},
			Kpm:          Person{Name: "Ursula Baum", Address: "Daham", Telephone: "06649146735", Role: "kapellmeister", Email: "kapellmeister@mvl.at"},
			King:         "Arnheim",
			Marches:      []string{"Arnheim", "Koline Koline", "Castaldo", "Attila", "Olympia", "Rázně Vpřed", "Florentinský", "Muziky, Muziky", "Slavnostní", "Fanfarovy"},
			Color:        "#134474FF",
			CalendarName: "Musikverein Leopoldsdorf",
			Name:         "Musikverein Leopoldsdorf/M.",
			Address:      "A-2285 Leopoldsdorf/M. Kempfendorf 2",
			ZVR:          "ZVR - Zahl: 091786949",
			Role:         "events",
			Timezone:     "Europe/Vienna"}
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

func threadControlChooser() (threadControl threadControl) {
	switch conf.ThreadType {
	case StaticThreads:
		threadControl = staticThreads
	case DynamicThreads:
		threadControl = dynamicThreads
	case MainThread:
		threadControl = mainThread
	default:
		threadControl = staticThreads
	}
	return
}

//Struct which holds the configuration.
type Configuration struct {
	Host         string   `json:"host"`
	Port         uint16   `json:"port"`
	RestHost     string   `json:"restHost"`
	ThreadType   string   `json:"threadType"`
	Threads      int      `json:"threads"`
	Obm          Person   `json:"obm"`
	Kpm          Person   `json:"kpm"`
	King         string   `json:"king"`
	Marches      []string `json:"marches"`
	Color        string   `json:"color"`
	CalendarName string   `json:"calendarName"`
	Name         string   `json:"name"`
	HomePage     string   `json:"homepage"`
	Address      string   `json:"address"`
	ZVR          string   `json:"zvr"`
	Role         string   `json:"role"`
	Timezone     string   `json:"timezone"`
}

type Person struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}
