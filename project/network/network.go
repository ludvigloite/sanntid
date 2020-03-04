package network

import (
  "net"
  "../setup"
)
func UDPInit(localPort, bcastPort string){

  //Setting up broadcast address
  bcastAddr, err = ResolveUDPAddr("udp", broadcastAddr + bcastPort)
  if err != nil{
    return err
  }

  //Setting up broadcast listening conn
  bcastListenConn, err = net.ListenUDP("udp",bcastAddr)
  if err != nil{
    return err
  }

  //Setting up local address
  localConn, err := net.DialUDP("udp4", nil, bcastAddr)
  if err != nil{
    return err
  }

  localAddr, err := net.ResolveUDPAddr("udp", (localConn.LocalAddr()).String())
  localAddr.Port() = localPort
  defer localConn

  if err != nil{
    return err
  }

  //Setting up local listening conn
  bcastListenConn, err := net.ListenUDP("udp",localAddr)
  if err != nil{
    return err
  }
}
