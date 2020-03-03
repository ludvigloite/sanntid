
	Feil
		Motor mister strøm(heis klarer ikke bevege seg, kommunikasjon fremdeles intakt)
			Watchdog-timer merker at det har gått mer enn 3 sek mellom etasjer -> Master vet at heisen ikke har strøm
				Hvis heisen er i en etasje -> Åpne døra
				Hvis heisen er mellom to etasjer -> Vent
				Alle ordre som aktuell heis har vil bli gitt til andre.
				Alle cabin orders som kommer inn og har kommet inn skal lagres. Når strømmen kommer tilbake igjen blir disse utført.

		Når heis får strøm:
			Aktuell heis vil igjen begynne å motta ordre fra Master.

	%TODO TA HENSYN TIL OM DET ER MASTER SOM KRASJER/MISTER NETT
		Program krasjer/PC mister strøm
			Master mottar ikke lenger UDP-pakker og vil skjønne at heis er død
				Alle ordre som er gitt til aktuell heis vil bli gitt til andre
				Heisen vil ikke motta nye ordre.

			Når heis kommer tilbake:
				Master har lagret cab calls. Aktuell heis utfører disse ordrene. Begynner også å motta ordre.


		Nettverksfeil:
			Master vet ikke om en feil er nettverksfeil eller programkrasj.
				Gjør det samme som ved programkrasj.
			Heis som mister nettverk mottar ikke lenger UDP-pakker fra Master
				???????Heis som mistet nettverk vil ikke utføre noen hall calls. ???????
				Alle cabin calls som finnes fra før og som kommer inn skal utføres.

			Når heis får nettverk igjen:
				Må gi beskjed til Master om den har hatt programfeil eller nettverksfeil
					Ved Nettverkfeil:
						Aktuell heis sender sine ordre til Master
					Ved Programkrasj:
						Master sender cab calls til aktuell heis.




	Hall order knapp på heis_x blir trykket:
		1. Ordreliste på heis_x blir oppdatert
		2. Endringer blir syncet til master. Bruker ID til å se at denne har kommet senere, og vil derfor godta, og oppdatere egen ordreliste.
		3. Oppdatert ordreliste blir syncet ut til alle heiser.
		4. Master regne ut hvilken heis som skal ta oppdraget ut ifra en cost function.
		5. Master sender ut en oppgave på et en task skal utføres til en heis. Alle heisene har max en current_order. Master kan sende ut en ny current_order til en av heisene, og da vil denne være gjeldende. Current_order fra ella heiser må være del av syncen.
		6. Aktuell heis utfører ordren.
		7. Aktuell heis gir beskjed til Master når oppgaven er utført. 


	Cab order knapp inni heis_x blir trykket:
		Nøyaktig samme skjer som hvis hall order blir trykket. Cab orders kan dog bare utføres av heis_x





	FAULT HANDLING
		Packet loss:
			Bruker UDP som sender masse pakker med likt innhold. Bruker packet_ID og lift_ID for å skille mellom pakker etter når de ble sendt ut, ved motstridende info

		Syncroniseringsfeil / motstridende info
			Bruker packet_ID og lift_ID for å skille mellom pakker etter når de ble sendt ut, ved motstridende info



				

