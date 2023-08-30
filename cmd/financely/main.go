package main

import (
	"github.com/jollyboss123/finance-tracker/pkg/server"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stdout, "Financely ", log.LstdFlags)
	s := server.New(server.WithLogger(l))
	s.Init()
	s.Run()
}
