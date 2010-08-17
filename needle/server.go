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
	udp *net.UDPConn	// Socket that receives UDP pings from the clients
	ids map[string]*client  // Id-to-client map
	lk  sync.Mutex          // Lock for ids field
}

// client describes real-time information for a given client
type client struct {
	id       string
	lastSeen int64
	addr     *net.UDPAddr
}

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
	err = conn.SetReadTimeout(ExpirePeriod / 2)
	if err != nil {
		return nil, err
	}

	// Start server
	s := &Server{
		udp: conn,
		ids: make(map[string]*client),
	}
	_, err = makeQueryAPI(
		haddr, 
		func(q string) (string, os.Error) { return s.Query(q) }, 
		MaxFD,
	)
	if err != nil {
		conn.Close()
		return nil, err
	}
	go s.loop()

	return s, nil
}

// expire removes all client structures that have not been refreshed recently
func (s *Server) expire(now int64) {
	s.lk.Lock()
	defer s.lk.Unlock()

	for id, cl := range s.ids {
		if now - cl.lastSeen > ClientFreshness {
			s.ids[id] = nil, false
		}
	}
}

func (s *Server) updateClient(id string, now int64, addr *net.UDPAddr) {
	s.lk.Lock()
	defer s.lk.Unlock()

	cl, ok := s.ids[id]
	if ok {
		cl.lastSeen = now
		cl.addr = addr
	} else {
		s.ids[id] = &client{
			id:       id,
			lastSeen: now,
			addr:     addr,
		}
	}
}

func (s *Server) poll() os.Error {

	// Read next UDP packet
	b := make([]byte, MaxIdLen + 32)
	n, addr, err := s.udp.ReadFromUDP(b)
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

func int64ToString(i int64) string {
	pbuf := pb.NewBuffer(nil)
	pbuf.EncodeFixed64(uint64(i))
	return hex.EncodeToString(pbuf.Bytes())
}

func stringToInt64(s string) (int64, os.Error) {
	buf, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	pbuf := pb.NewBuffer(buf)
	ui, err := pbuf.DecodeFixed64()
	if err != nil {
		return 0, err
	}
	return int64(ui), nil
}

func (s *Server) Query(query string) (string, os.Error) {
	qid := strings.TrimSpace(query)
	if len(qid) > MaxIdLen {
		return "", ErrIdLen
	}
	s.lk.Lock()
	cl, ok := s.ids[qid]
	result := "!no-entry"
	if ok {
		result = cl.addr.String()
	}
	s.lk.Unlock()
	return result, nil
}
