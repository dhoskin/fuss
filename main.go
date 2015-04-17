package main
import "code.google.com/p/govt/vt"
import "code.google.com/p/govt/vt/vtclnt"
import "fmt"
import "net"
import "flag"

var addr = flag.String("addr", ":7777", "port to listen");
var serve = flag.Bool("serve", true, "whether to listen");
var venti = flag.String("venti", "localhost:17034", "local venti");

// /usr/drh/src/ventirc.vac
var score = "552dd1ba32295a3daf3925aa4438adb3ed936f18";
var goatscore = "4aaddcaa3881e818d7b4ae0f1aca647994f4b05f";

var Self struct {
	Name string
	Host string
	Clnt *vtclnt.Clnt
	Peerchan chan *Peer
	Tagchan chan Tag
}

// copied from vtclnt.Connect,
// so we can use our own socket
func vtconnect(clnt *vtclnt.Clnt) *vt.Error {
	req := clnt.ReqAlloc()
	req.Done = make(chan *vtclnt.Req)
	tc := &req.Tc
	tc.Id = vt.Thello
	tc.Version = "02"
	tc.Uid = "anonymous"
	tc.Strength = 0
	tc.Crypto = make([]byte, 0)
	tc.Codec = tc.Crypto

	e := clnt.Rpcnb(req)
	if e != nil {
		return e
	}

	<-req.Done
	defer clnt.ReqFree(req)

	if req.Err != nil {
		return req.Err
	}

	return nil;
}

func main(){
	flag.Parse();

	_ = Self.Name

	clntpipe, listenpipe := net.Pipe();

	Self.Clnt = vtclnt.NewClnt(clntpipe);

	go syncproc(Self.Peerchan, Self.Tagchan);

	if(*serve){
		go listendispatch(listenpipe, *venti, *addr);
	}

	var args = flag.Args();
	for i := range args {
		go dialdispatch(listenpipe, *venti, args[i]);
	}

	vtconnect(Self.Clnt);

	myscore, _ := readscore(score);
	_ = myscore;
}

/*
func oldishmain(){
	src, e := vtclnt.Connect("tcp", *venti);
	if(e != nil){
		fmt.Println("could not connect to src: ", e);
		return;
	}
	dst, e := vtclnt.Connect("tcp", "localhost:7034");
	if(e != nil){
		fmt.Println("could not connect to dst: ", e);
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
*/

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
