package main

import(
  "project/network/udp"
)

func main(){
  var senderPacket Packet
  reciever, sender := Init(rPort, wPort)
  timeOut := make(chan bool)
  sendTimeOut(timeOut, sender, reciever)
  for i := 0; i < 10; i++ {
    sender<-senderPacket
  }
}
