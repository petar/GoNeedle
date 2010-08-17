// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package needle

import (
	//"fmt"
	"testing"
	"time"
)

func startServer(t *testing.T) {
	_, err := MakeServer(":62077", ":62070")
	if err != nil {
		t.Fatalf("Starting server %s\n", err)
	}
}

func startListener(id string, t *testing.T) {
	_, err := MakeListener(id, ":34000", "127.0.0.1:62077")
	if err != nil {
		t.Fatalf("Starting listener %s\n", err)
	}
}

func startConnect(id string, t *testing.T) {
	_, err := Dial("localhost:62070", id, "haha")
	if err != nil {
		t.Fatalf("Connect: %s", err)
	}
}

func TestServer(t *testing.T) {
	startServer(t)
	time.Sleep(1e9)
	startListener("1", t)
	time.Sleep(1e9)
	startConnect("1", t)
	startConnect("2", t)
	<-make(chan int)
}
