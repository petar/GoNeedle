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
	flagBind     = flag.String("bind", "", "Local UDP address to bind to")
	flagServer   = flag.String("server", "", "Address of Needle server")
	flagLocalId  = flag.String("local-id", "", "Local ID")
)

func main() {
	flag.Parse()
	fmt.Fprintf(os.Stderr, 
		"Starting Needle Listen, 2010 (C) Petar Maymounkov, " +
		"http://http://github.com/petar/GoNeedle\n")

	_, err := needle.MakePeer(*flagLocalId, *flagBind, *flagServer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem starting peer: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Listening for connections on %s, using server %s\n",
		*flagBind, *flagServer)

	<-make(chan int)
}
