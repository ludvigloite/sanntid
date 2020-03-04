# TODO
1. Nye ordre kommer ikke inn når døra er åpen
2. Bruke goroutines og channels istedet for basically C-kode
3. Gjøre koden mye penere. Få all order_managment til å bare skje inne i order_managment modulen. Altså ikke bruk order.func(order.GetX(),order.GetY())
4. Fikse nettverk



#For å kjøre ElevatorDriver
1. Gå i terminal
2. skriv "pwd",  "/home/student/" bør da komme opp. Hvis ikke, skriv "cd".
3. Skriv "cd .cargo/bin"
4. Skriv "./ElevatorServer"
5. Man kan nå kjøre main.go så vil heisen kjøre.



