package main

import (
	"flag"
)

var configFile = flag.String("f", "configs/server.yaml", "set config file which viper will loading.")


// @title ioj-server
// @version 0.0.2
// @description API server for infinity-oj
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /api/v1
// @contact.name Wycer
// @contact.email wycers@gmail.com
func main() {
	flag.Parse()

	app, err := CreateApp(*configFile)
	if err != nil {
		panic(err)
	}

	if err := app.Start(); err != nil {
		panic(err)
	}

	app.AwaitSignal()
}
