package main

import(
  "encoding/binary"
  "fmt"
  "net"
  "log"
  t "time"
  "os/exec"
)

var buffer = make([]byte, 16)
var counter uint64

func spawn_backup(){ //kjør phoenix.go i nytt terminalvindu.
	//file_path :=  "/home/student/Desktop/G14/exercise-6-gruppe-14/phoenix.go"

	//exec.Command("go","run", file_path)
	(exec.Command("gnome-terminal","-x","sh","-c","go run phoenix.go")).Run()
	fmt.Println("A backup has been spawned")
}

func main(){

  //Set up connection
  addr, err := net.ResolveUDPAddr("udp", "localhost:6789")
  conn, err := net.ListenUDP("udp", addr)
  dialConn,err := net.DialUDP("udp", nil, addr) //Kan kanskje ha denne på starten??

  isPrimary := false

  if err != nil{
    log.Println("Error: Couldn't establish connection.")
  }

	fmt.Println("I am a backup :-)")

  //Backup
  for !(isPrimary){ //while
    conn.SetReadDeadline(t.Now().Add(2*t.Second)) //Check if not getting a "respond" in 2.5 sec
    n, _, err := conn.ReadFromUDP(buffer)             //Read from buffer
    if err != nil{
      isPrimary = true                               //Primary down, backup is now set to primary
    } else{
      counter = binary.BigEndian.Uint64(buffer[:n]) //%Counter is set to last digit in buffer
    }
  }
  conn.Close()

  //Spawn Backup
  fmt.Println("I am now primary!")
  spawn_backup()

  for (isPrimary){
  	fmt.Println("counter: ",counter)
  	counter++
  	binary.BigEndian.PutUint64(buffer,counter) //Plasserer counter til bufferen
  	_,_ = dialConn.Write(buffer) //Skriver counteren til UDP
  	t.Sleep(500*t.Millisecond) //sov i et halvt sekund

  }
  dialConn.Close()


}
