package main

import (
	"github.com/jollyboss123/finance-tracker/pkg/server"
)

var version = "v0.1.0"

func main() {
	s := server.New(server.WithVersion(version))
	s.Init()
	s.Run()
}
