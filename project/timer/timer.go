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
	if !WatchDogTimer.Stop() && elevatorMap[elevID].CurrentOrder.Floor != -1 && elevatorMap[elevID].CurrentState == config.ACTIVE {
		<-WatchDogTimer.C
	}

	go func(){
		for{
			if elevatorMap[elevID].CurrentState != config.ACTIVE{
				WatchDogTimer.Stop()
				WatchDogTimer.Reset(watchDogTime)
				time.Sleep(time.Second)
			}
		}
	 }()

	for {
		select {
		case <-fsmCh.Watchdog_updater:
			WatchDogTimer.Stop()
			WatchDogTimer.Reset(watchDogTime)
			//fmt.Println("WatchDog Reset")

		case <-WatchDogTimer.C:
			fmt.Println("WatchDog Released")
			elevatorMap[elevID].CurrentOrder.Floor = -1
			elevatorMap[elevID].Stuck = true
			go func(){fsmCh.New_state <- *elevatorMap[elevID]}()
					
		}
	}


}