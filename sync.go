package main

import "fmt"
import "code.google.com/p/govt/vt"
import "code.google.com/p/govt/vt/vtclnt"

const(
	Thash = 20
	Rhash = 21
	Tpeer = 22
	Rpeer = 23
)

// sz[2] Thash tag[s] btype[1] score[20]
// sz[2] Tpeer name[s]

type Syncmsg struct {
	Id byte
	Btype byte
	Score vt.Score
	Name string
	Peer *Peer
}

func (s *Syncmsg) String() string {
	switch s.Id {
	case Thash:
		return fmt.Sprintf("Thash tag %v score %v type %v",
			s.Name, s.Score, s.Btype);
	case Tpeer:
		return fmt.Sprintf("Tpeer %s", s.Name);
	default:
		return fmt.Sprintf("unknown syncmsg %d", s.Id);
	}
}

func syncparse(buf []byte) *Syncmsg {
	var msg *Syncmsg

	if(len(buf) < 3){
		return nil;
	}

	n, buf := vt.Gint16(buf);
	if(len(buf) != int(n)){
		return nil;
	}

	msg.Id, buf = vt.Gint8(buf);
	switch(msg.Id){
	case Thash:
		msg.Name, buf = vt.Gstr(buf);
		msg.Btype, buf = vt.Gint8(buf);
		msg.Score, buf = vt.Gscore(buf);
	case Tpeer:
		msg.Name, buf = vt.Gstr(buf);
	}

	if(len(buf) != 0){
		fmt.Println("failed to parse syncmsg");
		return nil;
	}

	fmt.Println("syncmsg ", msg);
	return msg;
}

func sync2wire(buf []byte, msg *Syncmsg) int{

	buf2 := vt.Pint8(msg.Id, buf[2:]);
	switch(msg.Id){
	case Thash:
		buf2 = vt.Pstr(msg.Name, buf2);
		buf2 = vt.Pint8(msg.Btype, buf2);
		buf2 = vt.Pscore(msg.Score, buf2);
	case Tpeer:
		buf2 = vt.Pstr(msg.Name, buf2);
	default:
		break;
	}

	n := uint16(len(buf) - len(buf2) - 2)
	vt.Pint16(n, buf);

	return int(n);
}

/* XXX should have a way to get from multiple peers. */
func gotsyncscore(clnt *vtclnt.Clnt, msg *Syncmsg){
	if vthasblock(clnt, msg.Score, msg.Btype) {
		return;
	}

	e := walk(msg.Peer.VTClnt, clnt, msg.Score, msg.Btype)
	if e != nil {
		fmt.Println("gotsyncscore walk: ", e)
		return
	}
	/* XXX add to score table */
	fmt.Println("gotsyncscore ", msg.Score,
		" from ", msg.Peer.Addr);
}

func onsyncmsg(clnt *vtclnt.Clnt, msg *Syncmsg){
	switch msg.Id {
	case Thash:
		go gotsyncscore(clnt, msg);
	case Tpeer:
		// XXX go addpeer(msg.Name);
	}
}

/* We'll need another chan here. */
/* scorechan and peerchan? */
/* Still possibly a separate scoremanager interface / proc. */
/* Interactions: */
/* Receive syncmsg, simple dispatch. */
/* New peer from listen() but without full Tpeer yet */
/* Sync successfully enrolled score in scorelist. */
func syncproc(peerchan chan *Peer, tagchan chan Tag){
	var peers = make(map[string]*Peer)
	var scores = make(map[string]vt.Score)
	_ = peers
	_ = scores
	for {
		select {
		case peer := <- peerchan:
			peers[peer.Name] = peer
		case tag := <- tagchan:
			scores[tag.Name] = tag.Score
		}
	}
}

/*
pseudofunc updateset(id vt.Score){
	/* add to my rootblock

	for each p in peers{
		send syncmsg (peer, id)
	}
}

pseudofunc onaddpeer(p peer){
	for each q in peers{
		sendpeermsg(p, q)
		sendpeermsg(q, p)
	}

	for each id in rootblock{
		sendsyncmsg(p, id);
	}
}
*/
