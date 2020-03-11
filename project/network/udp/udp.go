/*
This module makes it possible to listen to packets sent over the channel
reciever using listen(...), and broadcast over the channel sender using
bcast(...). Init initializes the channels and the goroutines.
*/

package udp

import(
  "fmt"
  "strconv"
  "net"
  "../bcast"
  "../localip"
  //".../setup"
)

//Initialize channels for reciving and sending packets
//rPort is the port used to read
//wPort is the port used to write



//Burde legges over i setup, men finner ikke setup
const buffer = 1024
const broadcastAddr = "255.255.255.255"

//Trenger vi ha denne et annet sted???
type Packet struct{
  ID          int
  timestamp         int
  error_id          int
  state             int
  current_order     int
  message_nr        int
  order_list        [3][4]int
  confirmed_orders  [3][4]int
}


func Init(rPort int, wPort int) (<-chan Packet, chan<- Packet){
  reciever := make(chan Packet, buffer)
  sender := make(chan Packet, buffer)
  go listen(reciever, rPort)
  go broadcast(sender, wPort)
  return reciever, sender
}

func listen(reciever chan Packet, port int){
  //Set up connection to listen to
  localAddr, _ := net.ResolveUDPAddr("udp", strconv.Itoa(port))
  conn, err := net.ListenUDP("udp", localAddr)
  if err != nil{
    fmt.Println("Error while listening")
  }
  defer conn.Close()

  //Infinite loop waiting for reciving packets
  for {
    bcast.Receiver(port, reciever) //litt usikker her???
  }
}

func broadcast(sender chan Packet, port int){
  localIPString, err := localip.LocalIP()
  localIP := net.ParseIP(localIPString)
  var localIPAddr net.UDPAddr
  localIPAddr.IP = localIP

  if err != nil{
    fmt.Println("Error while setting up loacalIP")
  }

  destinationAddr, _ := net.ResolveUDPAddr("udp", broadcastAddr + strconv.Itoa(port))
  conn, err := net.DialUDP("udp", &localIPAddr, destinationAddr) //skal det være destinationAddr eller broadcastAddr vi prøver å kontakte???
  if err != nil{
    fmt.Println("Error while bcast")
  }
  defer conn.Close()

  //Infinite loop waiting for sending packets
  for{
    bcast.Transmitter(port, sender) //litt usikker her???
  }
}
