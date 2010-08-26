// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package needle

import (
	"os"
)

const (
	PingPeriod      = 3e9      // How often a peer sends a ping, every 3 secs
	DialTimeout     = 20e9     // How long to wait until dial succeeds, 20 secs
	MaxPacketSize   = 32*1024  // Maximum packet size for ping/pong/cargo messages

	Lifetime  = 3*PingPeriod   // Lifetime of presence markers, dial and ring requests
	//MaxIdLen        = 64     // Maximum number of characters in a node ID
	//ExpirePeriod    = 30e9   // Run expiration loop every 30 secs
	//ClientFreshness = 5e9    // Expire clients who haven't pinged in the past 5 secs
)

var (
	ErrTimeout      = os.NewError("timeout")
	//ErrIdLen        = os.NewError("Id exceeds max len")
)
