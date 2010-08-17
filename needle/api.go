// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package needle

import (
	//"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"github.com/petar/GoHTTP/http"
)

type queryAPI struct {
	l         net.Listener
	fdlim     http.FDLimiter
	queryFunc queryFunc
}

type queryFunc func(string) (string, os.Error)

func makeQueryAPI(listenAddr string, qf queryFunc, maxfd int) (*queryAPI, os.Error) {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	qa := &queryAPI{
		l:         l,
		queryFunc: qf,
	}
	qa.fdlim.Init(maxfd)
	go qa.acceptLoop()
	return qa, nil
}

func (qa *queryAPI) acceptLoop() {
	for {
		qa.fdlim.Lock()
		conn, err := qa.l.Accept()
		if err != nil {
			qa.fdlim.Unlock()
			continue
		}
		conn = http.NewConnRunOnClose(conn, func() { qa.fdlim.Unlock() })
		go qa.serveAndCloseConn(conn)
	}
}

func (qa *queryAPI) serveAndCloseConn(conn net.Conn) {

	defer conn.Close()

	sc := http.NewServerConn(conn, nil)
	defer sc.Close()

	req, err := sc.Read()
	if err != nil {
		return
	}

	query := req.URL.RawQuery
	args, err := http.ParseQuery(query)
	if err != nil {
		return
	}

	q, ok := args["q"]
	if !ok || len(q) != 1 {
		return
	}

	r, err := qa.queryFunc(q[0])
	if err != nil {
		return
	}

	resp := buildResp(r)
	sc.Write(resp)
}

func fetchAPI(serverAddr, query string) (string, os.Error) {

	conn, err := net.Dial("tcp", "", serverAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	cc := http.NewClientConn(conn, nil)
	defer cc.Close()

	url, err := http.ParseURL("http://" + serverAddr + "/?q=" + http.URLEscape(query))
	if err != nil {
		return "", err
	}
	req := &http.Request{
		Method: "GET",
		URL: url, 
		Proto: "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		//Close: true, XXX possible HTTP bug. Does not work when this is uncommented.
		Host: serverAddr,
		UserAgent: "GoNeedle-ClientConnect",
	}

	//d, err := http.DumpRequest(req, true)
	//fmt.Printf("REQ:\n%s\n", string(d))

	err = cc.Write(req)
	if err != nil {
		return "", err
	}

	resp, err := cc.Read()
	if resp == nil {
		return "", err
	}
	if resp.Body == nil {
		return "", nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

var (
	respOK = &http.Response{
		Status: "OK",
		StatusCode: 200,
		Proto: "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		RequestMethod: "GET",
		Close: true,
	}
	lkRespOK sync.Mutex
)

func buildResp(html string) *http.Response {
	lkRespOK.Lock()
	defer lkRespOK.Unlock()
	resp, err := http.DupResp(respOK)
	if err != nil {
		panic("needle, DupResp")
	}
	resp.Body = http.StringToBody(html)
	resp.ContentLength = int64(len(html))
	return resp
}
