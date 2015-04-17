package main

import "code.google.com/p/govt/vt"

/* Too heavy-handed? return -1? */
func rhex(c byte) (byte, *vt.Error) {
	if c >= '0' && c <= '9' {
		return c - '0', nil;
	}else if c >= 'a' && c <= 'f' {
		return c - 'a' + 10, nil;
	}else if c >= 'A' && c <= 'F' {
		return c - 'A' + 10, nil;
	} else {
		return 0, &vt.Error{"not hex"};
	}
}

/* XXX submit as patch for govt */
/* see /sys/src/libventi/parsescore.c */
func readscore(s string) (vt.Score, *vt.Error) {
	/* How is this done in idiomatic go? */
	buf := make(vt.Score, vt.Scoresize)

	return buf, readScore(s, buf);
}

func readScore(s string, buf vt.Score) *vt.Error {
	/* XXX should handle vac: prefix */
	if len(s) != vt.Scoresize * 2 {
		return &vt.Error{"bad score"}
	}
	if len(buf) < vt.Scoresize {
		return &vt.Error{"bad score buf"}
	}

	for i:= 0; i < vt.Scoresize; i++ {
		c, e := rhex(s[2*i])
		if e != nil {
			return e;
		}
		buf[i] = c << 4;
		c, e = rhex(s[2*i + 1]);
		if e != nil {
			return e;
		}
		buf[i] |= c;
	}

	return nil;
}
