package main

import (
	"fmt"
	"os"
	"log"
	"time"
	
	"net/http"
	"github.com/gorilla/websocket"

	"database/sql"
	"github.com/lib/pq"

	"strings"
	"strconv"

)

func openDB() *sql.DB {
	url := os.Getenv("DATABASE_URL")
	connection, _ := pq.ParseURL(url)
	connection += " sslmode=require"

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Println(err)
	}

	return db
}

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

	db := openDB()

	for {
		_,data,err := conn.ReadMessage()
		if len(data) > 0{
			go process(data,db,conn)
		} else if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(1*time.Second)
	}
}
func process(data []byte, db *sql.DB, conn *websocket.Conn){
	mt := websocket.TextMessage
	str := strings.Fields(string(data))
	fmt.Println(str)
	var age int
	name := make([]byte,50)
	switch string(str[0]) {
		case "get":
			l := copy(name,str[1])
rows, err := db.Query(`SELECT name, age FROM users WHERE name = $1;`,name[0:l])
			if err != nil {
				log.Println(err)
conn.WriteMessage(mt,[]byte("process switch case 'get' DB select error."))
				return
			}
			rows.Next()
			rows.Scan(&name,&age)
out := []byte(string(name)+" is "+strconv.Itoa(age)+" years old.")
			conn.WriteMessage(mt,out)
			break;

		case "store":
			l := copy(name,str[1])
			age,err := strconv.Atoi(string(str[2]))
			if err != nil {
				log.Println(err)
				return
			}
db.QueryRow(`INSERT INTO users VALUES($1,$2);`,name[0:l],int(age))
			mt = websocket.TextMessage
			out := "Row Inserted."
			conn.WriteMessage(mt,[]byte(out))
			break;
	}
}


