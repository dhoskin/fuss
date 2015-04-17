package main

import "fmt"
import "code.google.com/p/govt/vt"

const(
	Thash = 20
	Rhash = 21
	Tpeer = 22
	Rpeer = 23
)

// sz[2] Thash score[20]
// sz[2] Tpeer name[s]

type Syncmsg struct {
	Id byte
	Score vt.Score
	Name string
}

func (s *Syncmsg) String() string {
	switch s.Id {
	case Thash:
		return fmt.Sprintf("Thash %v", s.Score);
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

	