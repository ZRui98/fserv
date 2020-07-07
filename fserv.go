package main

import (
	"os"
	"flag"
	"github.com/zrui98/fserv/server"
	"github.com/zrui98/fserv/config"
)

func usage() {
    flag.PrintDefaults()
    os.Exit(2)
}

func init() {
    flag.Usage = usage
    flag.Set("logtostderr", "true")
    flag.Set("v", "2")
    flag.Parse()
}
func main() {
	c := config.LoadConfig()
	s := server.New(c)
	s.SetupRoutes()
	s.ListenAndServe()
}
