package main

import "net"
import "fmt"
import "bufio"
import "strconv"
import "time"

func main() {

  // connect to this socket
  connection, _ := net.Dial("tcp", "127.0.0.1:33546")
  i := 0
  for { 
    msg := strconv.Itoa(i)
        i++
    // send to socket
    connection.Write([]byte(msg + "\n"))
    time.Sleep(time.Second * 1)
    // listen for reply
    message, _ := bufio.NewReader(connection).ReadString('\n')
    fmt.Print("(FROM SERVER)"+message)
    
  }
}