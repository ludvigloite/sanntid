package comm

import(
  "time"
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
  go sendTimeOut(timOut, sender)
  //Must finish
}


//Function that lets a packet be sent up to ten times if no ack is recived
func sendTimeOut(timeOut chan bool, sender chan Packet, reciever chan Packet){
  var senderPacket Packet
  switch expression {
  sender<-senderPacket
  case senderPacket.message_nr < 10:
    if recieverPacket :<- reciever & senderPacket.ID == recieverPacket.ID { //Funker dette??
      time.Sleep(0.05 * time.Second)
      fmt.Println("Retrying to send package. message_nr: ", message_nr)
      senderPacket.message_nr++
      sender<-senderPacket
    }else{
      true->timeout
      fmt.Println("Packet delivered!")
    }
  case timeOut<-true:
    fmt.Println("Timeout. Packet used too long time to be sent.")
  }
}
