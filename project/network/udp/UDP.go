/*
This module makes it possible to listen to packets sent over the channel
reciever using listen(...), and broadcast over the channel sender using
bcast(...). Init initializes the channels and the goroutines.
*/

package network

import(
  "bytes"
  "fmt"
  "net"
  "../bcast"
  "../localip"
  "../conn"
  ".../setup"
)

//Initialize channels for reciving and sending packets
//rPort is the port used to read
//wPort is the port used to write

func Init(rPort string, wPort string) (<-chan Packet, chan<- Packet){
  reciever := make(chan Packet, buffer)
  sender := make(chan Packet, buffer)
  go listen(reciever, rPort)
  go bcast(sender, wPort)
  return reciever, sender
}

func listen(reciever chan Packet, port string){
  //Set up connection to listen to
  localAddr, _ := net.ResolveUDPAddr("udp", port)
  conn, err := net.ListenUDP("udp", localAddr)
  if err != nil{
    return err
  }
  defer conn.Close()

  var packet Packet
  //Infinite loop waiting for reciving packets
  for {
    Receiver(int(port), reciever) //litt usikker her???
  }
}

func bcast(sender chan Packet, port string){
  localIP, err = LocalIP()
  if err != nil{
    return err
  }

  destinationAddr, _ = net.ResolveUDPAddr("udp", broadcastAddr + port)
  conn, err = net.DialUDP("udp", localIP, destinationAddr) //skal det være destinationAddr eller broadcastAddr vi prøver å kontakte???
  if err != nil{
    return err
  }
  defer conn.Close()

  //Infinite loop waiting for sending packets
  for{
    Transmitter(int(port), sender) //litt usikker her???
  }
}
