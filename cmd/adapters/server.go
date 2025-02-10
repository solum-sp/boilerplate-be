package adapters

import (

	"proposal-template/pkg/logger"
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
		
		
		server := httpserver.NewHTTPServer(
			httpserver.WithLogger(logger),
		)
	
		// fmt.Println("HTTPServer successfully registered in IoC") ==> Debugging
		return server
	})
}