package main

import "net"
import "fmt"
import "code.google.com/p/govt/vt/vtclnt"

type Peer struct{
	Name string
	Addr string
	Conn net.Conn
	VTClnt *vtclnt.Clnt
	Msgchan chan *Syncmsg
	Server bool
}

func newpeer(addr string, conn net.Conn) *Peer{
	p := new(Peer)

	p.Addr = addr
	p.Conn = conn
	p.VTClnt = vtclnt.NewClnt(p.Conn)
	p.Msgchan = make(chan *Syncmsg, 32)

	/* XXX does the first tpeer need to be first? */
	vtconnect(p.VTClnt)

	/* add p to peerlist! */

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
