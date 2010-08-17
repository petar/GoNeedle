// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package needle

import (
	"fmt"
	"net"
	"os"
	//"sync"
	//"time"
	//"github.com/petar/GoNeedle/needle/proto"
	//pb "goprotobuf.googlecode.com/hg/proto"
)

type Conn struct {
	udp *net.UDPConn
}

func Dial(needleServerAddr, targetId, targetPort string) (*Conn, os.Error) {

	// Query for target node's UDP address
	answer, err := fetchAPI(needleServerAddr, targetId)
	if err != nil {
		return nil, err
	}
	fmt.Printf("a=%s\n", answer)

/*
	// Bind/dial target node's UDP address
	conn, err := net.DialUDP("udp", nil, saddr)
	if err != nil {
		return nil, err
	}

	// Prepare packet
	connectPayload := &proto.Connect{Port:&targetPort}
	packet, err := pb.Marshal(connectPayload)
	if err != nil {
		return nil, err
	}
*/
	// XXX ?
	return nil, nil
}
