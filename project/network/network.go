package network

import (
  
)

func Sender(fsmCh config.FSMChannels, netCh config.NetworkChannels){
  //KJØRES SOM GOROUNTINE
  order := config.Order{}
  Msg := config.Packet{}
  for{
    select{
      case buttonpress := <- fsmCh.Drv_buttons: //Fått inn knappetrykk
        fmt.Println("Knapp er trykket! ", int(buttonpress.Button), buttonpress.Floor)

        //Vil sende dette knappetrykket ut!!
        
        Msg.Elev_ID = orderhandler.GetElevID()
        Msg.Elev_rank = orderhandler.GetElevRank()

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
        /*
        type Elevator struct{
          ElevID int
          ElevRank int
          CurrentOrder Order
           CurrentFloor int
          CurrentState int
        }*/

        netCh.TransmittElevStateCh <- newState

      case newCurrentOrder := <- fsmCh.New_current_order
        //newCurrentOrder er en Order struct:
        /*
          type Order struct{
            Floor int
            ButtonType int
            Type_action int //-1 hvis ordre skal slettes, 1 hvis ordre blir lagt til.
            Packet_id int
            Approved bool
            Receiver_elev int
            }
        */
        netCh.TransmittCurrentOrderCh <- newCurrentOrder


    }
  }
}



func Receiver(){

}