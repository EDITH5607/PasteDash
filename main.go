package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w,r)
		return 
	}
	w.Write([]byte("Hello from PasteDash...."))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json" )
	w.Header()["Date"] = nil
	w.Header().Add("name","aadharsh")
	w.Write([]byte(`{"name":"aadharsh"}`))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // this way we can send tokens for verification.
		http.Error(w, "post is not allowed ", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new Snippet !!"))
}

func main() {
	// router
	mux := http.NewServeMux()


	//routes 
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)


	//server initialization and starting
	log.Println("Starting Server on :127.0.0.1:4000")
	err := http.ListenAndServe("127.0.0.1:4000", mux)
	log.Fatal(err)

}
