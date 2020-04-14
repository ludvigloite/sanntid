package network

import (
	"fmt"
  	"strconv"


  	"../config"
  	"../elevio"
)

func Sender(fsmCh config.FSMChannels, netCh config.NetworkChannels, elevID int, elevatorMap map[int]*config.Elevator){
  //KJØRES SOM GOROUNTINE
  order := config.Order{}
  order.Sender_elev_ID = elevID
  for{
    select{
      case buttonPress := <- fsmCh.Drv_buttons: //Fått inn knappetrykk
                
        order.Sender_elev_rank = elevatorMap[elevID].ElevRank
        order.Floor = buttonPress.Floor
        order.ButtonType = buttonPress.Button
        order.Should_add = true //Det er en ordre som skal legges til
        order.Receiver_elev = elevID

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
        

      case newState := <-fsmCh.New_state:
        //newState er en Elevator Struct.
        for i:=0;i<config.NUM_PACKETS_SENT;i++{
        	netCh.TransmittElevStateCh <- newState
        } 
        

      case newCurrentOrder := <- fsmCh.New_current_order:
        //newCurrentOrder er en Order struct.
        for i:=0;i<config.NUM_PACKETS_SENT;i++{
        	netCh.TransmittCurrentOrderCh <- newCurrentOrder
        } 


      case floor := <-fsmCh.Stopping_at_floor:

        order.Sender_elev_rank = elevatorMap[elevID].ElevRank
        order.Floor = floor
        order.Should_add = false //Det er en ordre som skal fjernes
        order.Receiver_elev = elevID

        if elevatorMap[elevID].HallOrders[floor][elevio.BT_HallUp]{ //Sender bare RemoveOrder om det er en ordre
          order.ButtonType = elevio.BT_HallUp
          elevatorMap[elevID].HallOrders[floor][elevio.BT_HallUp] = false
          	for i:=0;i<config.NUM_PACKETS_SENT;i++{
          		netCh.TransmittOrderCh <- order
        	} 
        }
        if elevatorMap[elevID].HallOrders[floor][elevio.BT_HallDown]{ //Sender bare RemoveOrder om det er en ordre
          	order.ButtonType = elevio.BT_HallDown
          	elevatorMap[elevID].HallOrders[floor][elevio.BT_HallDown] = false
          	for i:=0;i<config.NUM_PACKETS_SENT;i++{
        		netCh.TransmittOrderCh <- order
        	} 
        }
        if elevatorMap[elevID].CabOrders[floor]{ //Sender bare RemoveOrder om det er en ordre
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

  cabOrdersBackup := make(map[string][config.NUM_FLOORS]bool) //måtte være indeksert på string for å sendes med JSON

  for{
    select{
    case p := <-netCh.PeerUpdateCh:
      fmt.Printf("Peer update:\n")
      fmt.Printf("  Peers:    %q\n", p.Peers)
      fmt.Printf("  New:      %q\n", p.New)
      fmt.Printf("  Lost:     %q\n", p.Lost)

      //GÅR GJENNOM ALLE SOM ER ACTIVE
      for _, peerStr := range p.Peers{
        peerInt, _ := strconv.Atoi(peerStr)
        elevatorMap[peerInt].Active = true
      }

      if len(p.Peers) == 0{ //network failure!!
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


      //HAR DET KOMMET NOEN FLERE ELEVATORS TIL?
      if len(p.New) > 0{
        fsmCh.New_state <- *elevatorMap[elevID] //om det kommer noen nye må du sende ut deg selv sånn at de kan legge deg til!
        for i:=0;i<config.NUM_PACKETS_SENT;i++{
        	netCh.TransmittCabOrderBackupCh <- cabOrdersBackup //sender cabOrderBackup slik at alle vet om det.
        } 
        SendOrdersWhenComeback(netCh, elevatorMap, p.New, elevID, cabOrdersBackup)
      }
    

      //HAR VI MISTET NOEN ELEVATORS?
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
        fsmCh.New_state <- *elevatorMap[elevID] //Alle må vite om din nye rank
      }


    case newCabOrderBackup := <-netCh.ReceiveCabOrderBackupCh:
      cabOrdersBackup = MergeCaborders(cabOrdersBackup, newCabOrderBackup)


    case receivedOrder := <-netCh.ReceiveOrderCh:
    	fmt.Println("Har fått inn en ordre. Should_add: ",receivedOrder.Should_add, " . Etasje: ",receivedOrder.Floor)
      
      if receivedOrder.ButtonType == elevio.BT_Cab{ //Mottar cab orders fra andre heiser kun når du selv har vært nede..
      	if receivedOrder.Receiver_elev == elevID{
      		if elevatorMap[elevID].HasRecentlyBeenDown{
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
        if !receivedOrder.Should_add && elevatorMap[elevID].CurrentOrder.Floor == receivedOrder.Floor && elevatorMap[elevID].CurrentOrder.ButtonType != elevio.BT_Cab{
      		elevatorMap[elevID].CurrentOrder.Floor = -1
      		fmt.Println("Min CurrentOrder fjernes!!")
      	}
      }

      go func(){fsmCh.New_state <- *elevatorMap[elevID]}()


    case elevator := <-netCh.ReceiveElevStateCh:
      if elevator.ElevID != elevID{
        *elevatorMap[elevator.ElevID] = elevator
        fmt.Println("Har fått en ny state fra NR ",elevator.ElevID," hvor CurrentOrderFloor: ", elevator.CurrentOrder.Floor)
      }

    case newCurrentOrder := <-netCh.ReceiveCurrentOrderCh:
    	fmt.Println("Jeg har fått ny CurrentOrder")
      elevatorMap[newCurrentOrder.Receiver_elev].CurrentOrder = newCurrentOrder



    }
  }
}

func SendOrdersWhenComeback(netCh config.NetworkChannels, elevatorMap map[int]*config.Elevator, comebackElev string, senderElev int, cabOrdersBackup map[string][config.NUM_FLOORS]bool){
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

func MergeCaborders(cabOrders1 map[string][config.NUM_FLOORS]bool, cabOrders2 map[string][config.NUM_FLOORS]bool) map[string][config.NUM_FLOORS]bool{
  cabOrders := make(map[string][config.NUM_FLOORS]bool)
  var list [config.NUM_FLOORS]bool
  i_str := ""
  for i := 1; i < config.NUM_ELEVATORS+1; i++{
    i_str = strconv.Itoa(i)
    for j := 0; j < config.NUM_FLOORS; j++{
      list[j] = cabOrders1[i_str][j] || cabOrders2[i_str][j] //hvis en av dem blir true returnerer vi true
    }
    cabOrders[i_str] = list
  }
  return cabOrders
}