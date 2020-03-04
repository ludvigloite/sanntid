package main

import "net"
import "fmt"
import "bufio" 

func main() {

  fmt.Println("Launching server...")
 

  // listen 
  listen, _ := net.Listen("tcp", ":33546")

  // accept connection 
  connection, _ := listen.Accept()

  // run loop forever (or until ctrl-c)
  for {
    // will listen for message to process ending in newline (\n)
    message, _ := bufio.NewReader(connection).ReadString('\n')
    // output message received
    fmt.Print("Message Received:", string(message))
    // sample process for string received
    newmessage := "Received:" + message
    // send new string back to client
    connection.Write([]byte(newmessage + "\n"))
  }
}