package adapters

import (
	utils "proposal-template/pkg/utils/config"

	"github.com/golobby/container/v3"
)

func IoCConfig() {
	container.Singleton(func() utils.AppConfig {
		cfg, err := utils.LoadConfig()
		if err != nil {
			panic(err)
		}
		// fmt.Println("Config successfully registered in IoC:", cfg.Httpserver.Port)
		return *cfg
	})
}