package main

        

import (
 "fmt"
 "net"
 "time"
)

//Connect udp

func main() {

        buffer := make([]byte, 1024)

        /*writeConn, err := net.Dial("udp", "10.100.23.147:20014")

        if err != nil {
                fmt.Println("feil med writeConn")
                return
        }
        */

        ServerAddr,err := net.ResolveUDPAddr("udp","10.100.23.129:20014")
        readConn, err := net.ListenUDP("udp",ServerAddr)

        if err != nil {
                fmt.Println("feil med readConn")
                return
        }
        defer readConn.Close()


        for{

                //writeConn.Write([]byte("hei verden"))
                length,err := readConn.Read(buffer)


                if err != nil {
                        fmt.Println("feil med read")
                        return
                }
                fmt.Println(string(buffer[0:length]))
                time.Sleep(100*time.Millisecond)

        }
}