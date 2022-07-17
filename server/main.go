package main

import (
	"fmt"
	"log"
	"net/http"
	"syscall"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析參數，預設是不會解析的
	// fmt.Println(r.Form) //這些資訊是輸出到伺服器端的列印資訊
	// fmt.Println("path", r.URL.Path)
	// fmt.Println("scheme", r.URL.Scheme)
	// fmt.Println(r.Form["url_long"])
	// for k, v := range r.Form {
	// 	fmt.Println("key:", k)
	// 	fmt.Println("val:", strings.Join(v, ""))
	// }
	fmt.Fprintf(w, "Hello astaxie!") //這個寫入到 w 的是輸出到客戶端的
}

func setRlimitNOFile(nofile uint64) error {
	if nofile == 0 {
		return nil
	}
	var lim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &lim); err != nil {
		return err
	}
	if nofile <= lim.Cur {
		return nil
	}
	lim.Cur = nofile
	lim.Max = nofile
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &lim)
}

func main() {
	setRlimitNOFile(100000)
	http.HandleFunc("/", sayhelloName)
	log.Println("Listen at http://localhost:9090/")
	err := http.ListenAndServe(":9090", nil) //設定監聽的埠
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
