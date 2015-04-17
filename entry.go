package main
import "code.google.com/p/govt/vt"
import "fmt"

/* /sys/src/libventi/entry.c */

const (
	EntryActive = 1 << 0
	_EntryDir = 1 << 1
	_EntryDepthShift = 2
	_EntryDepthMask = 7 << 2
	EntryLocal = 1 << 5
)

type VtEntry struct {
	gen uint32
	psize uint16
	dsize uint16
	t byte
	flags byte
	size uint64
	score vt.Score
}

func vtentrypack(entry *VtEntry, data []byte, index int) {
	return;
}

//func vtentryunpack(entry *VtEntry, data []byte, index int) {
//	return;
//}

func vtentryunpack(data []byte, index int) *VtEntry {
	var e = new(VtEntry);
	data = data[index * vt.Entrysize:index * vt.Entrysize + vt.Entrysize]
	e.gen, data = vt.Gint32(data);
	e.psize, data = vt.Gint16(data);
	e.dsize, data = vt.Gint16(data);
	e.flags, data = vt.Gint8(data);
	if((e.flags & _EntryDir) == 1){
		e.t = vt.DirBlock;
	}else{
		e.t = vt.DataBlock;
	}
	e.t += (e.flags & _EntryDepthMask) >> _EntryDepthShift;
	e.flags &= ^(byte)(_EntryDir | _EntryDepthMask);
	data = data[5:];
	e.size, data = vt.Gint48(data);
	e.score = make(vt.Score, vt.Scoresize);
	copy(e.score, data);
	data = data[vt.Scoresize:];

	if(len(data) != 0){
		fmt.Println("entry size wrong! ", len(data));
		return nil;
	}
	return e;
}

func vtentryunpackseq(data []byte) (*VtEntry, []byte) {
	return vtentryunpack(data, 0), data[vt.Entrysize:];
}
