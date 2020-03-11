package main

import(
  "./network"

)

func main(){
  var senderPacket Packet
  reciever, sender := Init(rPort, wPort)
  timeOut := make(chan bool)
  SendTimeOut(timeOut, sender, reciever)
  for i := 0; i < 10; i++ {
    sender<-senderPacket
  }
}
