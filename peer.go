package main

import "net"
import "fmt"
import "code.google.com/p/govt/vt/vtclnt"

type Peer struct{
	Addr string
	Conn net.Conn
	VTClnt *vtclnt.Clnt
}

func newpeer(addr string, conn net.Conn) *Peer{
	p := new(Peer)

	p.Addr = addr
	p.Conn = conn
	p.VTClnt = vtclnt.NewClnt(p.Conn)
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

/* We need to check if we have the peer before adding it. */
/* Otherwise we will end up with some crazy power set. */
/* Or possibly even an infinite flood of peer-adds. */
func havepeer(name string) bool{
	fmt.Println("XXX havepeer unimplemented!")
	return false
}
