// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

// The Needle System is a simple general method for punching thru NATs and 
// Firewalls and establishing peer end-to-end reliable transport, using the
// help of an intermediate thin Needle UDP server.

package needle

import (
	"net"
	"os"
	"sync"
	"time"
	"github.com/petar/GoNeedle/needle/proto"
	pb "goprotobuf.googlecode.com/hg/proto"
)

// TODO:
//   -- Use LLRB for expiration algorithm

type Server struct {
	conn  *net.UDPConn	  // Socket that receives UDP pings from the clients
	peers map[string]*client  // Id-to-client map
	lk    sync.Mutex          // Lock for peers field
}

// client describes real-time information for a given client
type client struct {
	id           string
	lastSeen     int64
	addr         *net.UDPAddr

	// @dials maps the ID of the peer, that this peer has requested to connect to, to
	// the time when the request was communicated to the server.
	// @rings maps the ID of the peer, requesting to connect to this peer, to 
	// the time when the request was communicated to the server.
	dials, rings map[string]int64
}

func makeClient(id string, lastSeen int64, addr *net.UDPAddr) *client {
	return &client{
		id:       id,
		lastSeen: lastSeen,
		addr:     addr,
		dials:    make(map[string]int64),
		rings:    make(map[string]int64),
	}
}

// o When ping-ing, read out the dial requests and update clients involved

func MakeServer(uaddr string) (*Server, os.Error) {
	
	// Resolve UDP address
	uaddr2, err := net.ResolveUDPAddr(uaddr)
	if err != nil {
		return nil, err
	}

	// Bind and setup UDP connection
	conn, err := net.ListenUDP("udp", uaddr2)
	if err != nil {
		return nil, err
	}

	// Start server
	s := &Server{
		conn: conn,
		peers: make(map[string]*client),
	}
	go s.listenLoop()
	go s.expireLoop()
	go s.pongLoop()

	return s, nil
}

// lookupAddrNeedsLock() returns the address of the given ID.
func (s *Server) lookupAddr_NL(id string) *net.UDPAddr {
	c, ok := s.peers[id]
	if !ok {
		return nil
	}
	return c.addr
}

// expire() removes all client structures, dial and ring records that 
// have not been refreshed recently
func (s *Server) expire() {
	s.lk.Lock()
	defer s.lk.Unlock()

	now := time.Nanoseconds()
	for id, cl := range s.peers {
		if now - cl.lastSeen > Lifetime {
			s.peers[id] = nil, false
			continue
		}
		for id2, l2 := range cl.dials {
			if now - l2 > Lifetime {
				cl.dials[id2] = 0, false
			}
		}
		for id2, l2 := range cl.rings {
			if now - l2 > Lifetime {
				cl.rings[id2] = 0, false
			}
		}
	}
}

// expireLoop() runs in its own go-routine and periodically cleans up stale records
func (s *Server) expireLoop() {
	for {
		s.expire()
		time.Sleep(Lifetime)
	}
}

// makePongPacket() prepares a pong packet for the given id
func (s *Server) makePongPacket_NL(id string) []byte {
	c, ok := s.peers[id]
	if !ok {
		return nil
	}

	payload := &proto.PeerBound{}
	payload.Pong = &proto.Pong{}

	prep := make(map[string]string)

	for cid, _ := range c.rings {
		a := s.lookupAddr_NL(cid)
		if a == nil {
			continue
		}
		prep[cid] = a.String()
	}

	for did, _ := range c.dials {
		a := s.lookupAddr_NL(did)
		if a == nil {
			continue
		}
		prep[did] = a.String()
	}

	// XXX: Make sure that packet does not exceed allowed size
	payload.Pong.Punches = make([]*proto.PunchPoint, len(prep))
	k := 0
	for pId, pAddr := range prep {
		payload.Pong.Punches[k].Id = &pId
		payload.Pong.Punches[k].Address = &pAddr
		k++
	}

	packet, err := pb.Marshal(payload)
	if err != nil {
		return nil
	}
	return packet
}

// pong() sends a packet to the desired node
func (s *Server) pong(id string) {
	s.lk.Lock()
	taddr := s.lookupAddr_NL(id)
	if taddr == nil {
		s.lk.Unlock()
		return
	}
	packet := s.makePongPacket_NL(id)
	if packet == nil {
		s.lk.Unlock()
		return
	}
	s.lk.Unlock()
	s.conn.WriteToUDP(packet, taddr)
}

// pongLoop() sends pongs to all its clients at regular intervals
func (s *Server) pongLoop() {
	lastPong := int64(0)
	for {
		sleep := PongPeriod - (time.Nanoseconds() - lastPong)
		if sleep > 0 {
			time.Sleep(sleep)
		}

		// Pong all peers
		lastPong = time.Nanoseconds()
		s.lk.Lock()
		pp := make([]string, len(s.peers))
		k := 0
		for id, _ := range s.peers {
			pp[k] = id
			k++
		}
		s.lk.Unlock()
		for _, id := range pp {
			s.pong(id)
		}
	}
}

// processPing() updates data structures in light of the data from a ping packet
func (s *Server) processPing(id string, now int64, addr *net.UDPAddr, dialing []string) {
	s.lk.Lock()
	defer s.lk.Unlock()

	c, ok := s.peers[id]
	if ok {
		c.lastSeen = now
		c.addr = addr
	} else {
		c = makeClient(id, now, addr)
		s.peers[id] = c
	}

	for _, d := range dialing {
		c.dials[d] = now
		h, ok := s.peers[d]
		if ok {
			h.rings[id] = now
		}
	}
}

// receivePing() reads a single ping packet, makes necessary adjustments to the
// realtime data structures and issues some pongs if necessary
func (s *Server) receivePing() {

	// Read next UDP packet
	b := make([]byte, MaxPacketSize + 1)
	n, addr, err := s.conn.ReadFromUDP(b)
	Logf("Pack recv'd\n")
	if err != nil || n >= MaxPacketSize + 1 {
		return
	}
	
	// Decode packet contents
	ping := &proto.Ping{}
	err = pb.Unmarshal(b[0:n], ping)
	if err != nil || len(*ping.Id) > MaxIdLen {
		return
	}

	// Make necessary updates
	s.processPing(*ping.Id, time.Nanoseconds(), addr, ping.Dialing)
}

// listenLoop() runs in its own go-routine and repeatedly listens for incoming pings
func (s *Server) listenLoop() {
	for {
		s.receivePing()
	}
}
