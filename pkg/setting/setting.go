package setting

import (
	"log"

	"github.com/go-ini/ini"
)

type Wagent struct {
	Port int
	Host string
}

type Guarder struct {
	Port int
	Host string
}

var WagentSetting = &Wagent{}

var GuarderSetting = &Guarder{}

var cfg *ini.File

func Setup() {
	var err error
	cfg, err = ini.Load("conf/wagent.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/wagent.ini': %v", err)
	}
	mapTo("wagent", WagentSetting)
	mapTo("guarder", GuarderSetting)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
