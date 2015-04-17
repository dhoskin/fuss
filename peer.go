package main

import "net"
import "fmt"
//import "code.google.com/p/govt/vt/vtclnt"

type Peer struct{
	sync.Mutex
	Name string
	Addr string
	Conn net.Conn
	Msgchan, Reply chan *Syncmsg
	Tbuf, Rbuf abuf
	Server bool
}

func newpeer(name, addr string, conn net.Conn) *Peer{
	p := new(Peer)

	p.Name = name
	p.Addr = addr
	p.Conn = conn
	p.Msgchan = make(chan *Syncmsg, 32)
	p.Reply = make(chan *Syncmsg, 32)
	p.Tbuf = newabuf()
	p.Rbuf = newabuf()

	return p
}

func addpeer(name string) *Peer{
	/* if havepeer() return; */
	/* newpeer() */
	/* dialdispatch() */
	/* send our scores to peer */
	/* send our peers to peer */
	fmt.Println("XXX addpeer unimplemented!")
	return nil;
}
