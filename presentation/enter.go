package presentation

import (
	// "proposal-template/pkg/logger"
	httpserver "proposal-template/presentation/http"
	"sync"

	"github.com/golobby/container/v3"
)

type server struct {
	httpServer httpserver.HTTPServer
	//grpcServer...
	//...
}

func NewServer() *server {

	var hs *httpserver.HTTPServer
	err := container.Resolve(&hs)
	if err != nil {
		panic(err)
	}

	return &server{
		httpServer: *hs,
	}
}

// Run starts all the servers in a separate goroutine and waits for any one of them to return an error.
func (s *server) Run() error {
	var wg sync.WaitGroup


	errChan := make(chan error, 1)

	wg.Add(1) 
	//If there are more than http server, this "1" value should be increased, for example
	// if we added a grpc server, this value should be "2"
	go func ()  {
		defer wg.Done()
		if err := s.httpServer.Start(); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()
	
	//== if there are more than a server running, we add a goroutine like above
	//== example
	// go func ()  {
	// 	defer wg.Done()
	// 	if err := s.grpcServer.Start(); err != nil {
	// 		errChan <- err
	// 	} else {
	// 		errChan <- nil
	// 	}
	// }

	wg.Wait()
	err := <-errChan
	return err
}