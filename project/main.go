package main

import(
  "./network/udp"
  "./network/comm"
 
)

func main(){
	const rPort = 1234
	const wPort = 5678
	var senderPacket udp.Packet
	reciever, sender := udp.Init(rPort, wPort)
	timeOut := make(chan bool)
	comm.SendTimeOut(timeOut, sender, reciever)
	for i := 0; i < 10; i++ {
	sender<-senderPacket
	}
}
