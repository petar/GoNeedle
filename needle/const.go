// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package needle

import (
	"os"
)

const (
	MaxIdLen        = 64       // Maximum number of characters in a node ID
	ExpirePeriod    = 30e9     // Run expiration loop every 30 secs
	ClientFreshness = 5e9      // Expire clients who haven't pinged in the past 5 secs
	ListenerRefresh = 3e9      // How often the client sends a ping, every 3 secs
	MaxFD           = 200      // Maximum number of concurrent FDs used by HTTP API
	MaxPacketSize   = 32*1024  // Maximum packet size for payload packets
)

var (
	ErrIdLen        = os.NewError("Id exceeds max len")
)
