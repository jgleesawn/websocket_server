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
	//"github.com/jmoiron/sqlx"

	"strings"
	"strconv"
	"reflect"

	"encoding/json"

)
/*
type Quest struct {
	Questid		int
	Name		[]byte
	Description	[]byte
	Category	[]byte
	Recurring	bool
	Xpvalue		int
	Requiredquests	[]int
	Attributes	[][]byte
}*/

/*
func (quest Quest) New(qid int, name string, desc string, cat string, rec bool, xpval int, reqquests []int, attr []string) {
	quest.Questid = qid
	quest.Name= name
	quest.Description = desc
	quest.Category = cat
	quest.Recurring = rec
	quest.Xpvalue = xpval
	quest.Requiredquests = reqquests
	quest.Attributes = attr
}*/
type Quest struct {
	Questid		int
	Name		string
	Description	string
	Category	string
	Recurring	bool
	Xpvalue		int
	Requiredquests	[]int
	Attributes	[]string
}
func (quest Quest) New(qid int, name string, desc string, cat string, rec bool, xpval int, reqquests []int, attr []string) {
	quest.Questid = qid
	quest.Name= name
	quest.Description = desc
	quest.Category = cat
	quest.Recurring = rec
	quest.Xpvalue = xpval
	quest.Requiredquests = reqquests
	quest.Attributes = attr
}
/*
type User struct {
	Username	[]byte
	Firstname	[]byte
	Lastname	[]byte
	Xp		int
	Completedquests	[]int
	Attributes	[][]byte
}*/

/*
func (user User) New(u string,f string, l string, a []string) {
	user.Username = u
	user.Firstname = f
	user.Lastname = l
	user.Xp = 0
	user.Completedquests = make([]int,0)
	user.Attributes = a
}*/
type User struct {
	Username	string	`db:"Username"`
	Firstname	string	`db:"Firstname"`
	Lastname	string	`db:"Lastname"`
	Xp		int	`db:"Xp"`
	Completedquests	[]int	`db:"Completedquests"`
	Attributes	[]string `db:"Attributes"`
}
func (user User) New(u string,f string, l string, a []string) {
	user.Username = u
	user.Firstname = f
	user.Lastname = l
	user.Xp = 0
	user.Completedquests = make([]int,0)
	user.Attributes = a
}

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
	//	fmt.Println(reflect.ValueOf(conn.ReadMessage).Type())
		_,data,err := conn.ReadMessage()
		if len(data) > 0{
			process(data,db,conn)
		} else if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(1*time.Second)
	}
}
func process(data []byte, db *sql.DB, conn *websocket.Conn){
	mt := websocket.TextMessage
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
			success := addUser(db,&u)
			if success {
			conn.WriteMessage(mt,[]byte("User added."))
			} else {
			conn.WriteMessage(mt,[]byte("Couldn't add user."))
			}
			break
		case "add Quest":
			q := Quest{}
			err := json.Unmarshal([]byte(str[1]),&q)
			if err != nil {
				log.Println("Struct doesn't match command.")
			}
			q.Attributes = append(q.Attributes,"")
			fmt.Println(q.Attributes)
			success := addQuest(db,&q)
			if success {
			conn.WriteMessage(mt,[]byte("Quest added."))
			} else {
			conn.WriteMessage(mt,[]byte("Error on adding quest."))
			}
			break
		//case "update Auth":
		case "update User":
			u := User{}
			err := json.Unmarshal([]byte(str[1]),&u)
			if err != nil {
				log.Println("Struct doesn't match command.")
			}
			success := updateUser(db,&u)
			if success {
			conn.WriteMessage(mt,[]byte("User Updated."))
			} else {
			conn.WriteMessage(mt,[]byte("Error on Updating User."))
			}
			break
		case "update Quest":
			q := Quest{}
			err := json.Unmarshal([]byte(str[1]),&q)
			if err != nil {
				log.Println("Struct doesn't match command.")
			}
			success := updateQuest(db,&q)
			if success {
			conn.WriteMessage(mt,[]byte("Quest Updated."))
			} else {
			conn.WriteMessage(mt,[]byte("Error on Updating Quest."))
			}
			break
	}
}

func addAuth(db *sql.DB,username string,password string) (bool) {
_,err := db.Query(`INSERT INTO auth VALUES($1,$2);`,username,password) 
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func addUser(db *sql.DB,u *User) (bool){
	strq := []string{`INSERT INTO users VALUES(`,`,`,`,`,`,`,`, ARRAY[`,`],ARRAY[`,`]);`}
	str,varargs := unroll_query(strq,u.Username,u.Firstname,u.Lastname,u.Xp,u.Completedquests,u.Attributes)
	_,err := db.Query(str,varargs...)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func addQuest(db *sql.DB,q *Quest) (bool) {
row := db.QueryRow(`SELECT questid FROM quests ORDER BY questid DESC LIMIT 1;`)
	var prev_id int
	err := row.Scan(&prev_id)
	if err != nil {
		q.Questid = 0
	} else {
		q.Questid = prev_id + 1
	}
	strq := []string{`INSERT INTO quests VALUES(`,`,`,`,`,`,`,`,`,`,`,`,ARRAY[`,`],ARRAY[`,`]);`}
	str,varargs := unroll_query(strq,q.Questid,q.Name,q.Description,q.Category,q.Recurring,q.Xpvalue,q.Requiredquests,q.Attributes)
	fmt.Println(str,varargs)
	_,err = db.Query(str,varargs...)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func updateUser(db *sql.DB,u *User) (bool){
	strq := []string{`UPDATE users SET (firstname, lastname, xp, completedquests, attributes) = (`,`,`,`,`,`,ARRAY[`,`],ARRAY[`,`]) WHERE username = `,`;`}
	str,varargs := unroll_query(strq,u.Firstname,u.Lastname,u.Xp,u.Completedquests,u.Attributes,u.Username)
	_,err := db.Query(str,varargs...)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func updateQuest(db *sql.DB,q *Quest) (bool){
	strq := []string{`UPDATE quests SET (name, description, category, recurring, xpvalue, requiredquests, attributes) = (`,`,`,`,`,`,`,`,`,`,ARRAY[`,`],ARRAY[`,`]) WHERE questid = `,`;`}
	str,varargs := unroll_query(strq,q.Name,q.Description,q.Category,q.Recurring,q.Xpvalue,q.Requiredquests,q.Attributes,q.Questid)
	_,err := db.Query(str,varargs...)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}



/*
func getQuestByQid(db *sql.DB,qid int) (*Quest) {
	quest := Quest{}
	err := db.Get(&quest,`SELECT * FROM quests WHERE questid = $1`,qid)
	if err != nil {
		log.Println(err)
	}
	return &quest
}
func getQuestByName(db *sql.DB,name string) ([]Quest){
	quests := []Quest{}
	err := db.Select(&quests,`SELECT * FROM quests WHERE name = $1`,name)
	if err != nil {
		log.Println(err)
	}
	return quests
}
func getQuestByCategory(db *sql.DB,category string) ([]Quest) {
	quests := []Quest{}
	err := db.Select(&quests,`SELECT * FROM quests WHERE category = $1`,category)
	if err != nil {
		log.Println(err)
	}
	return quests
}
*/
//takes comma seperated query, seperations lie on variable placement
//fills out this query with $1 $2 etc.
//unrolls arrays and slices in varargs, so you can insert arrays
func unroll_query(strq []string, varargs ...interface{}) (string,[]interface{}) {
	stro := ""
	argo := make([]interface{},0,len(strq)-1)
	count := 1
	for i := range varargs {
		stro += strq[i]
		k := reflect.ValueOf(varargs[i])
		if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		   if k.Len() > 0 {
			for c := 0; c < k.Len()-1; c++ {
				v := reflect.ValueOf(k.Index(c).Interface())
//				if v.Kind() == reflect.String {
//					stro += `'`
//				}
				fmt.Println(v.Kind())
				argo = append(argo,k.Index(c).Interface())
				stro += `$`+strconv.Itoa(count)
				count++
				if v.Kind() == reflect.Int {
					stro += `::integer`
				} else if v.Kind() == reflect.String {
					stro += `::text`
				}
				stro += `,`
			}
			v := reflect.ValueOf(k.Index(k.Len()-1).Interface())
			fmt.Println(v.Kind())
//			if v.Kind() == reflect.String {
//				stro += `'`
//			}
			stro += `$`+strconv.Itoa(count)
			count++
			if v.Kind() == reflect.Int {
				stro+= `::integer`
			} else if v.Kind() == reflect.String {
				stro += `::text`
			}

			argo = append(argo,k.Index(k.Len()-1).Interface())
		   }
		} else {
			stro += `$`+strconv.Itoa(count )
			count++
			argo = append(argo,varargs[i])
		}
	}
	stro += strq[len(strq)-1]
	fmt.Println(stro)
	for i := range argo {
		k := reflect.ValueOf(argo[i]).Kind()
		fmt.Println(k,": ",argo[i])
	}
	fmt.Println(len(argo))
	fmt.Println(argo)
	return stro,argo
}
