package Utilities

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
)

// func Load(configfile *string)  {
// 	if configfile != nil{
// 		viper.SetConfigFile(*configfile)
// 	}
// 	fmt.Sprintln("Loading Config data...")
// 	viper.AddConfigPath("./config")
// 	viper.AddConfigPath("../config")
// 	err := viper.ReadConfig()

// 	if err != nil{
// 		panic(fmt.Errorf("Fatal error while reading the config file : %s",err.Error()))
// 	}else{
// 		fmt.Sprintln("Config file:", viper.ConfigFileUsed())
// 	}

// }
var koan = koanf.New(".")

func Loadconfig() {
	f := file.Provider("Config.json")
	if err := koan.Load(f, json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
}

func GetKaonf() *koanf.Koanf {
	return koan
}
