<script type='text/javascript'> 
ws.onmessage = function (event) {
	curDiv = addElement();
	document.getElementById(curDiv).innerHTML = event.data;
	/*if (reader.readyState == 1){
		queue.push(event.data);
	} else {
		reader.readAsBinaryString(event.data);
	}*/
};
function get(){ 
	ws.send("get "+document.getElementById("name").value) 
};
function store(){ 
	ws.send("store "+document.getElementById("name").value+" "+document.getElementById("age").value); 
};
function getallUsers(){ 
	ws.send("get all Users; a");
};
function getallQuests(){ 
	ws.send("get all Quests; a");
};
function prepUser(){
	uname = '"'+document.getElementById('username').value+'"'
	ufname = '"'+document.getElementById('firstname').value+'"'
	ulname = '"'+document.getElementById('lastname').value+'"'
	str = '{"Username":'+uname+',"Firstname":'+ufname+',"Lastname":'+ulname+',"Xp":0,"Completedquests":[0],"Attributes":[""]}'
	return str
}
function AddUser(){
	ws.send("add User;"+prepUser())
}
function UpdateUser(){
	ws.send("update User;"+prepUser())
}
function prepQuest(){
	qid = parseInt(document.getElementById('qid').value)
	qname = '"'+document.getElementById('questname').value+'"'
	qdesc = '"'+document.getElementById('description').value+'"'
	qcat = '"'+document.getElementById('category').value+'"'
	qrec = document.getElementById('recurringtrue').value == 'true' ? true : false
	qxp = parseInt(document.getElementById('Xp').value)
	qreq = '['+reqquest+']'
	qattr = '['+questattr+']'

	str = '{"Questid":'+qid+',"Name":'+qname+',"Description":'+qdesc+',"Category":'+qcat+',"Recurring":'+qrec+',"Xp":'+qxp+',"Requiredquests":'+qreq+',"Attributes":'+qattr+'}'
	return str
}
function AddQuest(){
	ws.send("add Quest;"+prepQuest())
}
function UpdateQuest(){
	ws.send("update Quest;"+prepQuest())
}
function AppendQuestReq() {
	var elem = document.getElementById("reqquest")
	reqquest.push(parseInt(elem.value))
}
function AppendAttribute() {
	var elem = document.getElementById("attribute")
	questattr.push('"'+elem.value+'"')
}
var questattr = []
var reqquest = [0]
</script>
<div id='input'>
User: <br>
username<input type='text' id='username' name='username' value='username'><br>
firstname<input type='text' id='firstname' name='firstname' value='firstname'><br>
lasname<input type='text' id='lastname' name='lastname' value='lastname'><br>
<button onclick='AddUser()'>Add User</button><button onclick='UpdateUser()'>Update User</button><br>

Quest: <br>
qid<input type='number' id='qid' name='qid' min='0' max='10000'> <br>
questname<input type='text' id='questname' name='questname' value='questname'><br>
description<input type='text' id='description' name='description' value='description'><br>
category<input type='text' id='category' name='category' value='category'><br>
recurring<br><input type='radio' id='recurringtrue' name='recurring' value='true'>true<br>
<input type='radio' id='recurringfalse' name='recurring' value='false'>false<br>
xp<input type='number' id='Xp' name='Xp' min='0' max='10000'> <br>
add required quest<input type='number' id='reqquest' name='reqquest' min='0' max='10000'><button onclick='AppendQuestReq()'>Add Req Quest </button><br>
add attribute<input type='text' id='attribute' name='attribute' value='attribute'><button onclick='AppendAttribute()'>Add Attribute</button><br>
<button onclick='AddQuest()'>Add Quest</button><button onclick='UpdateQuest()'>Update Quest</button> <br>


<button onclick='store()'>Store</button>
<button onclick='getallUsers()'>All Users</button>
<button onclick='getallQuests()'>All Quests</button>
<button onclick='removeElements()'>Clear</button>
</div>
<div id='output'></div>
<script type='text/javascript'> 
function addElement() {
	var ni = document.getElementById('output');
	var newdiv = document.createElement('div');
	var div_id = Math.random().toString(36).substring(7);
	newdiv.setAttribute('id',div_id);
	ni.appendChild(newdiv);
	return div_id;};
function removeElements() {
	var out = document.getElementById('output');
      	for (i = out.childElementCount-1;i>=0;i--) {
		out.removeChild(out.childNodes[i])
	};
};
/*var reader = new FileReader();
var queue = [];
reader.onload = function(e) {
	document.getElementById(curDiv).innerHTML += reader.result;
	if (queue.length == 0)  return ;
	var i = queue.shift()
	console.log(i)
	this.readAsBinaryString(i);
}
*/
</script>
