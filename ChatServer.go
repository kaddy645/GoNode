
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

 
 
    
	
)

const BUFFERSIZE int = 1024

var buffer [BUFFERSIZE]byte

type User struct {
	Username string
	Login    bool
	Key      string
}

 type userInput struct {
		Type string
            To string
		Message string
	}


var allClient_conns = make(map[net.Conn]string)

var allLoggedIn_conns = make(map[net.Conn]interface{})
var lostclient = make(chan net.Conn)
var newclient = make(chan net.Conn)

var currentLoggedUser User
var currentLoggedUsername string
var userlist []string
var usernamelist string
var myconfig []User

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <port>\n", os.Args[0])
		os.Exit(0)
	}
	port := os.Args[1]
	if len(port) > 5 {
		fmt.Println("Invalid port value. Try again!")
		os.Exit(1)
	}

	server, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Printf("Cannot listen on port '" + port + "'!\n")
		os.Exit(2)
	}

	fmt.Println("ChatServer in GoLang developed by Kartik Desai, SecAD")

	fmt.Printf("ChatServer is listening on port '%s' ...\n", port)

	go func() {

		for {

			client_conn, _ := server.Accept()

			welcomemessage := fmt.Sprintf("A new client is connected from'%s' Waiting for login!\n", client_conn.RemoteAddr().String())

			fmt.Println(welcomemessage)

			go AuthenticateUser(client_conn)

		}
	}()
	for {
		select {
                   
		      case client_conn := <-newclient:
			allClient_conns[client_conn] = client_conn.RemoteAddr().String()
			allLoggedIn_conns[client_conn] = currentLoggedUsername
                   fmt.Println("No of Online Users", len(allLoggedIn_conns))

			if allLoggedIn_conns[client_conn] != "" {
                      
				go client_goroutine(client_conn)

			} 
                             case client_conn := <- lostclient:
                              
					go logout(client_conn)
                

		}

	}
}
func client_goroutine(client_conn net.Conn) {

	var buffer [BUFFERSIZE]byte


	messageForUser := fmt.Sprintf(" New user %s logged into Chat System  from %s.  %s  (from %d connections)", currentLoggedUsername, client_conn.RemoteAddr().String(), getUserList(),len(userlist))

		fmt.Println(messageForUser)
            sendtoAll([]byte (messageForUser))

	fmt.Printf("Connected Clients: %d\n", len(allLoggedIn_conns))

	go func() {
		for {
			byte_received, read_err := client_conn.Read(buffer[0:])
			if read_err != nil {
				lostclient <- client_conn
				return
			}
			fmt.Printf("Received data: %s\n", buffer[0:byte_received])

			handleUserRequest(client_conn, buffer[0:byte_received])

		}
	}()

}

func AuthenticateUser(client_conn net.Conn) {
	byte_received, read_err := client_conn.Read(buffer[0:])

	if read_err != nil {
		fmt.Println("Error in receiving...")
		lostclient <- client_conn
		return
	}
	fmt.Printf("Got data : %s Expecting Login Data\n", buffer[0:byte_received])
	status, Username, message := checklogin(buffer[0:byte_received])
 
      
	if status {
		
		currentLoggedUser = User{Username: Username, Login: true,Key:client_conn.RemoteAddr().String()}
		currentLoggedUsername = Username
        		     
		 fmt.Println(currentLoggedUser)
		newclient <- client_conn
           
		userlist = append(userlist, currentLoggedUser.Username)
            myconfig = append(myconfig, currentLoggedUser)
		usernamelist = usernamelist + ", " + currentLoggedUser.Username
       

	} else {
	failedLogin := fmt.Sprintf("Authentication_failed_Please_Try_Again! Invalid username or password")
	client_conn.Write([]byte(failedLogin))
      go AuthenticateUser(client_conn)
	
}

  fmt.Println(message)




	
}


func privateMsg(sender net.Conn, receiver string, msg string){
      counter := 0
	for client_conn,_ := range allLoggedIn_conns{
	        counter++
		if allLoggedIn_conns[client_conn] == receiver {
	
			recieving_user := client_conn
	

			incomingMsg := fmt.Sprintf("%s: %s",allLoggedIn_conns[sender], msg)
			
			sendtoOne(recieving_user, []byte(incomingMsg))
                  return
		} 
	}

       
           failedSending := fmt.Sprintf("Receiver is not online at the moment! Please Try Again!")
             sender.Write([]byte(failedSending))

}




func handleUserRequest(client_conn net.Conn, data []byte) {
   
	 var userIn userInput

       err := json.Unmarshal(data, &userIn)
	
     
	if err == nil &&  userIn.Type == "userlist" {
        
		client_conn.Write([]byte(getUserList()))
		return
	} 

      if err == nil &&  userIn.Type == "public" {
            fmt.Println("coming")
            publicMsg := []byte(userIn.Message)
		sendtoAll(publicMsg)
		return
	} 

      
      if err == nil &&  userIn.Type == "private" {
		privateMsg(client_conn, userIn.To, userIn.Message)
		return
	} 


       if err == nil &&  userIn.Type == "exit" {

		lostclient <- client_conn

		return

	}

}

func sendtoAll(data []byte) {

	for u, _ := range allLoggedIn_conns {
            
           	fmt.Printf("To All: %s\n", data)
		sendtoOne(u, data)
	}


}

func getUserList() string {

	var allUserList string
	for n, _ := range userlist {
		allUserList = allUserList + " " + userlist[n]
	}
	return "Online Users: " + allUserList

}

func sendtoOne(client_conn net.Conn, data []byte) {

	_, write_err := client_conn.Write(data)
			

	if write_err != nil {
		fmt.Println("DEBUG>Error in sending...to "+ client_conn.RemoteAddr().String())
		return
	}
}

func checkAccount(Username string, Password string) bool {

	users := []string{"k", "john", "smith", "jenny"}

	password := "123"

	for _, U := range users {
		if Username == U && Password == password {
			return true
		}
	}
	return false

}

func checklogin(data []byte) (bool, string, string) {

	type Account struct {
		Username string
		Password string
	}

	var account Account

	err := json.Unmarshal(data, &account)

	if err != nil || account.Username == " " || account.Password == " " {
		fmt.Printf("JSON parsing error : %s\n", err)
		return false, " ", `[BAD LOGIN] Expected:  {"Username":"..","Password":".."}`
	}

	fmt.Printf("DEBUG>Got: account =%s\n", account)
	fmt.Printf("DEBUG>Got: username=%s,password=%s\n", account.Username, account.Password)

	if checkAccount(account.Username, account.Password) {
		return true, account.Username, "logged"

	}

	return false, "", "Invalid username or password"

}





 func deleteDisconnedted(slices []string, name string,index int) []string{
  fmt.Println(index)  
  fmt.Println(len(slices))  
     
                    if(len(slices)>=index) {
	          
				if slices[index] == name   {
                             
					slices = append(slices[:index], slices[index+1:]...)
				}

}
				

		return slices
}


func getUsernameofLoggedout (client_conn net.Conn ) (string,int) {

  var username string
  var index int

  for i,user := range myconfig {
        if user.Key == client_conn.RemoteAddr().String() {
       fmt.Println("Hello")
       username  =  user.Username
       index  = i
        }
    }

     return username,index

}


func logout(client_conn net.Conn) ([]string){
	
		exitmsg := fmt.Sprintf("%s client disconned", allLoggedIn_conns[client_conn])
		go sendtoAll([]byte(exitmsg))
           
           username , i := getUsernameofLoggedout(client_conn)
       
	     userlist = deleteDisconnedted(userlist,username ,i )

		 fmt.Println(userlist)      


		go delete(allLoggedIn_conns, client_conn)
       
		
		client_conn.Close()

		return userlist

}














  
	
