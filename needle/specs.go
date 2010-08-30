// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package needle

import (
	"os"
)

const (
	PingPeriod      = 3e9            // How often a peer sends a ping, every 3 secs
	PongPeriod      = 3e9            // How often every peer is being ponged, every 3 sec
	DialTimeout     = 20e9           // How long to wait until dial succeeds, 20 secs
	MaxPacketSize   = 64*1024        // Maximum packet size for ping/pong/cargo messages
	MaxIdLen        = 64             // Maximum number of characters in a node ID
	Lifetime        = 2*PingPeriod   // Lifetime of presence markers, dial and ring requests
)

var (
	ErrTimeout      = os.NewError("timeout")
	//ErrIdLen        = os.NewError("Id exceeds max len")
)
