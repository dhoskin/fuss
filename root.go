package main
import "code.google.com/p/govt/vt"

/* /sys/src/libventi/root.c */

const (
	RootSize = 300
	RootVersion = 2
)

type Root struct {
	name [128]byte
	t [128]byte
	score vt.Score
	blocksize uint16
	prev vt.Score
}

/* XXX vtrootpack */

func vtrootunpack(buf []byte) *Root{
	var root = new(Root);

	vers, buf := vt.Gint16(buf);
	if(vers != RootVersion){
		return nil;
	}

	copy(root.name[0:len(root.name)], buf);
	buf = buf[128:];
	copy(root.t[0:len(root.t)], buf);
	buf = buf[128:];
	root.score = make(vt.Score, vt.Scoresize);
	copy(root.score, buf);
	buf = buf[vt.Scoresize:];
	root.blocksize, buf = vt.Gint16(buf);
	/* XXX checksize(root.blocksize) */
	root.prev = make(vt.Score, vt.Scoresize);
	copy(root.prev, buf);
	buf = buf[vt.Scoresize:];

	return root;
}
