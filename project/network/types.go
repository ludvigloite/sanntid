
package network

type Packet struct{
  packetID          int
  timestamp         int
  error_id          int
  state             int
  current_ordr      int
  order_list        [3][4]int
  confirmed_orders  [3][4]int
}
