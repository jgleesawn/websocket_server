<script type='text/javascript'> 
ws.onmessage = function (event) {
	curDiv = addElement();
	document.getElementById(curDiv).innerHTML = event.data;
	reader.readAsBinaryString(event.data);
};
function get(){ 
	ws.send("get "+document.getElementById("name").value) 
};
function store(){ 
	ws.send("store "+document.getElementById("name").value+" "+document.getElementById("age").value); 
};
function getall(){ 
	ws.send("all");
};
</script>
<div id='input'>
name:<input type='text' id='name' name='name' value='oldman'>
age:<input type='text' id='age' name='age' value='132'>
<button onclick='get()'>Get</button>
<button onclick='store()'>Store</button>
<button onclick='getall()'>Entire Table</button>
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
var reader = new FileReader();
reader.onload = function(e) {
	console.log(reader.result);
}
</script>
