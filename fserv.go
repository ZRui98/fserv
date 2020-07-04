package main

import (
	"os"
	"flag"
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
	s := server.New()
	s.ListenAndServe()
}
