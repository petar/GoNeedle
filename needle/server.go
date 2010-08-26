// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

// The Needle System is a simple general method for punching thru NATs and 
// Firewalls and establishing peer end-to-end reliable transport, using the
// help of an intermediate thin Needle UDP server.

package needle

import (
	"encoding/hex"
	"net"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/petar/GoNeedle/needle/proto"
	pb "goprotobuf.googlecode.com/hg/proto"
)

// TODO:
//   -- Use LLRB for expiration algorithm

type Server struct {
	conn *net.UDPConn	// Socket that receives UDP pings from the clients
	ids  map[string]*client  // Id-to-client map
	lk   sync.Mutex          // Lock for ids field
}

// client describes real-time information for a given client
type client struct {
	id           string
	lastSeen     int64
	addr         *net.UDPAddr

	// Ring source (or dial target) maps to the time when the ring (or dial)
	// request was communicated to the server
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

// o When pong-ing, attach dial/ring list
// o When ping-ing, read out the dial requests and update clients involved

func MakeServer(uaddr string, haddr string) (*Server, os.Error) {
	
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
		ids: make(map[string]*client),
	}
	go s.listenLoop()
	go s.expireLoop()

	return s, nil
}

// lookupAddrNeedsLock() returns the address of the given ID.
// This routine must be called inside a lock on s.lk
func (s *Server) lookupAddrNeedsLock(id string) *net.UDPAddr {
	c, ok := s.ids[id]
	if !ok {
		return nil
	}
	return c.addr
}

// expire() removes all client structures that have not been refreshed recently
func (s *Server) expire() {
	s.lk.Lock()
	defer s.lk.Unlock()

	now := time.Nanoseconds()
	for id, cl := range s.ids {
		if now - cl.lastSeen > Lifetime {
			s.ids[id] = nil, false
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

func (s *Server) expireLoop() {
	for {
		s.expire()
		time.Sleep(Lifetime)
	}
}

// makePongPacket() prepares a pong packet for the given id
func (s *Server) makePongPacket_NL(id string) []byte {
	c, ok := s.ids[id]
	if !ok {
		return nil
	}

	payload := &proto.PeerBound{}
	payload.Pong := &proto.Pong{}

	prep := make(map[string]string)

	for cid, _ := range c.rings {
		a := s.lookupAddrNeedsLock(cid)
		if a == nil {
			continue
		}
		prep[cid] := a.String()
	}

	for did, _ := range c.dials {
		a := s.lookupAddrNeedsLock(did)
		if a == nil {
			continue
		}
		prep[did] := a.String()
	}

	// XXX: Make sure that packet does not exceed allowed size
	payload.Pong.Punches := make([]*proto.PunchPoint, len(prep))
	k := 0
	for pId, pAddr := range prep {
		payload.Pong.Punches[k] = &pAddr
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
	defer s.lk.Unlock()

	taddr := s.lookupAddr_NL(id)
	if taddr == nil {
		return
	}
	packet := s.makePongPacket_NL(id)
	if packet == nil {
		return
	}
	n, err := s.conn.WriteToUDP(b, taddr)
}

func (s *Server) updateClient(id string, now int64, addr *net.UDPAddr) {
	s.lk.Lock()
	defer s.lk.Unlock()

	cl, ok := s.ids[id]
	if ok {
		cl.lastSeen = now
		cl.addr = addr
	} else {
		c := makeClient(id, now, addr)
		s.ids[id] = c
		?
	}
}

func (s *Server) poll() os.Error {

	// Read next UDP packet
	b := make([]byte, MaxIdLen + 32)
	n, addr, err := s.conn.ReadFromUDP(b)
	if err != nil {
		return err
	}
	
	// Decode packet contents
	payload := &proto.Ping{}
	err = pb.Unmarshal(b[0:n], payload)
	if err != nil || len(*payload.Id) > MaxIdLen {
		return err
	}

	// Make necessary updates
	s.updateClient(*payload.Id, time.Nanoseconds(), addr)

	return nil
}

func (s *Server) loop() {
	lastExpire := time.Nanoseconds()
	for {
		s.poll()
		now := time.Nanoseconds()
		if now - lastExpire > ExpirePeriod {
			s.expire(now)
			lastExpire = time.Nanoseconds()
		}
	}
}
