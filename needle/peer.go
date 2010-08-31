// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package needle

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
	"github.com/petar/GoNeedle/needle/proto"
	pb "goprotobuf.googlecode.com/hg/proto"
)

type Peer struct {
	id         string        // ID of this peer
	serverAddr *net.UDPAddr  // UDP address of needle server
	conn       *net.UDPConn  // UDP connection for all purposes
	dials      dialBook      // Structure of outstanding dial requests
	lk         sync.Mutex
}

func MakePeer(id, bindAddr, serverAddr string) (*Peer, os.Error) {

	saddr, err := net.ResolveUDPAddr(serverAddr)
	if err != nil {
		return nil, err
	}

	baddr, err := net.ResolveUDPAddr(bindAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", baddr)
	if err != nil {
		return nil, err
	}

	p := &Peer{
		id:         id,
		serverAddr: saddr,
		conn:       conn,
	}
	p.dials.Init()

	go p.pingLoop()
	go p.listenLoop()
	go p.expireDialsLoop()

	return p, nil
}

//func (p *Peer) Dial(dialToId string) (Conn, os.Error) {
func (p *Peer) Dial(dialToId string) {
	c := make(chan os.Error)
	p.lk.Lock()
	p.dials.Add(dialToId, time.Nanoseconds(), c)
	p.lk.Unlock()
	p.ping()
	<-c

	// XXX: to be continued
	return
}

// makePingPacket() creates a current ping packet for sending to needle server
func (p *Peer) makePingPacket() []byte {
	p.lk.Lock()
	defer p.lk.Unlock()

	payload := &proto.Ping{
		Id:      &p.id,
		Dialing: p.dials.GetIds(),
	}
	buf, err := pb.Marshal(payload)
	if err != nil {
		return nil
	}
	return buf
}

// ping() sends a ping to the needle server
func (p *Peer) ping() {
	packet := p.makePingPacket()
	if packet != nil {
		p.conn.WriteToUDP(packet, p.serverAddr)
		Logf("Ping OK\n")
	}
}

func (p *Peer) pingLoop() {
	for {
		p.ping()
		time.Sleep(PingPeriod)
	}
}

func (p *Peer) expireDialsLoop() {
	for {
		p.lk.Lock()
		expired := p.dials.Expire(time.Nanoseconds(), DialTimeout)
		p.lk.Unlock()
		for e := expired.Front(); e != nil; e = e.Next() {
			t := e.Value.(*dialTicket)
			t.notify <- ErrTimeout
			close(t.notify)
		}
		time.Sleep(PingPeriod)
	}
}

// listenLoop() receives pongs from the needle server and cargo packets from peers
func (p *Peer) listenLoop() {
	for {
		buf := make([]byte, MaxPacketSize)
		n, addr, err := p.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		Logf("Packet from %s/%d\n", addr.String(), n)

		payload := &proto.PeerBound{}
		err = pb.Unmarshal(buf[0:n], payload)
		if err != nil {
			continue
		}

		switch {
		case payload.Pong != nil:
			for _, m := range payload.Pong.Punches {
				// Resolve remote peer
				addr, err := net.ResolveUDPAddr(*m.Address)
				if err != nil {
					break
				}

				// Prepare cargo packet
				payload := &proto.PeerBound{}
				payload.Cargo = &proto.Cargo{}
				payload.Cargo.OriginId = &p.id
				packet, err := pb.Marshal(payload)
				if err != nil {
					break
				}
				// Start sending empty cargos
				go func() {
					for {
						p.conn.WriteToUDP(packet, addr)
						time.Sleep(1e9)
					}
				}()
			}
		case payload.Cargo != nil:
			fmt.Printf("  cargo from %s\n", *payload.Cargo.OriginId)
		}

	}
}
