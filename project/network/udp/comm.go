package comm

import(
  "net"
  "fmt"
  "./udp"
)

//Må flyttes
const rPort := 1234
const wPort := 5678

func messaging(){
  reciever, sender := udp.Init(rPort, wPort)
  timeOut := make(chan bool) //ha en måte å se at det har tatt for lang tid før man har fått svar
  go senderTimeOut(timOut, sender)
}


func senderTimeOut(timeOut chan bool, sender chan Packet){
  var packet Packet
  switch expression {
  case sender<-packet:
    if packet.packetID < 10{
      packet.packetID++
      packet->sender
    }else{
      true->timeout
    }
  case timeOut<-false:
    fmt.Println("Timeout. Packet used too long to send.")
  }
}

func ackTimeOut(timeOut chan bool, ack chan Packet){
  var packet Packet
  
}
