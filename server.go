package main
//Quest and Users must have one element in their arrays if you are trying to add or update
import (
	"fmt"
	"os"
	"log"
	"time"
	
	"net/http"
	"io"
	"github.com/gorilla/websocket"

	//"database/sql"
	//"github.com/lib/pq"
	//"github.com/jmoiron/sqlx"

	"strings"
	//"strconv"
	//"reflect"

	"encoding/json"

	"github.com/jgleesawn/ECC_Conn"
)


func main() {

	fmt.Println("listening...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
 	}
	http.HandleFunc("/",webHandler)
	http.HandleFunc("/ws",wsHandler)
	http.HandleFunc("/noencws",wsnoencHandler)
	err := http.ListenAndServe(":"+port,nil)
	if err != nil {
		panic(err)
	}
}
func webHandler(res http.ResponseWriter, req *http.Request){
	url := `onelyfe.herokuapp.com`
	url = req.Host
	fmt.Fprintln(res, `<script type='text/javascript'> ws = new WebSocket('ws://`+url+`/noencws'); ws.onmessage = function (event) {curDiv = addElement(); document.getElementById(curDiv).innerHTML = event.data;}; function get(){ ws.send("get "+document.getElementById("name").value) }; function store(){ ws.send("store "+document.getElementById("name").value+" "+document.getElementById("age").value); }; function getall(){ ws.send("all");}; </script>`)
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
type wsinterface interface {
//int is websocket.messageType
	WriteMessage(int,[]byte)(error)
	ReadMessage()(int,[]byte,error)
}
type noencWs struct {
	wsinterface
	PacketSize	int
	PayloadLen	int
}
func (x *noencWs) Write(p []byte) (n int, err error){
	start := 0
	end := len(p)
	if end > x.PayloadLen {
		end = x.PayloadLen
	}
	for end < len(p) {
		err = x.WriteMessage(websocket.BinaryMessage,p[start:end])
		if err != nil {
			return start,err
		}
		start = end
		end += x.PayloadLen
	}
	err = x.WriteMessage(websocket.BinaryMessage,p[start:len(p)])
	if err != nil {
		return start,err
	}
	return len(p),err
}
func (x *noencWs) Read(p []byte) (n int, err error){
	_,data,err := x.ReadMessage()
	l := len(data)
	copy(p[:l],data)
	return l,err
}
func wsnoencHandler(res http.ResponseWriter, req *http.Request){
	conn, err := upgrader.Upgrade(res,req,nil)
	if err != nil {
		log.Println(err)
		return
	}
	//makes payload same size as encrypted packet
	noencConn := noencWs{conn,1024,1024-(32+10)}
	
	noencConn.Write([]byte("Websocket connected."))
	db := OpenDB()
	data := make([]byte,noencConn.PacketSize)
	for {
		//fmt.Println(reflect.ValueOf(conn.ReadMessage).Type())
		n,err := noencConn.Read(data)
		if n > 0{
			process(data[:n],db,&noencConn)
		} else if err != nil {
			log.Println(err)
			return
		}
		//time.Sleep(1*time.Second)
	}
}
func wsHandler(res http.ResponseWriter, req *http.Request){
	conn, err := upgrader.Upgrade(res,req,nil)
	if err != nil {
		log.Println(err)
		return
	}

	dh_conn := new(ECC_Conn.ECC_Conn)
	dh_conn.Connect(conn)
	fmt.Println("Outside diffie.")

	dh_conn.Write([]byte("Websocket connected."))
	db := OpenDB()
	data := make([]byte,dh_conn.PacketSize)
	for {
		//fmt.Println(reflect.ValueOf(conn.ReadMessage).Type())
		n,err := dh_conn.Read(data)
		if n > 0{
			process(data[:n],db,dh_conn)
		} else if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(1*time.Second)
	}
}


func process(data []byte, db Custom_db, conn io.ReadWriter){//*ECC_Conn.ECC_Conn){
	str := strings.Split(string(data),";")
	cmd := strings.Fields(str[0])
	//fmt.Println(str)
	switch strings.Join(cmd," ") {
		//case "add Auth":
		case "add User":
			u := User{}
			err := json.Unmarshal([]byte(str[1]),&u)
			if err != nil {
				log.Println("Struct doesn't match command.")
			}
			u.Xp = 0
			u.Completedquests = []int{0}
			success := db.AddUser(&u)
			if success {
				conn.Write([]byte("User added."))
			} else {
				conn.Write([]byte("Couldn't add user."))
			}
			break
		case "add Quest":
			q := Quest{}
			fmt.Println(string(str[1]))
			err := json.Unmarshal([]byte(str[1]),&q)
			if err != nil {
				log.Println(err)
				log.Println("Struct doesn't match command.")
			}
//Removed for consistency with updateQuest
//			q.Attributes = append(q.Attributes,"")
			success := db.AddQuest(&q)
			if success {
				conn.Write([]byte("Quest added."))
			} else {
				conn.Write([]byte("Error on adding quest."))
			}
			break
		//case "update Auth":
		case "update User":
			u := User{}
			err := json.Unmarshal([]byte(str[1]),&u)
			if err != nil {
				log.Println("Struct doesn't match command.")
			}
			success := db.UpdateUser(&u)
			if success {
				conn.Write([]byte("User Updated."))
			} else {
				conn.Write([]byte("Error on Updating User."))
			}
			break
		case "update Quest":
			q := Quest{}
			err := json.Unmarshal([]byte(str[1]),&q)
			if err != nil {
				log.Println("Struct doesn't match command.")
			}
			success := db.UpdateQuest(&q)
			if success {
				conn.Write([]byte("Quest Updated."))
			} else {
				conn.Write([]byte("Error on Updating Quest."))
			}
			break
		case "get User":
			strv := strings.Fields(str[1])[0]
			fmt.Println(strv)
			db.GetUser(strv)
			break
		case "get Quest":
			break
	}
}

