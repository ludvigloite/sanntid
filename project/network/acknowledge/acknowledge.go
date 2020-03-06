package acknowledge



import (
   "fmt"
   "time"
   "os"
   "sync"
   //+++
)

type AckMessage struct {
  node_id string
  timestamp string
  error_id string
  state int
  order_list string
  confirmed_orders string
  current_orders string
}

type SentMessage struct {
  UpdateMessage map[int]
  StatusMessage map[int]
  NumberOfSentMessage map[int]int
  NotRecFromSlave map[int][]string
}

type AckStruct struct {
  AckMsg AckMessage
  AckTimer *time.Timer
}
