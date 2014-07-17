package main

import (
	"testing"
	"os"
	"fmt"
)

func TestSuite(t *testing.T) {
	os.Setenv("DATABASE_URL","postgres://tc:tc@localhost/test_db")
	db := OpenDB()
	
	var q *Quest
	q = &Quest{0,"Test","Testing Function","Test",false,0,[]int{0},[]string{""}}
	succ := db.AddQuest(q) 
	if succ == false {
		t.Fail()
	}
	
	var u *User
	u = &User{"testing_username","test","test",100,[]int{0},[]string{""}}
	succ = db.AddUser(u) 
	if succ == false {
		//t.Fail()
	}

	q.Questid = 0
	q.Name = "update"
	q.Requiredquests = []int{0,1}
	q.Attributes = []string{"updates","quest"}
	succ = db.UpdateQuest(q) 
	if succ == false {
		t.Fail()
	}

	u.Firstname = "update"
	u.Completedquests = []int{3000,200,0}
	u.Attributes = []string{"fast","strong"}
	succ = db.UpdateUser(u)
	if succ == false {
		t.Fail()
	}

	fmt.Println("Getting Data.")

	ret,err := db.GetUser("testing_username")
	ru := User(ret.([]User)[0])
	//fmt.Println(ru)
	if ru.Attributes[0] != "fast" && ru.Attributes[0] != "strong" {
		//fmt.Println(ru.Attributes)
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
	fmt.Println(ru)
	quest,err := db.GetQuest(0)
	if err != nil {
		t.Fail()
	}
	fmt.Println(quest)

}
