package network

import (
  "../config"
  //"../elevio"
  "fmt"
  "math/rand"
  "strconv"
  "../elevcontroller"
  
)

func Sender(fsmCh config.FSMChannels, netCh config.NetworkChannels, elevID int, elevatorMap map[int]*config.Elevator){
  //KJØRES SOM GOROUNTINE
  order := config.Order{}
  order.Sender_elev_ID = elevID
  for{
    select{
      case buttonPress := <- fsmCh.Drv_buttons: //Fått inn knappetrykk
        
        fmt.Println("Knapp er trykket! Floor: ", buttonPress.Floor," buttonType: ", buttonPress.Button)
        
        order.Sender_elev_rank = elevatorMap[elevID].ElevRank
        order.Floor = buttonPress.Floor
        order.ButtonType = buttonPress.Button
        order.Packet_id = rand.Intn(10000)
        order.Should_add = true //Det er en ordre som skal legges til
        //order.Approved = false

        netCh.TransmittOrderCh <- order
        fmt.Println("Har nå sendt avgårde pakke om at knapp er trykket!")

        

      case newState := <-fsmCh.New_state: //her må det opprettes ny intern channel. denne skal skrives Elevator til når det er en ny endring. Denne skal også sendes med en gang en heis går online.
        //newState er en Elevator Struct.
        netCh.TransmittElevStateCh <- newState
        fmt.Println("Sendt ny Elevstate ut")

      case newCurrentOrder := <- fsmCh.New_current_order:
        //newCurrentOrder er en Order struct.
        netCh.TransmittCurrentOrderCh <- newCurrentOrder

      case floor := <-fsmCh.Stopping_at_floor:

        order.Sender_elev_rank = elevatorMap[elevID].ElevRank
        order.Floor = floor
        order.Packet_id = rand.Intn(10000)
        order.Should_add = false //Det er en ordre som skal legges til
        //order.Approved = false

        netCh.TransmittOrderCh <- order
    }
  }
}



func Receiver(ch config.NetworkChannels, fsmCh config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator, activeElevators map[int]bool){

  for{
    select{
    case p := <-ch.PeerUpdateCh:
      fmt.Printf("Peer update:\n")
      fmt.Printf("  Peers:    %q\n", p.Peers)
      fmt.Printf("  New:      %q\n", p.New)
      fmt.Printf("  Lost:     %q\n", p.Lost)

      //HAR DET KOMMET NOEN FLERE ELEVATORS TIL?


      for _, peerStr := range p.Peers{
        fmt.Println(peerStr)
        peerInt, _ := strconv.Atoi(peerStr)
        fmt.Println(peerInt)
        elevatorMap[peerInt].Active = true
        activeElevators[peerInt] = true
      }
    
      //HAR VI MISTET NOEN ELEVATORS?
      if len(p.Lost) > 0{
        for _, peerStr := range p.Lost{
          peerInt, _ := strconv.Atoi(peerStr)
          elevatorMap[peerInt].Active = false
          activeElevators[peerInt] = false
        }
      }

      //HVIS MASTER DÆVVER MÅ NOE SPENNENDE SKJE.

      //må jeg sende ut noe på elevState Channel nå??


    case receivedOrder := <-ch.ReceiveOrderCh:
      /*if receivedOrder.Sender_elev_ID == elevID || (receivedOrder.Sender_elev_rank != 1 && *elevatorMap[elevID].ElevRank != 1){
        break //drit i ordre HVIS enten 1.Ordre kommer fra deg selv. 2. Slave prøver å sende til slave.
      }*/
      //hvis melding kommer fra master skal den merkes med Approved og sendes tilbake.
      //hvis melding kommer til master skal den 1. hvis approved: LAGRES 2. hvis ikke approved, sendes rett ut igjen
      fmt.Println("Mottar melding om ny ordre")
      elevatorMap[elevID].HallOrders[receivedOrder.Floor][receivedOrder.ButtonType] = receivedOrder.Should_add
      go func(){fsmCh.New_state <- *elevatorMap[elevID]}()


    case elevator := <-ch.ReceiveElevStateCh:
      fmt.Println("Mottatt ny elevState")
      //elevcontroller.PrintElevator(elevator)
      *elevatorMap[elevator.ElevID] = elevator
      fmt.Println("HEIS fra RECEIVE ELEV STATE:")
      fmt.Println(elevatorMap)
      elevcontroller.PrintElevator(*elevatorMap[elevator.ElevID])

    case newCurrentOrder := <-ch.ReceiveCurrentOrderCh:
      fmt.Println("Mottatt ny currentOrder")
      fmt.Println(newCurrentOrder)
      elevatorMap[newCurrentOrder.Receiver_elev].CurrentOrder = newCurrentOrder
    }
  }
}













