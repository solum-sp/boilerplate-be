package adapters

import (
	"proposal-template/pkg/logger"
	utils "proposal-template/pkg/utils/config"
	"proposal-template/presentation/http"

	"github.com/golobby/container/v3"
)



func IoCServer() {
	container.Singleton(func() *httpserver.HTTPServer {
		var (
			logger  logger.ILogger
		)

		err := container.Resolve(&logger)
		if err != nil {
			panic(err)
		}
		
		var appConfig utils.AppConfig
		container.Resolve(&appConfig)
		server := httpserver.NewHTTPServer(
			httpserver.WithLogger(logger),
			httpserver.WithConfig(appConfig.Httpserver),
		)
		
		// fmt.Println("HTTPServer successfully registered in IoC") ==> Debugging
		return server
	})
}