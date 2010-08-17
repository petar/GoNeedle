// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/petar/GoNeedle/needle"
)

var (
	flagServer = flag.String("server", "", "Address of Needle server HTTP API")
	flagId     = flag.String("id", "", "Target ID to connect to")
	flagPort   = flag.String("port", "", "Target port to connect to")
)

func main() {
	flag.Parse()
	fmt.Fprintf(os.Stderr, 
		"Starting Needle Connect, 2010 (C) Petar Maymounkov, " +
		"http://http://github.com/petar/GoNeedle\n")

	_,err := needle.Dial(*flagServer, *flagId, *flagPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Dialing %s:%s, using server %s\n",
		*flagId, *flagPort, *flagServer)

	<-make(chan int)
}
