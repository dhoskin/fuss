package main

import "fmt"
import "code.google.com/p/govt/vt"
import "code.google.com/p/govt/vt/vtclnt"

func vthasblock(c *vtclnt.Clnt, hash vt.Score, btype byte) bool{
	buf, e := c.Get(hash, btype, vt.Maxblock);
	if(e != nil){
		return false;
	}
	return len(buf) > 0;
}

func walk(src *vtclnt.Clnt, dst *vtclnt.Clnt,
	hash vt.Score, btype byte) (*vt.Error) {

	fmt.Println("walk type ", btype);
	fmt.Println(hash);
	buf, e := dst.Get(hash, btype, vt.Maxblock)
	if(e != nil){
		return e;
	}
	if(len(buf) > 0){
		return nil;
	}

	buf, e = src.Get(hash, btype, vt.Maxblock);
	if(e != nil){
		return e;
	}
	//if(len(buf) == 0){
	//	return nil;
	//}

	switch(btype){
	case vt.RBlock:{
		var root = vtrootunpack(buf);
		if(root == nil){
			return &vt.Error{"bad root"};
		}
		e = walk(src, dst, root.score, vt.DirBlock);
		if(e != nil){
			return e;
		}
		/* XXX root.prev */
		break;
	}
	case vt.DirBlock:{
		for i := 0; i * vt.Entrysize < len(buf); i++ {
			entry := vtentryunpack(buf, i);
			if(entry == nil){
				return &vt.Error{"bad entry"};
			}
			e = walk(src, dst, entry.score, entry.t);
			if(e != nil){
				return e;
			}
		}
		break;
	}
	case vt.DataBlock:{
		break;
	}
	default:{ /* pointers */
		for i := 0; i < len(buf);i += vt.Scoresize{
			e = walk(src, dst, buf[i:], btype - 1);
			if(e != nil){
				return e;
			}
		}
		break;
	}
	}

	newscore, e := dst.Put(btype, buf);
	if(e != nil){
		return e;
	}
	fmt.Println(newscore);
	for i := 0; i < vt.Scoresize; i++{
		if(hash[i] != newscore[i]){
			fmt.Println("score mismatch!");
			return &vt.Error{"score mismatch!"};
		}
	}
	return nil;
}
