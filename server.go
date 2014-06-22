package main

import (
	"fmt"
	"os"
	"log"
	"time"
	
	"net/http"
	"github.com/gorilla/websocket"

	"database/sql"
	_ "github.com/lib/pq"

	"strings"
	"strconv"

)

func main() {

	fmt.Println("listening...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
 	}
	http.HandleFunc("/",webHandler)
	http.HandleFunc("/ws",wsHandler)
	err := http.ListenAndServe(":"+port,nil)
	if err != nil {
		panic(err)
	}
//		go handleConnection(conn)
}
func webHandler(res http.ResponseWriter, req *http.Request){
	fmt.Fprintln(res, "webpage")
}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
func wsHandler(res http.ResponseWriter, req *http.Request){
	conn, err := upgrader.Upgrade(res,req,nil)
	if err != nil {
		log.Println(err)
		return
	}
	mt := websocket.TextMessage
	conn.WriteMessage(mt,[]byte("Websocket connected."))

	db_url := os.Getenv("DATABASE_URL")
	db_name := "websocket_db"
	sslmode := "sslmode=disable"
	db, err := sql.Open("postgres", db_url+"/"+db_name+"?"+sslmode)
	if err != nil {
		panic(err)
	}

	for {
		_,data,err := conn.ReadMessage()
		if len(data) > 0{
			go process(data,db,conn)
		}
		if err != nil {
			panic(err)
		}
		time.Sleep(1*time.Second)
	}
}
func process(data []byte, db *sql.DB, conn *websocket.Conn){
	mt := websocket.TextMessage
	str := strings.Fields(string(data))
	fmt.Println(str)
	var age int
	var name []byte
	switch str[0] {
		case "get":
//			fmt.Println("get")
			copy(name,[]byte(str[1]))
			fmt.Println(string(name))
rows, err := db.Query(`SELECT name, age FROM users WHERE name = $1;`,string(name))
			if err != nil {
				panic(err)
			}
			rows.Next()
			rows.Scan(&name,&age)
out := []byte(string(name)+" is "+strconv.Itoa(age)+" years old.")
			conn.WriteMessage(mt,out)
			break;
		case "store":
			copy(name,[]byte(str[1]))
			age,_ := strconv.Atoi(string(str[2]))
db.QueryRow(`INSERT INTO users VALUES($1,$2);`,string(name),int(age))
			fmt.Println(string(name))
			fmt.Println(strconv.Itoa(age))
			break;
	}
}


