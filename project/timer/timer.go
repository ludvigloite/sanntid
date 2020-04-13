package timer

import(
  "time"
  "fmt"


  "../config"
)

func DoorTimer(finished chan<- bool, start <-chan bool, doorOpenTime time.Duration) {

	doorTimer := time.NewTimer(doorOpenTime)

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

func WatchDogTimer(fsmCh config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator, watchDogTime time.Duration) {
	WatchDogTimer := time.NewTimer(watchDogTime)

	if !WatchDogTimer.Stop() && elevatorMap[elevID].CurrentOrder.Floor != -1 && elevatorMap[elevID].CurrentFsmState == config.ACTIVE {
		<-WatchDogTimer.C
	}

	go func(){ //oppdaterer WatchDog hvert minutt så lenge den er IDLE eller DOOR_OPEN
		for{
			if elevatorMap[elevID].CurrentFsmState != config.ACTIVE{
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

		case <-WatchDogTimer.C:
			fmt.Println("WatchDog Released")
			elevatorMap[elevID].CurrentOrder.Floor = -1
			elevatorMap[elevID].Stuck = true
			go func(){fsmCh.New_state <- *elevatorMap[elevID]}()			
		}
	}
}