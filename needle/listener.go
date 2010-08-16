// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package needle

import (
	"net"
	"os"
	"sync"
	"time"
	"tonika/needle/proto"
	pb "goprotobuf.googlecode.com/hg/proto"
)

type Listener struct {
	serverAddr *net.UDPAddr
	pingPacket []byte         // Prepared ping packet
	conn       *net.UDPConn   // UDP connection for pings to needle server and all peer traffic
	connLk     sync.Mutex     // Lock for conn
}

func MakeListener(id int64, bindAddr, serverAddr string) (*Listener, os.Error) {

	// Resolve UDP addresses
	saddr, err := net.ResolveUDPAddr(serverAddr)
	if err != nil {
		return nil, err
	}

	baddr, err := net.ResolveUDPAddr(bindAddr)
	if err != nil {
		return nil, err
	}

	// Bind UDP port
	conn, err := net.ListenUDP("udp", baddr)
	if err != nil {
		return nil, err
	}

	// Prepare ping packet
	payload := &proto.Ping{Id:&id}
	pingPacket, err := pb.Marshal(payload) 
	if err != nil {
		return nil, err
	}

	// All init went OK, start client
	l := &Listener{
		serverAddr: saddr,
		pingPacket: pingPacket,
		conn:       conn,
	}
	go l.pingLoop()

	return l, nil
}

func (l *Listener) pingLoop() {
	for {
		l.connLk.Lock()
		l.conn.WriteToUDP(l.pingPacket, l.serverAddr)
		l.connLk.Unlock()
		time.Sleep(ListenerRefresh)
	}
}
