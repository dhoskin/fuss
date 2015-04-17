package main

import "net"
import "io"
import "bytes"
import "code.google.com/p/govt/vt"

type abuf struct {
	bytes.Buffer
	data []byte
}

func newabuf() abuf {
	var a abuf
	a.data = make([]byte, vt.Maxblock)
	return a
}

func (b abuf) Fill(r io.Reader) error {
	n, e := r.Read(b.data)
	if e != nil {
		return e
	}
	_, e = b.Write(b.data[0:n])
	return e
}

func dial(peername, peeraddr string){
	peerconn, e := net.Dial("tcp", peeraddr)
	if e != nil {
		return
	}

	peer := newpeer(peername, peeraddr, peerconn)
	_ = peer
	// announce my own peer structure
	talk(peer)
}

func listen(port string){
	l, e := net.Listen("tcp", port)
	if e != nil {
		return
	}

	for {
		peerconn, e := l.Accept()
		if e != nil {
			return
		}

		// recv Tpeer
		// send Rpeer?
		go talk(peer)
	}
}


/* Will the other cases actually go here? */
/* Rread */
func send(peer *Peer){
	var buf = peer.Tbuf.data
	for {
		select {
		case msg := <- peer.Msgchan:
			n := sync2wire(buf, msg)
			peer.Conn.Write(buf[:n])
		}
	}
}

func recv(peer *Peer){
	var buf = peer.Rbuf
	var conn = peer.Conn

	for {
		for buf.Len() < 2 {
			buf.Fill(conn)
		}

		pktsz16, _ := vt.Gint16(buf.Bytes())
		pktsz := int(pktsz16) + 2

		for buf.Len() < pktsz {
			buf.Fill(conn)
		}

		_ = pktsz
		handlemsg(peer, buf.Next(pktsz));
	}
}

func talk(peer *Peer){
	go send(peer)
	go recv(peer)
}
