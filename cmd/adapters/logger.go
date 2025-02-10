package adapters

import (
	"proposal-template/pkg/logger"

	"github.com/golobby/container/v3"
)

func IoCLogger(){
	container.Singleton(func() logger.ILogger {
		return logger.NewLogger("debug")
	})
}