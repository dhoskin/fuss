package main

import "net"
import "io"
import "sync"
import "code.google.com/p/govt/vt"

/* XXX how to cleanup the mess if the connection goes down? */

func dialdispatch(myclnt net.Conn, srvaddr, peeraddr string){
	peer, e := net.Dial("tcp", peeraddr);
	if(e != nil){
		return;
	}
	mysrv, e := net.Dial("tcp", srvaddr);
	if(e != nil){
		return;
	}
	dispatch(myclnt, mysrv, peer);
}

func listendispatch(myclnt net.Conn, srvaddr, port string){
	l, e := net.Listen("tcp", port);
	if(e != nil){
		return;
	}
	for {
		peer, e := l.Accept();
		if(e != nil){
			return;
		}
		mysrv, e := net.Dial("tcp", srvaddr);
		if(e != nil){
			return;
		}
		go dispatch(myclnt, mysrv, peer);
	}
}

func copydiscrete(dst, src net.Conn, mtx *sync.Mutex) {
	var buf = make([]byte, vt.Maxblock * 2)
	n := 0
	for {
		pktsz := 0
		if n < 2 {
			x, e := io.ReadAtLeast(src, buf, 2)
			if e != nil {
				return;
			}
			n += x
		}

		/* Thanks a lot. */
		pktsz16, _ :=  vt.Gint16(buf)
		pktsz = int(pktsz16) + 2

		/* XXX what about packet containing maxblock? */
		if pktsz > vt.Maxblock{
			return
		}
		for n < pktsz {
			x, e := io.ReadAtLeast(src, buf, pktsz - n);
			if(e != nil){
				return;
			}
			n += x;
		}

		mtx.Lock()
		dst.Write(buf[0:pktsz])
		mtx.Unlock()

		/* XXX inefficient when multiple msg per pkt */
		if(n > pktsz){
			copy(buf, buf[n-pktsz:n]);
		}
		n = n - pktsz;
	}
}

func syncmsgproc(peer net.Conn, msgchan chan *Syncmsg, mtx *sync.Mutex) {
	buf := make([]byte, vt.Maxblock)

	for {
		msg := <- msgchan
		n := sync2wire(buf, msg)

		mtx.Lock()
		peer.Write(buf[:n])
		mtx.Unlock()
	}
}

func dispatch(myclnt, mysrv, peerconn net.Conn){
	var peer Peer
	_ = peer
	/* send Tpeer */
	/* recv Tpeer */

	msgchan := make(chan *Syncmsg, 32)
	var mtx sync.Mutex
	go copydiscrete(peerconn, myclnt, &mtx)
	go copydiscrete(peerconn, mysrv, &mtx)
	go syncmsgproc(peerconn, msgchan, &mtx)

	//go io.Copy(peerconn, myclnt);
	//go io.Copy(peerconn, mysrv);

	var buf = make([]byte, vt.Maxblock * 2);
	n := 0;
	for {
		pktsz := 0;

		if(n < 2){
			x, e := io.ReadAtLeast(peerconn, buf, 2);
			if(e != nil){
				return;
			}
			n += x;
		}

		/* Thanks a lot. */
		pktsz16, _ :=  vt.Gint16(buf);
		pktsz = int(pktsz16) + 2;

		/* XXX what about packet containing maxblock? */
		if(pktsz > vt.Maxblock){
			return;
		}
		for n < pktsz {
			x, e := io.ReadAtLeast(peerconn, buf, pktsz - n);
			if(e != nil){
				return;
			}
			n += x;
		}

		if(buf[2] >= 18){
			onsyncmsg(Self.Clnt, syncparse(buf));
		}else if((buf[2] % 2) == 1){
			myclnt.Write(buf[0:pktsz]);
		}else{
			mysrv.Write(buf[0:pktsz]);
		}

		/* XXX inefficient when multiple msg per pkt */
		if(n > pktsz){
			copy(buf, buf[n-pktsz:n]);
		}
		n = n - pktsz;
	}
}
