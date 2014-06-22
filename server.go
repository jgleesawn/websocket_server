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
	url := `onelyfe.herokuapp.com`
	url = req.Host
	fmt.Fprintln(res, `<script type='text/javascript'> ws = new WebSocket('ws://`+url+`/ws'); ws.onmessage = function (event) {curDiv = addElement(); document.getElementById(curDiv).innerHTML = event.data;}; function get(){ ws.send("get "+document.getElementById("name").value) }; function store(){ ws.send("store "+document.getElementById("name").value+" "+document.getElementById("age").value); }; function getall(){ ws.send("all");}; </script>`)
	fmt.Fprintln(res, `<div id='input'>`)
	fmt.Fprintln(res, "name:<input type='text' id='name' name='name' value='oldman'>age:<input type='text' id='age' name='age' value='132'>")
	fmt.Fprintln(res, "<button onclick='get()'>Get</button>")
	fmt.Fprintln(res, "<button onclick='store()'>Store</button>")
	fmt.Fprintln(res, "<button onclick='getall()'>Entire Table</button>")
	fmt.Fprintln(res, "<button onclick='removeElements()'>Clear</button>")
	fmt.Fprintln(res, "</div>")
	fmt.Fprintln(res, `<div id='output'></div>`)
	fmt.Fprintln(res, `<script type='text/javascript'> function addElement() {var ni = document.getElementById('output'); var newdiv = document.createElement('div'); var div_id = Math.random().toString(36).substring(7); newdiv.setAttribute('id',div_id); ni.appendChild(newdiv); return div_id;};</script>`)
	fmt.Fprintln(res, `<script type='text/javascript'> function removeElements() {var out = document.getElementById('output');  for (i = out.childElementCount-1;i>=0;i--) {out.removeChild(out.childNodes[i])};};</script>`)
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
	//fmt.Println(str)
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
		case "all":
rows,err := db.Query(`SELECT * FROM users`)
			if err != nil {
				log.Println(err)
				return
			}
			for rows.Next() {
				rows.Scan(&name,&age)
				mt = websocket.TextMessage
out := []byte("Name: "+string(name)+" Age: "+strconv.Itoa(age))
				conn.WriteMessage(mt,out)
			}
	}
}


