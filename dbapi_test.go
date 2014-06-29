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

	q.Name = "update"
	succ = db.UpdateQuest(q) 
	if succ == false {
		t.Fail()
	}

	u.Firstname = "update"
	u.Completedquests = []int{3000,200,0}
	succ = db.UpdateUser(u)
	if succ == false {
		t.Fail()
	}

	name,err := db.GetUser("testing_username")
	if err != nil {
		t.Fail()
	}
	fmt.Println(name)
	quest,err := db.GetQuest(1)
	if err != nil {
		t.Fail()
	}
	fmt.Println(quest)

}
