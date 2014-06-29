package main

type Quest struct {
	Questid		int64
	Name		string
	Description	string
	Category	string
	Recurring	bool
	Xpvalue		int64
	Requiredquests	[]int
	Attributes	[]string
}
func (quest *Quest) New(vars []interface{}) {
	quest.Questid = vars[0].(int64)
	quest.Name = vars[1].(string)
	quest.Description = vars[2].(string)
	quest.Category = vars[3].(string)
	quest.Recurring = vars[4].(bool)
	quest.Xpvalue = vars[5].(int64)
	quest.Requiredquests = vars[6].([]int)
	quest.Attributes = vars[7].([]string)
}
type User struct {
	Username	string	`db:"Username"`
	Firstname	string	`db:"Firstname"`
	Lastname	string	`db:"Lastname"`
	Xp		int64	`db:"Xp"`
	Completedquests	[]int	`db:"Completedquests"`
	Attributes	[]string `db:"Attributes"`
}
func (user *User) New(vars []interface{}){
	user.Username = vars[0].(string)
	user.Firstname = vars[1].(string)
	user.Lastname = vars[2].(string)
	user.Xp = vars[3].(int64)
	user.Completedquests = vars[4].([]int)
	user.Attributes = vars[5].([]string)
}
