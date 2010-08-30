// Copyright 2010 GoNeedle Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package needle

import (
	"container/list"
	"os"
)

type dialBook struct {
	dials map[string]*list.List
}

type dialTicket struct {
	when   int64           // When was the dial request issued
	notify chan<- os.Error
}

func makeDialBook() *dialBook {
	db := &dialBook{}
	db.Init()
	return db
}

func (db *dialBook) Init() {
	db.dials = make(map[string]*list.List)
}

func (db *dialBook) Add(id string, when int64, notify chan<- os.Error) {
	l, ok := db.dials[id]
	if !ok {
		l = list.New()
		db.dials[id] = l
	}
	l.PushBack(&dialTicket{when, notify})
}

// RETURNS the list of tickets (list values are *dialTicket) that have been expired
func (db *dialBook) Expire(now, ageLimit int64) *list.List {
	r := list.New()
	for _, l := range db.dials {
		for e := l.Front(); e != nil; e = e.Next() {
			t := e.Value.(*dialTicket)
			if now - t.when > ageLimit {
				l.Remove(e)
				r.PushBack(t)
			}
		}
	}
	return r
}

func (db *dialBook) GetDialsTicketsForId(id string) *list.List {
	r, ok := db.dials[id]
	if !ok {
		return list.New()
	}
	return r
}

func (db *dialBook) GetIds() []string {
	r := make([]string, len(db.dials))
		k := 0
	for id, _ := range db.dials {
		r[k] = id
		k++
	}
	return r
}
