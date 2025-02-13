package main

import (
	"fmt"
	"proposal-template/cmd/adapters"
	"proposal-template/presentation"
)

func init() {
	fmt.Println("Initializing IoC container...") // Debugging
	adapters.IoCConfig()
	adapters.IoCLogger()
	adapters.IoCDatabase()
	adapters.IoCRepositories()
	adapters.IoCBiz()
	adapters.IoCServer()
	fmt.Println("IoC container initialized.") 
}
func main() { 
	app := presentation.NewServer()
	app.Run()
}