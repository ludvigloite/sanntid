package network

import (
  "../config"
  "../elevio"
  "fmt"
  "math/rand"
  "strconv"
  //"../elevcontroller"
  
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
        order.Packet_id = rand.Intn(10000)
        order.Should_add = true //Det er en ordre som skal legges til
        order.Receiver_elev = elevID

        for i:=0;i<config.NUM_PACKETS;i++{
        	netCh.TransmittOrderCh <- order
        }        
        

      case newState := <-fsmCh.New_state:
        //newState er en Elevator Struct.
        for i:=0;i<config.NUM_PACKETS;i++{
        	netCh.TransmittElevStateCh <- newState
        } 
        

      case newCurrentOrder := <- fsmCh.New_current_order:
        //newCurrentOrder er en Order struct.
        for i:=0;i<config.NUM_PACKETS;i++{
        	netCh.TransmittCurrentOrderCh <- newCurrentOrder
        } 


      case floor := <-fsmCh.Stopping_at_floor:

        order.Sender_elev_rank = elevatorMap[elevID].ElevRank
        order.Floor = floor
        order.Packet_id = rand.Intn(10000)
        order.Should_add = false //Det er en ordre som skal fjernes
        order.Receiver_elev = elevID

        if elevatorMap[elevID].HallOrders[floor][elevio.BT_HallUp]{ //Sender bare RemoveOrder om det er en ordre
          order.ButtonType = elevio.BT_HallUp
          	for i:=0;i<config.NUM_PACKETS;i++{
          		netCh.TransmittOrderCh <- order
        	} 
        }
        if elevatorMap[elevID].HallOrders[floor][elevio.BT_HallDown]{ //Sender bare RemoveOrder om det er en ordre
          	order.ButtonType = elevio.BT_HallDown
          	for i:=0;i<config.NUM_PACKETS;i++{
        		netCh.TransmittOrderCh <- order
        	} 
        }
        if elevatorMap[elevID].CabOrders[floor]{ //Sender bare RemoveOrder om det er en ordre
          	order.ButtonType = elevio.BT_Cab
          	for i:=0;i<config.NUM_PACKETS;i++{
        		netCh.TransmittOrderCh <- order
        	} 
        }

    }
  }
}

func SendOrdersWhenComeback(netCh config.NetworkChannels, elevatorMap map[int]*config.Elevator, comebackElev string, senderElev int, cabOrdersBackup map[string][config.NUM_FLOORS]bool){
  comebackElevInt,_ := strconv.Atoi(comebackElev)
  order := config.Order{}
  order.Sender_elev_ID = senderElev
  order.Sender_elev_rank = elevatorMap[senderElev].ElevRank
  order.Packet_id = rand.Intn(10000) //DETTE FUNKER IKKE!! ALLE PAKKENE FÅR SAMME PACKET_ID!
  order.Should_add = true //Det er en ordre som skal legges til
  order.Receiver_elev = comebackElevInt

  for i := 0; i < config.NUM_FLOORS; i++{
    if cabOrdersBackup[comebackElev][i]{
      	order.Floor = i
      	order.ButtonType = elevio.BT_Cab
      	for i:=0;i<config.NUM_PACKETS;i++{
        	netCh.TransmittOrderCh <- order
        } 
    }
    for j := elevio.BT_HallUp; j != elevio.BT_Cab; j++{
      if elevatorMap[senderElev].HallOrders[i][j]{
        order.Floor = i
        order.ButtonType = j
        for i:=0;i<config.NUM_PACKETS;i++{
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


func Receiver(ch config.NetworkChannels, fsmCh config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator){

  cabOrdersBackup := make(map[string][config.NUM_FLOORS]bool) //måtte være indeksert på string for å sendes med JSON

  for{
    select{
    case p := <-ch.PeerUpdateCh:
      fmt.Printf("Peer update:\n")
      fmt.Printf("  Peers:    %q\n", p.Peers)
      fmt.Printf("  New:      %q\n", p.New)
      fmt.Printf("  Lost:     %q\n", p.Lost)

      //GÅR GJENNOM ALLE SOM ER ACTIVE
      for _, peerStr := range p.Peers{
        peerInt, _ := strconv.Atoi(peerStr)
        elevatorMap[peerInt].Active = true
      }


      //HAR DET KOMMET NOEN FLERE ELEVATORS TIL?
      if len(p.New) > 0{
        fsmCh.New_state <- *elevatorMap[elevID] //om det kommer noen nye må du sende ut deg selv sånn at de kan legge deg til!
        for i:=0;i<config.NUM_PACKETS;i++{
        	ch.TransmittCabOrderBackupCh <- cabOrdersBackup //sender cabOrderBackup slik at alle vet om det.
        } 
        SendOrdersWhenComeback(ch, elevatorMap, p.New, elevID, cabOrdersBackup)
      }
    

      //HAR VI MISTET NOEN ELEVATORS?
      if len(p.Lost) > 0{

        for _, peerStr := range p.Lost{
          peerInt, _ := strconv.Atoi(peerStr)
          elevatorMap[peerInt].Active = false
          cabOrdersBackup[peerStr] = elevatorMap[peerInt].CabOrders
          elevatorMap[peerInt].CurrentOrder.Floor = -1
          fsmCh.New_current_order <- config.Order{Sender_elev_ID: elevID, Floor: -1, Receiver_elev: peerInt}


          //Istedet for det under kan jeg bare gå gjennom alle Active og så sette elevRank deretter.
          if elevatorMap[peerInt].ElevRank == 1{ //MASTER FIKSING
            elevatorMap[peerInt].ElevRank = 3
            elevatorMap[elevID].ElevRank --
          }else if elevatorMap[peerInt].ElevRank == 2 && elevatorMap[elevID].ElevRank == 3{
            elevatorMap[elevID].ElevRank = 2
          }
          fmt.Println("Jeg har nå rank ",elevatorMap[elevID].ElevRank)
        }
        fsmCh.New_state <- *elevatorMap[elevID] //Alle må vite om din nye rank
      }


    case newCabOrderBackup := <-ch.ReceiveCabOrderBackupCh:
      cabOrdersBackup = MergeCaborders(cabOrdersBackup, newCabOrderBackup)


    case receivedOrder := <-ch.ReceiveOrderCh:
      
      if receivedOrder.ButtonType == elevio.BT_Cab{
        elevatorMap[receivedOrder.Receiver_elev].CabOrders[receivedOrder.Floor] = receivedOrder.Should_add
        if receivedOrder.Receiver_elev == elevID{
           fsmCh.LightUpdateCh <- true
        }
      }else{
        if elevatorMap[elevID].HallOrders[receivedOrder.Floor][receivedOrder.ButtonType] != receivedOrder.Should_add{
          elevatorMap[elevID].HallOrders[receivedOrder.Floor][receivedOrder.ButtonType] = receivedOrder.Should_add
          fsmCh.LightUpdateCh <- true
        }
      }

      go func(){fsmCh.New_state <- *elevatorMap[elevID]}()


    case elevator := <-ch.ReceiveElevStateCh:
      if elevator.ElevID != elevID{
        *elevatorMap[elevator.ElevID] = elevator
      }

    case newCurrentOrder := <-ch.ReceiveCurrentOrderCh:
      elevatorMap[newCurrentOrder.Receiver_elev].CurrentOrder = newCurrentOrder



    }
  }
}













