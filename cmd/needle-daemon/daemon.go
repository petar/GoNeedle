
package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/petar/GoNeedle/needle"
)

var (
	flagPing = flag.String("ping", ":62077", "Address where server listens for UDP ping updates")
	flagHTTP = flag.String("http", ":62070", "Address of HTTP API")
)

func main() {
	flag.Parse()
	fmt.Fprintf(os.Stderr, 
		"Starting Needle Daemon, 2010 (C) Petar Maymounkov, " +
		"http://http://github.com/petar/GoNeedle\n")
	_,err := needle.MakeServer(*flagPing, *flagHTTP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Listening for pings on %s, accepting HTTP API queries on %s\n",
		*flagPing, *flagHTTP)
	<-make(chan int)
}
