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

        //fmt.Println("Knapp er trykket! Floor: ", buttonPress.Floor," buttonType: ", buttonPress.Button)

        order.Sender_elev_rank = elevatorMap[elevID].ElevRank
        order.Floor = buttonPress.Floor
        order.ButtonType = buttonPress.Button
        order.Packet_id = rand.Intn(10000)
        order.Should_add = true //Det er en ordre som skal legges til
        order.Receiver_elev = elevID

        netCh.TransmittOrderCh <- order
        //fmt.Println("Har nå sendt avgårde pakke om at knapp er trykket!")



      case newState := <-fsmCh.New_state:
        //newState er en Elevator Struct.
        netCh.TransmittElevStateCh <- newState
        //fmt.Println("Sendt ny Elevstate ut")

      case newCurrentOrder := <- fsmCh.New_current_order:
        //newCurrentOrder er en Order struct.
        netCh.TransmittCurrentOrderCh <- newCurrentOrder

      case floor := <-fsmCh.Stopping_at_floor:

        order.Sender_elev_rank = elevatorMap[elevID].ElevRank
        order.Floor = floor
        order.Packet_id = rand.Intn(10000)
        order.Should_add = false //Det er en ordre som skal fjernes
        order.Receiver_elev = elevID

        if elevatorMap[elevID].HallOrders[floor][elevio.BT_HallUp]{ //Sender bare RemoveOrder om det er en ordre
          order.ButtonType = elevio.BT_HallUp
          netCh.TransmittOrderCh <- order
        }

        if elevatorMap[elevID].HallOrders[floor][elevio.BT_HallDown]{ //Sender bare RemoveOrder om det er en ordre
          order.ButtonType = elevio.BT_HallDown
          netCh.TransmittOrderCh <- order
        }

        if elevatorMap[elevID].CabOrders[floor]{ //Sender bare RemoveOrder om det er en ordre
          order.ButtonType = elevio.BT_Cab
          netCh.TransmittOrderCh <- order
        }
    }
  }
}

func SendCabOrdersWhenComeback(netCh config.NetworkChannels, elevatorMap map[int]*config.Elevator, comebackElev int, senderElev int, cabOrdersBackup map[int][config.NUM_FLOORS]bool){
  order := config.Order{}
  order.Sender_elev_ID = senderElev
  order.Sender_elev_rank = elevatorMap[senderElev].ElevRank
  order.ButtonType = elevio.BT_Cab
  order.Packet_id = rand.Intn(10000)
  order.Should_add = true //Det er en ordre som skal legges til
  order.Receiver_elev = comebackElev

  for i := 0; i < config.NUM_FLOORS; i++{
    if cabOrdersBackup[comebackElev][i]{
      order.Floor = i
      netCh.TransmittOrderCh <- order
    }
  }
}



func Receiver(ch config.NetworkChannels, fsmCh config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator){

  cabOrdersBackup := make(map[int][config.NUM_FLOORS]bool)

  for{
    select{
    case p := <-ch.PeerUpdateCh:
      fmt.Printf("Peer update:\n")
      fmt.Printf("  Peers:    %q\n", p.Peers)
      fmt.Printf("  New:      %q\n", p.New)
      fmt.Printf("  Lost:     %q\n", p.Lost)

      //HAR DET KOMMET NOEN FLERE ELEVATORS TIL?
      for _, peerStr := range p.Peers{
        peerInt, _ := strconv.Atoi(peerStr)
        elevatorMap[peerInt].Active = true
        fsmCh.New_state <- *elevatorMap[elevID] //om det kommer noen nye må du sende ut deg selv sånn at de kan legge deg til!
        SendCabOrdersWhenComeback(ch, elevatorMap, peerInt, elevID, cabOrdersBackup)
        //Send dine hall_orders og lagrede Cab orders!
      }


      //HAR VI MISTET NOEN ELEVATORS?
      if len(p.Lost) > 0{
        for _, peerStr := range p.Lost{
          peerInt, _ := strconv.Atoi(peerStr)
          elevatorMap[peerInt].Active = false
          cabOrdersBackup[peerInt] = elevatorMap[peerInt].CabOrders
          elevatorMap[peerInt].CurrentOrder.Floor = -1

          
          NuActiveElevators := 0

          for _, elevator := range elevatorMap{
            if elevator.Active{
              NuActiveElevators++
            }
          }


          if elevatorMap[peerInt].ElevRank == 1 { //MASTER FIKSING
            //elevatorMap[peerInt].ElevRank = 3 //hardkodind, burde endres til NUMBER_ELEVATORS?
            //per naa saa vil en master som mistes faa rank 2 siden ranken forst settes til 3, men saa dekrementeres
            elevatorMap[elevID].ElevRank --
            elevatorMap[elevID].ElevRank = NuActiveElevators + 1//hardkodind, burde endres til NUMBER_ELEVATORS?
            fmt.Println("Jeg har nå rank ",elevatorMap[elevID].ElevRank)
          }
        }
      }



    case receivedOrder := <-ch.ReceiveOrderCh:
      /*if receivedOrder.Sender_elev_ID == elevID || (receivedOrder.Sender_elev_rank != 1 && *elevatorMap[elevID].ElevRank != 1){
        break //drit i ordre HVIS enten 1.Ordre kommer fra deg selv. 2. Slave prøver å sende til slave.
      }*/
      //hvis melding kommer fra master skal den merkes med Approved og sendes tilbake.
      //hvis melding kommer til master skal den 1. hvis approved: LAGRES 2. hvis ikke approved, sendes rett ut igjen

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
      //elevcontroller.PrintElevator(elevator)
      if elevator.ElevID != elevID{
        *elevatorMap[elevator.ElevID] = elevator
        //fmt.Println("MOTTATT HEIS:")
        //elevcontroller.PrintElevator(*elevatorMap[elevator.ElevID])
        fsmCh.LightUpdateCh <- true
      }

    case newCurrentOrder := <-ch.ReceiveCurrentOrderCh:
      //fmt.Println("Mottatt ny currentOrder")
      //fmt.Println(newCurrentOrder)
      elevatorMap[newCurrentOrder.Receiver_elev].CurrentOrder = newCurrentOrder
    }
  }
}
