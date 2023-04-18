package Utilities

import (
	"github.com/spf13/viper"
)

func Load() {
	viper.AddConfigPath(".")
	viper.SetConfigName("Config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	viper.ReadInConfig()
}

// var koan = koanf.New(".")

// func Loadconfig() {
// 	f := file.Provider("Config.json")
// 	if err := koan.Load(f, json.Parser()); err != nil {
// 		log.Fatalf("error loading config: %v", err)
// 	}
// }

// func GetKaonf() *koanf.Koanf {
// 	return koan
// }
