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
	flagId     = flag.String("id", "", "ID of listener")
	flagBind   = flag.String("bind", "", "UDP address to bind to")
	flagServer = flag.String("server", "", "UDP address of Needle server")
)

func main() {
	flag.Parse()
	fmt.Fprintf(os.Stderr, 
		"Starting Needle Listen, 2010 (C) Petar Maymounkov, " +
		"http://http://github.com/petar/GoNeedle\n")

	_,err := needle.MakeListener(*flagId, *flagBind, *flagServer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Listening for connections on %s, using server %s\n",
		*flagBind, *flagServer)

	<-make(chan int)
}
