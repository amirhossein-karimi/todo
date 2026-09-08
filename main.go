package main

import "github.com/amirhossein-karimi/todo/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
