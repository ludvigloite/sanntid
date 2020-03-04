package main

import(
  "./network/udp"
)

func main(){
    reciever, sender = network.Init("6789", "2345")
    network.Listen(reciever, "6228")
    network.Bcast(sender, "3345")


}
