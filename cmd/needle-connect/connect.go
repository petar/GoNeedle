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
	flagRemoteId = flag.String("remote-id", "", "Remote ID to connect to")
)

func main() {
	needle.InstallCtrlCPanic()
	flag.Parse()
	fmt.Fprintf(os.Stderr, 
		"Starting Needle Connect, 2010 (C) Petar Maymounkov, " +
		"http://github.com/petar/GoNeedle\n")

	peer, err := needle.MakePeer(*flagLocalId, *flagBind, *flagServer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem starting peer: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Dialing %s, using server %s\n", *flagRemoteId, *flagServer)
	peer.Dial(*flagRemoteId)

	<-make(chan int)
}
