// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package needle

import (
	"net"
	"os"
	"sync"
	"time"
	"github.com/petar/GoNeedle/needle/proto"
	pb "goprotobuf.googlecode.com/hg/proto"
)

type Conn struct {
}

func Dial(needleServerAddr, targetId, targetPort string) (*Conn, os.Error) {
}
