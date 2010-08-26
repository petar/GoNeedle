
package needle

// punchBook{} is a structure for keeping track of which
// remote peer connections are currently opened (punched) and
// what their addresses and last access times are.
type punchBook struct {
	ids map[string]*punchTicket
}

type punchTicket struct {
	addr                  *net.UDPAddr
	lastSend, lastReceive int64           // Times of last send and receive
}

func makePunchBook() *punchBook {
	pb := &punchBook{}
	pb.Init()
	return pb
}

func (pb *punchBook) Init() {
	pb.ids := make(map[string]*punchTicket)
}

func (pb *punchBook) Get() *punchTicket {
	?
}

func (pb *punchBook) Expire(now, ageLimit int64) []*punchTicket {
	?
}
