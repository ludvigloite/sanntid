# Elevator project

## TLDR
- The task was to control 3 elevators.
- They must be able to communicate, give each other orders, update each others lights etc.
- The system must be able to handle most faults. This includes loss of network, loss of motor power, packet loss and more.
- The program was made spring 2020 as a part of the NTNU course *TTK4145 - Real Time Programming*


## STATUS
1. Finished!

## How to run
1. run main.go with the flags elevID and port. Given the example that you would like to run an elevator with elevID=1 and port=14001:
```bash
$ go run main.go -elevID=1 -port=14001
```
elevID must be 1, 2 and 3 for the individual elevators. Several elevators can not be initialized with the same elevID. Port must be the same as your simulator.

2. You can also use the script *heis.sh*. Use the following syntax to start an elevator at port 14001 with elevID=1:
```bash
$ ./Heis.sh 14001
```

## Information about the elevator
1. We are using a master-slave system. At one point in time, there will always be one and only one Master. This Master is giving the other elevators on the network *currentOrders*. One elevator is always only executing one *currentOrder* at a time, but is checking for other relevant orders at every floor it passes.
2. Even though we use a master-slave system, the Master can change at any time. Therefore, all elevators are storing all information about all other elevators. Including Cab Orders. Cab Orders are only sent to its owner in case of elevator shutdown. 
3. To keep track of current state of all elevators, all elevators has a local elevatorMap. This maps an elevID to a pointer to an elevator struct. This elevator struct has all necessary info about the spesific elevator. Each elevator also has a cabOrdersBackup that maps an elevator to a backup of that elevators Cab Orders.
4. The elevator is initialized by going down until it hits a floor.
5. All hardcoded values are in the config file. Here you can also change the behaviour for "plugged out network cable" scenario.
6. Delivered code from these projects are *elevio.go* and all of the files in *network*-folder except *network.go*. We have used the libraries *strconv, fmt, flag* and *time*. All other code is written by us.

## Useful
1. Add scripts to PATH
```bash
$ cd
$ vim .bashrc
```
Add the following to the bottom of the file:
```bash
$ export PATH=$PATH:path/to/file    for example ~/sanntid/Simulator eller ~/sanntid/scripts
```
You can now run the program/script from wherever in the file system

2. Make scripts executable:
```bash
$ chmod +x <filename>
```
3. How does *Heis.sh* work?
To run an elevator on port 14001, run the following. Elev_ID will automatically be the last number, in this case 1.
```bash
$ Heis.sh 14001
```
4. FileOpener.sh does not work
You need to update the first line of code. It is originally
```bash
$ cd ~/sanntid/project/
```

## FAULT HANDLING
1. Motor Failure
Press 8 to stop motor. Press 7(down) or 9(up) to start motors.
2. Packet Loss
```bash
$ PacketLoss.sh
```
To flush filter chain:
```bash
$ sudo iptables -F
```
3. Total network shutdown
```bash
$ NetworkLoss.sh
```
