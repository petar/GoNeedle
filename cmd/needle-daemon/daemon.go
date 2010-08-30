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
	flagPing = flag.String("ping", ":62077", "Address where server listens for UDP ping updates")
)

func main() {
	flag.Parse()
	fmt.Fprintf(os.Stderr, 
		"Starting Needle Daemon, 2010 (C) Petar Maymounkov, " +
		"http://http://github.com/petar/GoNeedle\n")
	_,err := needle.MakeServer(*flagPing)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Listening for pings on %s ...\n", *flagPing)
	<-make(chan int)
}
