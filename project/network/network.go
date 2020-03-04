package network

import (
  "net"
  "set_up"
)
func UDPInit(localPort, bcastPort string){

  //broadcast address
  bcastAddr, err = ResolveUDPAddr("udp", "255.255.255.255" + bcastPort)
  if err != nil{
    return err
  }

  //local address
  localConn, err := net.DialUDP("udp4", nil, bcastAddr)
  if err != nil{
    return err
  }
  defer localConn
  localAddr, err := net.ResolveUDPAddr("udp", (localConn.LocalAddr()).String())
}
