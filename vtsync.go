package main
import "code.google.com/p/govt/vt"
import "code.google.com/p/govt/vt/vtclnt"
import "fmt"

// /usr/drh/src/ventirc.vac
var score = "552dd1ba32295a3daf3925aa4438adb3ed936f18";
var goatscore = "4aaddcaa3881e818d7b4ae0f1aca647994f4b05f";

/*
pseudofunc walk(src, dst, hash, type){
	if(dst.read(hash, type))
		return true;
	if((block = src.read(hash, type)) == error)
		return false;

	for(child hash in block){
		if(walk(src, dst, child, childtype) == error)
			return false
	}
	return dst.write(block, type) == hash;
}
*/

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

func main(){
	src, e := vtclnt.Connect("tcp", "localhost:17034");
	if(e != nil){
		fmt.Println("could not connect to src: ", e);
		return;
	}
	dst, e := vtclnt.Connect("tcp", "localhost:7034");
	if(e != nil){
		fmt.Println("could not connect to src: ", e);
		return;
	}

	myscore, _ := readscore(score);

	fmt.Println("prepare to walk");
	e = walk(src, dst, myscore, vt.RBlock);
	if(e != nil){
		fmt.Println("walk: ", e);
		return;
	}
	fmt.Println("walk done; prep to sync");
	e = dst.Sync();
	if(e != nil){
		fmt.Println("sync: ", e);
		return;
	}

	fmt.Println("end of main");
	return;
}

func oldmain(){
	clnt, e := vtclnt.Connect("tcp", "localhost:17034");
	if e != nil {
		fmt.Println("could not connect: ", e);
		return;
	}

	myscore, e := readscore(score);

	if e != nil {
		fmt.Println("score error: %s", e.Ename);
		fmt.Println(myscore);
	}else{
		fmt.Println("vac:", myscore);
	}

	var t uint8; // how to init in for loop?
	for t = 0; t <= vt.RBlock; t++ {
		data, e := clnt.Get(myscore, t, vt.Maxblock);
		if(e == nil){
			if(len(data) > 0){
				fmt.Println("found block, type ", t);
				fmt.Println("length ", len(data));
				fmt.Println(data);
				break;
			}
			fmt.Println("no block of type ", t);
		}else{
			fmt.Println("error:", e.Ename);
		}
	}
}
