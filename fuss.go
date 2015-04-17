package main
import "9fans.net/go/plan9"
import "code.google.com/p/govt/vt/vtclnt"
import "fmt"

func main(){
	_ = vtclnt.DefaultDebuglevel
	var v9p = plan9.VERSION9P
	fmt.Printf("%s goats\n", v9p);
}
