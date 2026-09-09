package main

import (
	"github.com/amirhossein-karimi/todo/cmd"
	_ "github.com/amirhossein-karimi/todo/docs"
)

// @title Todo API
// @version 1.0
// @description Todo REST API built with Go and Gin.
// @host localhost:8080
// @BasePath /
func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
