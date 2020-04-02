//This module updates which peers that is alive and lost + new

package peers

import (
	"../conn"
	"fmt"
	"net"
	"sort"
	"time"
)

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

const interval = 15 * time.Millisecond
const timeout = 50 * time.Millisecond

// Transmitter(.) writes to conn after time interval if transmitEnable is enabeled
func Transmitter(port int, id string, transmitEnable <-chan bool) {
	//fmt.Println("Started transmitting peer at port ",port,". I have ID ",id)

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
	}
}

//Reciever(.)

func Receiver(port int, peerUpdateCh chan<- PeerUpdate) {
	//fmt.Println("Started receiving peer at port ",port)

	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval)) //sets deadline after interval
		n, _, _ := conn.ReadFrom(buf[0:]) //reads from buffer on conn

		id := string(buf[:n])

		// Adding new connection
		p.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// Removing dead connection
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			peerUpdateCh <- p
		}
	}
}
