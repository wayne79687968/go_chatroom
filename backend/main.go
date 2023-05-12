package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"fmt"
	"math/rand"
	"time"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
)

var addr = flag.String("addr", ":8080", "http service address")

type User struct {
	Username string `json:"username"`
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	http.ServeFile(w, r, "index.html")
}

func generateUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(9999) + 1
	
	data := User{fmt.Sprintf("Anonymous#%d", num)}
	jData, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}


func main() {
	flag.Parse()

	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer db.Close()

	hub := newHub()
	go hub.run()

	http.HandleFunc("/", serveHome)
	
	http.HandleFunc("/api/getUsername", generateUsername)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}