package main



import (
 "fmt"
 "net"
 "time"
 //"bufio"
 //"strconv"
)

func checkError(err error){
	if err != nil {
		fmt.Println(err)
		//return
	}
}

func main() {
	SERVER_IP := "10.100.23.147"
	SERVER_PORT := "33546"
	LOCAL_IP := "10.100.23.129"
	LOCAL_PORT := "20014"


	server_addr,err := net.ResolveTCPAddr("tcp", SERVER_IP+":"+SERVER_PORT)
	checkError(err)

	server_conn, err := net.DialTCP("tcp", nil, server_addr)
	checkError(err)

	local_addr,err := net.ResolveTCPAddr("tcp",LOCAL_IP+":"+LOCAL_PORT)
	checkError(err)

	listener, err := net.ListenTCP("tcp",local_addr)
	checkError(err)

	connect_msg := "Connect to: " + LOCAL_IP + ":" + LOCAL_PORT + "\x00"
	_, err = server_conn.Write([]byte (connect_msg))
	checkError(err)

	client_conn, err := listener.AcceptTCP()
	checkError(err)

	i:=0
	for(i<10){
			msg := "Hello, this is g14, \x00"
			fmt.Println("Sending msg: ", msg,  "\n",i, "\n")
    	client_conn.Write([]byte(msg))

			buffer := make([]byte, 1024)
			client_conn.Read(buffer)
			fmt.Println("Msg recived: ", string(buffer), "\n")
    	time.Sleep(time.Second * 1)
			i++
	}
}
