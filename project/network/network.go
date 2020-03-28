package network

import (
  
)

func Sender(fsmCh config.FSMChannels, netCh config.NetworkChannels, elevID int, elevatorMap map[int]*config.Elevator){
  //KJØRES SOM GOROUNTINE
  order := config.Order{}
  elevator := *elevatorMap[elevID]
  for{
    select{
      case buttonPress := <- fsmCh.Drv_buttons: //Fått inn knappetrykk
        fmt.Println("Knapp er trykket! Floor: ", buttonPress.Floor," buttonType: ", buttonPress.Button)

        order.Floor = buttonpress.Floor
        order.ButtonType = int(buttonpress.Button)
        order.Packet_id = rand.Intn(10000)
        order.Type_action = 1 //Det er en ordre som skal legges til
        order.Approved = false

        Msg.New_order = order

        netCh.TransmittOrderCh <- Msg
        fmt.Println("Har nå sendt avgårde pakke om at knapp er trykket!")



      case newState := <-fsmCh.New_state: //her må det opprettes ny intern channel. denne skal skrives Elevator til når det er en ny endring. Denne skal også sendes med en gang en heis går online.
        //newState er en Elevator Struct.
        netCh.TransmittElevStateCh <- newState

      case newCurrentOrder := <- fsmCh.New_current_order
        //newCurrentOrder er en Order struct.
        netCh.TransmittCurrentOrderCh <- newCurrentOrder


    }
  }
}



func Receiver(ch config.NetworkChannels, elevID int, elevatorMap map[int]*config.Elevator){

  for{
    select{
    case peerUpdate := <-ch.PeerUpdateCh:
      fmt.Printf("Peer update:\n")
      fmt.Printf("  Peers:    %q\n", p.Peers)
      fmt.Printf("  New:      %q\n", p.New)
      fmt.Printf("  Lost:     %q\n", p.Lost)

      //HAR DET KOMMET NOEN FLERE ELEVATORS TIL?
      for _, peerStr := range p.New{
        peerInt, _ := strconv.Atoi(peerStr)
        elevatorMap[peerInt].Active = true
      }
    
      //HAR VI MISTET NOEN ELEVATORS?
      if len(p.Lost) > 0{
        for _, peerStr := range p.Lost{
          peerInt, _ := strconv.Atoi(peerStr)
          elevatorMap[peerInt].Active = true
        }
      }

      //må jeg sende ut noe på elevState Channel nå??


    case receivedMsg := <-ch.ReceiveOrderCh:
      if receivedMsg.Elev_ID == elevID{
        break //drit i ordre fra deg selv.
      }

    case elevator := <-ch.ReceiveElevStateCh:
      *elevatorMap[elevator.Elev_ID] = elevator

    case newCurrentOrder := <-ch.ReceiveCurrentOrderCh:
    }
  }
}













