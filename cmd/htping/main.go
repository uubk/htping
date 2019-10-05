package main

import (
	"flag"
	"github.com/uubk/htping/pkg"
)

func main() {
	obj := pkg.HTPing{}
	obj.RegisterArgs()
	flag.Parse()
	obj.Listen()
}
