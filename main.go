package main

import (
	"net/http"
)

// hello world output
func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<h1>Testing</h1>`))
}

func main() {
	//register url to servermux
	http.HandleFunc("/sayHello", sayHello)

	http.ListenAndServe(":3001", nil)
	//fmt.Println("Testing")

}
