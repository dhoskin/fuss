package main

import "net"
import "io"
import "code.google.com/p/govt/vt"

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

func dispatch(myclnt, mysrv, peer net.Conn){

	go io.Copy(peer, myclnt);
	go io.Copy(peer, mysrv);

	var buf = make([]byte, vt.Maxblock * 2);
	n := 0;
	for {
		pktsz := 0;

		if(n < 2){
			x, e := io.ReadAtLeast(peer, buf, 2);
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
		if(n < pktsz){
			x, e := io.ReadAtLeast(peer, buf, pktsz - n);
			if(e != nil){
				return;
			}
			n += x;
		}

		if(buf[2] >= 18){
			msg := syncparse(buf);
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