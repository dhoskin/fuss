package main

import "net"
import "fmt"
import "code.google.com/p/govt/vt/vtclnt"

type Peer struct{
	Name string
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

/* Uh oh, this isn't as obvious as it seemed... */
/* Seems like we need another round of UML diagrams! */
/* maybe general Manager rather than PeerManager? */
/* biggest problem seems to be the name... */
/* syncmanager? */
/* OO design: func newPeerMgr() */
/* sets up chans and runs go peermgr(chan, chan, ...) */
/* then returns PeerMgr object with addpeer etc methods */
func peermgr() {
/*
Should this function subsume syncmgr or whatever?
What about the tag manager?
where should the tag manager even go?
should it just be a list or something with a mutex?

Score added successfully:
- update master tag list
- send score to all peers

Peer added successfully:
- send it the whole tag list
- send it a Tpeer for each current peer
- send its Tpeer to each current peer

Insane idea: "are you also?" message
heavy-handed way to avoid peer duplication
can't just do host:port if we're
running multiple daemons on one host for debug...
(or just bind to IP addr rather than 0.0.0.0?)
*/
}

/* We need to check if we have the peer before adding it. */
/* Otherwise we will end up with some crazy power set. */
/* Or possibly even an infinite flood of peer-adds. */
func havepeer(name string) bool{
	fmt.Println("XXX havepeer unimplemented!")
	return false
}
