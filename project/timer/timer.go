package timer

import(
  "time"
  "fmt"
  "../config"
)



func DoorTimer(finished chan<- bool, start <-chan bool, doorOpenTime time.Duration) { //må kjøres som goroutine

	doorTimer := time.NewTimer(doorOpenTime)

  //empty the channel -> not concurrent receivers
	if !doorTimer.Stop() {
		<-doorTimer.C
	}

	for {
		select {
		case <-start:
			doorTimer.Reset(doorOpenTime)
		case <-doorTimer.C:
			finished <- true
		}
	}
}

func WatchDogTimer(fsmCh config.FSMChannels, netCh config.NetworkChannels, elevID int, elevatorMap map[int]*config.Elevator, watchDogTime time.Duration) {
	WatchDogTimer := time.NewTimer(watchDogTime)

  //empty the channel -> not concurrent receivers
	if !WatchDogTimer.Stop() {
		<-WatchDogTimer.C
	}

	for {
		select {
		case <-fsmCh.Drv_floors:
			//SJEKK OM DEN HAR EN CURRENT_ORDER!!! HVIS IKKE SKAL IKKE TIDEN RESETTES
			if elevatorMap[elevID].CurrentOrder.Floor != -1{
				WatchDogTimer.Reset(watchDogTime)
			}
		case <-WatchDogTimer.C:
			fmt.Println("WatchDog Released")
			elevatorMap[elevID].CurrentOrder.Floor = -1
			elevatorMap[elevID].Stuck = true
			go func(){fsmCh.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!
		}
	}


}