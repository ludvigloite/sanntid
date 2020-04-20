//This module contains the functions Sender() and Receiver() which are both ran as goroutines. It also has a couple of help-functions. 
package network

import (
	"fmt"
  	"strconv"


  	"../config"
  	"../elevio"
)


func Sender(fsmCh config.FSMChannels, netCh config.NetworkChannels, elevID int, elevatorMap map[int]*config.Elevator){
  order := config.Order{}
  order.Sender_elev_ID = elevID

  for{
    select{
      case buttonPress := <- fsmCh.Drv_buttons:
                
        order.Sender_elev_rank = elevatorMap[elevID].ElevRank
        order.Floor = buttonPress.Floor
        order.ButtonType = buttonPress.Button
        order.Should_add = true //You should ADD this order, not delete it.
        order.Receiver_elev = elevID //if it is cab order, this is necessary.

        if buttonPress.Button == elevio.BT_Cab{
        	elevatorMap[elevID].CabOrders[buttonPress.Floor] = true
        	fsmCh.LightUpdateCh <- true
        }else if config.ADD_HALL_ORDERS_WHEN_NETWORK_DOWN{
        	elevatorMap[elevID].HallOrders[buttonPress.Floor][buttonPress.Button] = true
        	fsmCh.LightUpdateCh <- true
        }

        for i:=0;i<config.NUM_PACKETS_SENT;i++{
        	netCh.TransmittOrderCh <- order
        }        
        

      case newState := <-fsmCh.New_state: //New_state contains config.Elevator structs
        for i:=0;i<config.NUM_PACKETS_SENT;i++{
        	netCh.TransmittElevStateCh <- newState
        } 
        

      case newCurrentOrder := <- fsmCh.New_current_order: //New_current_order contains config.Order structs
        for i:=0;i<config.NUM_PACKETS_SENT;i++{
        	netCh.TransmittCurrentOrderCh <- newCurrentOrder
        } 


      case floor := <-fsmCh.Stopping_at_floor:

        order.Sender_elev_rank = elevatorMap[elevID].ElevRank
        order.Floor = floor
        order.Should_add = false //The order should be REMOVED
        order.Receiver_elev = elevID

        if elevatorMap[elevID].HallOrders[floor][elevio.BT_HallUp]{ //if there is no order, no need to delete it.
          order.ButtonType = elevio.BT_HallUp
          elevatorMap[elevID].HallOrders[floor][elevio.BT_HallUp] = false
          	for i:=0;i<config.NUM_PACKETS_SENT;i++{
          		netCh.TransmittOrderCh <- order
        	} 
        }
        if elevatorMap[elevID].HallOrders[floor][elevio.BT_HallDown]{
          	order.ButtonType = elevio.BT_HallDown
          	elevatorMap[elevID].HallOrders[floor][elevio.BT_HallDown] = false
          	for i:=0;i<config.NUM_PACKETS_SENT;i++{
        		netCh.TransmittOrderCh <- order
        	} 
        }
        if elevatorMap[elevID].CabOrders[floor]{
          	order.ButtonType = elevio.BT_Cab
          	elevatorMap[elevID].CabOrders[floor] = false
          	for i:=0;i<config.NUM_PACKETS_SENT;i++{
        		netCh.TransmittOrderCh <- order
        	} 
        }
        fsmCh.LightUpdateCh <- true
    }
  }
}


func Receiver(fsmCh config.FSMChannels, netCh config.NetworkChannels, elevID int, elevatorMap map[int]*config.Elevator){

  cabOrdersBackup := make(map[string][config.NUM_FLOORS]bool)

  for{
    select{
    case p := <-netCh.PeerUpdateCh:
      fmt.Printf("Peer update:\n")
      fmt.Printf("  Peers:    %q\n", p.Peers)
      //fmt.Printf("  New:      %q\n", p.New)
      //fmt.Printf("  Lost:     %q\n", p.Lost)
      fmt.Println()

      for _, peerStr := range p.Peers{
        peerInt, _ := strconv.Atoi(peerStr)
        elevatorMap[peerInt].Active = true
      }

      if len(p.Peers) == 0{
      	fmt.Println("Network Lost!!")
      	elevatorMap[elevID].NetworkDown = true
      	fsmCh.LightUpdateCh <- true
      	for i := 1; i < config.NUM_ELEVATORS+1; i++{
      		elevatorMap[i].CurrentOrder.Floor = -1
      		elevatorMap[i].Active = false
      	}
      	elevatorMap[elevID].Active = true
      }else{
      	elevatorMap[elevID].NetworkDown = false
      }

      if len(p.New) > 0{
        fsmCh.New_state <- *elevatorMap[elevID] //If it is detected another elevator, you must send your own elevator struct, so that it can add it to its local elevatorMap.
        for i:=0;i<config.NUM_PACKETS_SENT;i++{
        	netCh.TransmittCabOrderBackupCh <- cabOrdersBackup //Also send cabOrdersBackup so that everyone knows everyones cabOrders.
        } 
        sendOrdersWhenComeback(netCh, elevatorMap, p.New, elevID, cabOrdersBackup)
      }
   
      if len(p.Lost) > 0{
        for _, peerStr := range p.Lost{
          peerInt, _ := strconv.Atoi(peerStr)
          if peerInt == elevID{
          	break
          }

          elevatorMap[peerInt].Active = false
          cabOrdersBackup[peerStr] = elevatorMap[peerInt].CabOrders
          elevatorMap[peerInt].CurrentOrder.Floor = -1
          fsmCh.New_current_order <- config.Order{Sender_elev_ID: elevID, Floor: -1, Receiver_elev: peerInt}
        }
        fsmCh.New_state <- *elevatorMap[elevID]
      }


    case newCabOrderBackup := <-netCh.ReceiveCabOrderBackupCh:
      cabOrdersBackup = mergeCaborders(cabOrdersBackup, newCabOrderBackup)


    case receivedOrder := <-netCh.ReceiveOrderCh:
      
      if receivedOrder.ButtonType == elevio.BT_Cab{
      	if receivedOrder.Receiver_elev == elevID{
      		if elevatorMap[elevID].HasRecentlyBeenDown{ //Recieve cab orders from others only when you have been down.
      			elevatorMap[elevID].CabOrders[receivedOrder.Floor] = receivedOrder.Should_add
      			fsmCh.LightUpdateCh <- true
      		}
        }else{
        	elevatorMap[receivedOrder.Receiver_elev].CabOrders[receivedOrder.Floor] = receivedOrder.Should_add
        }        
      }else{
        if elevatorMap[elevID].HallOrders[receivedOrder.Floor][receivedOrder.ButtonType] != receivedOrder.Should_add{
          elevatorMap[elevID].HallOrders[receivedOrder.Floor][receivedOrder.ButtonType] = receivedOrder.Should_add
          fsmCh.LightUpdateCh <- true
        }
        if !receivedOrder.Should_add && elevatorMap[elevID].CurrentOrder.Floor == receivedOrder.Floor && elevatorMap[elevID].CurrentOrder.ButtonType != elevio.BT_Cab && receivedOrder.Sender_elev_ID != elevID{
      		elevatorMap[elevID].CurrentOrder.Floor = -1
      	}
      }

      fsmCh.New_state <- *elevatorMap[elevID]


    case elevator := <-netCh.ReceiveElevStateCh:
      if elevator.ElevID != elevID{
        *elevatorMap[elevator.ElevID] = elevator
      }

    case newCurrentOrder := <-netCh.ReceiveCurrentOrderCh:
      elevatorMap[newCurrentOrder.Receiver_elev].CurrentOrder = newCurrentOrder



    }
  }
}

//Sends backup of caborders when elevator has been down
func sendOrdersWhenComeback(netCh config.NetworkChannels, elevatorMap map[int]*config.Elevator, comebackElev string, senderElev int, cabOrdersBackup map[string][config.NUM_FLOORS]bool){
  comebackElevInt,_ := strconv.Atoi(comebackElev)
  order := config.Order{}
  order.Sender_elev_ID = senderElev
  order.Sender_elev_rank = elevatorMap[senderElev].ElevRank
  order.Should_add = true
  order.Receiver_elev = comebackElevInt

  for i := 0; i < config.NUM_FLOORS; i++{
    if cabOrdersBackup[comebackElev][i]{
      	order.Floor = i
      	order.ButtonType = elevio.BT_Cab
      	for i:=0;i<config.NUM_PACKETS_SENT;i++{
        	netCh.TransmittOrderCh <- order
        } 
    }
    for j := elevio.BT_HallUp; j != elevio.BT_Cab; j++{
      if elevatorMap[senderElev].HallOrders[i][j]{
        order.Floor = i
        order.ButtonType = j
        for i:=0;i<config.NUM_PACKETS_SENT;i++{
        	netCh.TransmittOrderCh <- order
        } 
      }
    }
  }
}

func mergeCaborders(cabOrders1 map[string][config.NUM_FLOORS]bool, cabOrders2 map[string][config.NUM_FLOORS]bool) map[string][config.NUM_FLOORS]bool{
  cabOrders := make(map[string][config.NUM_FLOORS]bool)
  var list [config.NUM_FLOORS]bool
  i_str := ""
  for i := 1; i < config.NUM_ELEVATORS+1; i++{
    i_str = strconv.Itoa(i)
    for j := 0; j < config.NUM_FLOORS; j++{
      list[j] = cabOrders1[i_str][j] || cabOrders2[i_str][j] //If one of them are true, return true
    }
    cabOrders[i_str] = list
  }
  return cabOrders
}