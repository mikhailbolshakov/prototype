package configuro

import (
	"encoding/json"
	"github.com/sherifabdlnaby/configuro"
	"log"
	"testing"
)

type Obj struct {
	v1 string
	v2 string
}

type Config struct {
	// Primitive Types
	Number      int
	NumberList  []int `config:"number_list"`
	Word        string
	AnotherWord string            `config:"another_word"`
	WordMap     map[string]string `config:"word_map"`
	// Nested Objects (Ptr and None)
	Database *Database
	Logger   Logger
	ComplexMap map[string]interface{}
}

//Database A sub-config struct
type Database struct {
	Hosts    []string
	Username string
	Password string
}

//Logger Another sub-config struct
type Logger struct {
	Level string
	Debug bool
}

func Test_Config(t *testing.T) {
	// Create Configuro Object
	Loader, err := configuro.NewConfig(
		configuro.WithLoadFromConfigFile("./config.yml", true),
		configuro.WithLoadDotEnv("./.env"))
	if err != nil {
		panic(err)
	}

	// Create our Config holding Struct
	config := &Config{Word: "default value in struct."}

	// Load Our Config.
	err = Loader.Load(config)
	if err != nil {
		panic(err)
	}

	j, _ := json.Marshal(config)

	log.Println(string(j))

}

type Redis struct {
	Port string
	Host string
	Password string
	Ttl int
}

type Db struct {
	Dbname string
	User string
	Password string
	Port string
	HostRw string `config:"host-rw"`
	HostRo string `config:"host-ro"`
}

type SrvCfg struct {
	Database *Db
	Grpc interface{}
}

type Config2 struct {
	Redis *Redis
	Services map[string]*SrvCfg
}

func Test_Config2(t *testing.T) {
	// Create Configuro Object
	Loader, err := configuro.NewConfig(
		configuro.WithLoadFromConfigFile("./config2.yml", true))
	if err != nil {
		panic(err)
	}

	// Create our Config holding Struct
	config := &Config2{}

	// Load Our Config.
	err = Loader.Load(config)
	if err != nil {
		panic(err)
	}

	j, _ := json.Marshal(config)

	log.Println(string(j))

}
