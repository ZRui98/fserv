package main

import (
	"flag"
	"os"

	"github.com/zrui98/fserv/config"
	"github.com/zrui98/fserv/server"
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
	s.ListenAndServe()
}
