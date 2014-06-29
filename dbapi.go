package main

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"os"
	"reflect"
	"strings"

	"strconv"
	"fmt"
)

type Custom_db struct {
	sql.DB
}

func OpenDB() Custom_db {
	url := os.Getenv("DATABASE_URL")
	connection,_ := pq.ParseURL(url)
	connection += " sslmode=require"

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Println(err)
	}
	cdb := Custom_db{*db}

	return cdb
}

func (db *Custom_db) AddAuth(username string,password string) (bool) {
_,err := db.Query(`INSERT INTO auth VALUES($1,$2);`,username,password) 
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (db *Custom_db) AddUser(u *User) (bool){
	strq := []string{`INSERT INTO users VALUES(`,`,`,`,`,`,`,`, ARRAY[`,`],ARRAY[`,`]);`}
	str,varargs := unroll_query(strq,u.Username,u.Firstname,u.Lastname,u.Xp,u.Completedquests,u.Attributes)
	_,err := db.Query(str,varargs...)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (db *Custom_db) AddQuest(q *Quest) (bool) {
row := db.QueryRow(`SELECT questid FROM quests ORDER BY questid DESC LIMIT 1;`)
	var prev_id int64
	err := row.Scan(&prev_id)
	if err != nil {
		q.Questid = 0
	} else {
		q.Questid = prev_id + 1
	}
	strq := []string{`INSERT INTO quests VALUES(`,`,`,`,`,`,`,`,`,`,`,`,ARRAY[`,`],ARRAY[`,`]);`}
	str,varargs := unroll_query(strq,q.Questid,q.Name,q.Description,q.Category,q.Recurring,q.Xpvalue,q.Requiredquests,q.Attributes)
	_,err = db.Query(str,varargs...)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (db *Custom_db) UpdateUser(u *User) (bool){
	strq := []string{`UPDATE users SET (firstname, lastname, xp, completedquests, attributes) = (`,`,`,`,`,`,ARRAY[`,`],ARRAY[`,`]) WHERE username = `,`;`}
	str,varargs := unroll_query(strq,u.Firstname,u.Lastname,u.Xp,u.Completedquests,u.Attributes,u.Username)
	_,err := db.Query(str,varargs...)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (db *Custom_db) UpdateQuest(q *Quest) (bool){
	strq := []string{`UPDATE quests SET (name, description, category, recurring, xpvalue, requiredquests, attributes) = (`,`,`,`,`,`,`,`,`,`,ARRAY[`,`],ARRAY[`,`]) WHERE questid = `,`;`}
	str,varargs := unroll_query(strq,q.Name,q.Description,q.Category,q.Recurring,q.Xpvalue,q.Requiredquests,q.Attributes,q.Questid)
//	fmt.Println(str)
	_,err := db.Query(str,varargs...)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (db *Custom_db) GetUser(u string) (interface{},error) { //(*User){
	rows,err :=db.Query(`SELECT * FROM users WHERE username = $1;`,u)
	if err != nil {
		log.Println(err)
		return  nil,err
	}

	skel := reflect.ValueOf(User{"a","a","a",0,[]int{0},[]string{""}})
	data := RowData(skel,rows)
	out := make([]User,len(data))
	for i := 0; i<len(data); i++ {
		out[i].New(data[i])
	}
	return out,nil
}
func (db *Custom_db) GetQuest(qid int64) (interface{},error) {
	rows,err := db.Query(`SELECT * FROM quests WHERE questid = $1;`,qid)
	if err != nil {
		log.Println(err)
		return nil,err
	}

	skel := reflect.ValueOf(Quest{0,"a","a","a",true,0,[]int{0},[]string{""}})
	data := RowData(skel,rows)
	out := make([]Quest,len(data))
	for i := 0; i<len(data); i++ {
		out[i].New(data[i])
	}
	return out,nil
}
func (db *Custom_db) GetAllUsers() (interface{},error) {
	rows,err :=db.Query(`SELECT * FROM users WHERE username is not null;`)
	if err != nil {
		log.Println(err)
		return  nil,err
	}

	skel := reflect.ValueOf(User{"a","a","a",0,[]int{0},[]string{""}})
	data := RowData(skel,rows)
	out := make([]User,len(data))
	fmt.Println(data)
	fmt.Println(len(data))
	for i := 0; i<len(data); i++ {
		out[i].New(data[i])
	}
	return out,nil
}
func (db *Custom_db) GetAllQuests() (interface{},error) {
	rows,err := db.Query(`SELECT * FROM quests WHERE questid is not null;`)
	if err != nil {
		log.Println(err)
		return nil,err
	}

	skel := reflect.ValueOf(Quest{0,"a","a","a",true,0,[]int{0},[]string{""}})
	data := RowData(skel,rows)
	out := make([]Quest,len(data))
	for i := 0; i<len(data); i++ {
		out[i].New(data[i])
	}
	return out,nil
}

//skel is a skeleton object, of the type you want to store data in.
//needs one value in each field
func RowData(skel reflect.Value, rows *sql.Rows) [][]interface{} {
	//f := skel.NumField()
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	v := *new([][]interface{})
	c := 0
	for rows.Next() {
		a := make([]interface{},0,len(columns))
		//a := new([]interface{})
		v = append(v,a)
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		for i, _ := range columns {	//_ was col
			val := values[i]
			
			skel_val := skel.Field(i).Interface()
			k := skel.Field(i).Kind()
			b, ok := val.([]byte)
			if (ok) {
				t := reflect.ValueOf(skel_val).Index(0).Kind()
				if k != reflect.String {
//De-tabbed to prevent line-wrapping
		if t == reflect.Int {
			sep := strings.Split(string(b[1:len(b)-1]),",")
			intarr := make([]int,len(sep))
			for j,s := range sep {
				intarr[j],_ = strconv.Atoi(s)
			}
			v[c] = append(v[c],intarr)
		}else if t == reflect.String {
			sep := strings.Split(string(b[1:len(b)-1]),",")
			strarr := make([]string,len(sep))
			v[c] = append(v[c],strarr)
		}else {
			fmt.Println("Else")
		}
//Re-tabbed
				} else {
					v[c] = append(v[c],string(b))
				}
			} else if k == reflect.Int64 {
				v[c] = append(v[c],val.(int64))
			} else if k == reflect.Bool {
				v[c] = append(v[c],val.(bool))
			}
		}
		c++
	}
	return v
}

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
//				fmt.Println(v.Kind())
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
//			fmt.Println(v.Kind())
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
	//fmt.Println("Unrolled: ")
	//fmt.Println(stro)
//	for i := range argo {
//		k := reflect.ValueOf(argo[i]).Kind()
		//fmt.Println(k,": ",argo[i])
//	}
//	fmt.Println(len(argo))
//	fmt.Println(argo)
	return stro,argo
}
