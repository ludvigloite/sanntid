PORT=$1
ELEV_ID=$(echo "${PORT: -1}")

if [ $ELEV_ID = 0 ]
then
	echo "Du har ikke gitt input!!!"
	echo "Avslutter..."
else
	echo "Starter heis med ID $ELEV_ID på port $PORT"
	cd ~/sanntid/project
	go run main.go -elevID=$ELEV_ID -port=$PORT
fi
