package acknowlegde

import (
  "./network"
  "time"
)

const Slave_Timeout = 5 * time.Second

func waitingForBackup () {

}

func listenForClientTimeout(timer *time.Timer, id network.ID, timeout chan network.id){
  for  {
    select {
      case <- timer.C:
        timeout <- id
    }
  }
}

func addNewOrders(request,) {
  for _,r = range(request){


  }

}

func deleteDoneOrders() {
  for _,r = range(request){

  }
}
