package udp

import(
  "time"
  "fmt"
)

//Må flyttes
const rPort = 1234
const wPort = 5678

/*
func messaging(){
  reciever, sender := udp.Init(rPort, wPort)
  timeOut := make(chan bool) //ha en måte å se at det har tatt for lang tid før man har fått svar
  go sendTimeOut(timOut, sender)
  //Must finish
}
*/

//Function that lets a packet be sent up to ten times if no ack is recived
func SendTimeOut(timeOut chan bool, sender chan Packet, reciever chan Packet){
  var senderPacket Packet
  var recieverPacket Packet
  sender<-senderPacket
  switch {
  case senderPacket.message_nr < 10:
    reciever<-recieverPacket
    if (recieverPacket.ID == 0 || !(recieverPacket.ID==senderPacket.ID)) { //Funker dette??
      time.Sleep(500 * time.Millisecond)
      senderPacket.message_nr++
      fmt.Println("Retrying to send package. message_nr: ", senderPacket.message_nr)
      sender<-senderPacket
    }else{
      timeOut<-false
      fmt.Println("Packet sent!")
    }
case senderPacket.message_nr >= 10:
    timeOut<-true
    fmt.Println("Timeout. Packet used too long time to be sent.")
  }
}
