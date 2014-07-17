package main
//Quest and Users must have one element in their arrays if you are trying to add or update
import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"time"
	
	"net/http"
	"io"
	"github.com/gorilla/websocket"

	//"database/sql"
	//"github.com/lib/pq"
	//"github.com/jmoiron/sqlx"

	"strings"
	"strconv"
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
	fmt.Fprintln(res,`<script type='text/javascript'>ws = new WebSocket('ws://`+url+`/noencws');</script>`)
	js_code_arr, err := ioutil.ReadFile("index.js")
	if err != nil {
		fmt.Fprintln(res,[]byte("Error loading page."))
	}
	fmt.Fprintln(res, string(js_code_arr))
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
	/*
	start := 0
	end := len(p)
	if end > x.PayloadLen {
		end = x.PayloadLen
	}
	for end < len(p) {
		err = x.WriteMessage(websocket.TextMessage,p[start:end])
		if err != nil {
			return start,err
		}
		start = end
		end += x.PayloadLen
	}
	*/
	err = x.WriteMessage(websocket.TextMessage,p)
	if err != nil {
		return 0,err
		//return start,err
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
	//log.Println(req.Header)
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
	out := []byte("Didn't Process Command.")
	str := strings.Split(string(data),";")
	if len(str) < 2 {
		out = []byte("Command missing arguments.")
		conn.Write(out)
		return
	}
	cmd := strings.Fields(str[0])
	//fmt.Println(str)
	switch strings.Join(cmd," ") {
		//case "add Auth":
		case "add User":
			u := User{}
			err := json.Unmarshal([]byte(str[1]),&u)
			if err != nil {
				log.Println("Struct doesn't match command.")
				break
			}
			u.Xp = 0
			u.Completedquests = []int{0}
			success := db.AddUser(&u)
			if success {
				out = []byte("User added.")
			} else {
				out = []byte("Couldn't add user.")
			}
			break
		case "add Quest":
			q := Quest{}
			fmt.Println(string(str[1]))
			err := json.Unmarshal([]byte(str[1]),&q)
			if err != nil {
				log.Println(err)
				log.Println("Struct doesn't match command.")
				break
			}
//Removed for consistency with updateQuest
//			q.Attributes = append(q.Attributes,"")
			success := db.AddQuest(&q)
			if success {
				out = []byte("Quest added.")
			} else {
				out = []byte("Error on adding quest.")
			}
			break
		//case "update Auth":
		case "update User":
			u := User{}
			err := json.Unmarshal([]byte(str[1]),&u)
			if err != nil {
				log.Println("Struct doesn't match command.")
				break
			}
			success := db.UpdateUser(&u)
			if success {
				out = []byte("User Updated.")
			} else {
				out = []byte("Error on Updating User.")
			}
			break
		case "update Quest":
			q := Quest{}
			err := json.Unmarshal([]byte(str[1]),&q)
			if err != nil {
				log.Println("Struct doesn't match command.")
				break
			}
			success := db.UpdateQuest(&q)
			if success {
				out = []byte("Quest Updated.")
			} else {
				out = []byte("Error on Updating Quest.")
			}
			break
		case "get User":
			username := string("")
			//err := json.Unmarshal([]byte(str[1]),&username)
			err := error(nil)
			username = strings.Fields(str[1])[0]
			if err != nil {
				out = []byte("Couldn't Unmarshal request.")
				break
			}
			users,err := db.GetUser(username)
			if err != nil {
				log.Println(err)
				out = []byte("Error getting user from db.")
			}
			out,err = json.Marshal(users)
			if err != nil {
				out = []byte("Couldn't put users in json form.")
				break
			}
			break
		case "get Quest":
			questid := int64(0)
			//err := json.Unmarshal([]byte(str[1]),&questid)
			questid, err := strconv.ParseInt(str[1],10,64)
			if err != nil {
				out = []byte("Couldn't Unmarshal request.")
			}
			quests,err := db.GetQuest(questid)
			if err != nil {
				log.Println(err)
				out = []byte("Error getting quest from db.")
			}
			out,err = json.Marshal(quests)
			if err != nil {
				out = []byte("Couldn't put quests in json form.")
				break
			}
			break
		case "get all Users":
			users,err := db.GetAllUsers()
			if err != nil {
				log.Println(err)
				out = []byte("Error getting users from db.")
			}
			out,err = json.Marshal(users)
			if err != nil {
				out = []byte("Couldn't put users in json form.")
				break
			}
			break
		case "get all Quests":
			quests,err := db.GetAllQuests()
			if err != nil {
				log.Println(err)
				out = []byte("Error getting quests from db.")
			}
			out,err = json.Marshal(quests)
			if err != nil {
				out = []byte("Couldn't put quests in json form.")
				break
			}
			break
	}
	//fmt.Println(out)
	conn.Write(out)
}

