

var net = require('net');
 
if(process.argv.length != 4){
console.log(process.argv.length);
	console.log("Usage: node %s <host> <port>", process.argv[1]);
	process.exit(1);	
}

var host=process.argv[2];
var port=process.argv[3];

if(host.length >253 || port.length >5 ){
	console.log("Invalid host or port. Try again!\nUsage: node %s <port>", process.argv[1]);
	process.exit(1);	
}





var readlineSync =require('readline-sync');

var user ={};

var validateRegex = /[0-9a-zA-Z]{5,}/;

var client = new net.Socket();
console.log("Simple telnet.js developed by Kartik Desai, SecAD");
console.log("Connecting to localhost:%s", port);
console.log("Connected to: %s:%s", host, port);
console.log("You need to login before sending/receiving messages");

client.connect(port,host, connected);






function connected(){

loginsync();
	
}
var counter = 0;


client.on("data",(data) => {

var Message = (data.toString());
var failedLogin = "Authentication_failed_Please_Try_Again! Invalid username or password";



if(Message!==failedLogin){

console.log("Received data:" + Message);
if(!counter){
console.log("Welcome to chat system. Type anything to send to public chat");
console.log("Type .userlist to see online users");
console.log("Type .exit to logout and close the connection");
console.log("Type pm for private chat");
callCommand()
}

counter++;


}
else {
console.log(failedLogin);
loginsync();
}



})

client.on("error",(err) => {
console.log("Error");
process.exit(2);

})

client.on("close",(data) => {
console.log("Connection has been disconnected");
process.exit(3);

})





function loginsync(){

user.username = readlineSync.question('Username ',{


});
user.password = readlineSync.question('Password ',{

hideEchoBack: true
});

client.write(JSON.stringify(user));




}


function callCommand(){

const keyboard = require('readline').createInterface({

input: process.stdin,
output:process.stdout


})

keyboard.on('line', (input) => {



if (input ===".exit"){

client.destroy();
console.log("Disconnected");
process.exit();

}


else if(input ===".userlist") {
var userList={}
userList.Type = "userlist";
userList.Message = "getUserList";
client.write(JSON.stringify(userList));

}


else if(input ==="pm") {

var privateChat={}
console.log("Now you can send message to specific user")
privateChat.Type = "private";
 keyboard.question('To ', (answer) => {
   privateChat.To = answer;
   keyboard.question('Message ', (answer) => {
   privateChat.Message = answer;
client.write(JSON.stringify(privateChat));
 })

})
  


}

else{
var publicChat={}

publicChat.Type = "public"
publicChat.Message  = input;
client.write(JSON.stringify(publicChat));
}

})




  }


       














