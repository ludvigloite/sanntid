# STATUS
1. Nærmer oss veldig ferdig nå!!
2. Ved "ethernet-utdragning" vil heisen fullføre de nåværende ordre, men vil ikke ta flere.
3. Trenger å teste at alt fungerer som planlagt og etter oppgavespesifikasjonene.
4. Må kjøre vår egen Final Acceptance Test (FAT)
5. Koden må ryddes og kommenteres.

## NYTTIG
1. Kill eldste prosess fra sanntids-PC
```bash
$ pkill -o -u ludvig sshd
```
2. Fra Mac:
```bash
$ osascript sanntid_terminal_opener.scpt
```
3. Fra remote:
```bash
$ ./open_files.sh
```

## NYTTIG
1. Legg til scripts i PATH
```bash
$ cd
$ vim .bashrc
```
Legg til følgende på bunn av fila:
```bash
$ export PATH=$PATH:path/to/file    for eksempel ~/sanntid/Simulator eller ~/sanntid/scripts
```
Du kan nå kjøre programmet/scriptet uansett hvor du er i filsystemet

2. Gjør scripts kjørbare:
```bash
$ chmod +x <filename>
```
3. Hvordan funker 'Heis.sh' ?
For å kjøre en heis på port 14001, kjør følgende. Elev_ID vil automatisk bli det siste sifferet, i dette tilfellet 1.
```bash
$ Heis.sh 14001
```
4. FileOpener.sh funker ikke
Du må oppdatere første linja med kode, som hos meg er
```bash
$ cd ~/sanntid/project/
```
Her kan du også velge hvilken rekkefølge filene åpnes i. Ikke alle filene i prosjektet blir åpnet.


## FAULT HANDLING
1. Motor Failure
Trykk 8 for å stoppe motor. Trykk 7(ned) eller 9(opp) for å starte igjen
2. Pakketap
```bash
$ PacketLoss.sh
```
For å flushe filter chain:
```bash
$ sudo iptables -F
```
3. Nettverkstrøbbel
```bash
$ NetworkLoss.sh
```

## For å kjøre ElevatorDriver(fysisk på sanntidssal)
1. Gå i terminal
2. skriv "pwd",  "/home/student/" bør da komme opp. Hvis ikke, skriv "cd".
3. Skriv "cd .cargo/bin"
4. Skriv "./ElevatorServer"
5. Man kan nå kjøre main.go så vil heisen kjøre.


## Hvordan update branch til Master
https://gist.github.com/santisbon/a1a60db1fb8eecd1beeacd986ae5d3ca

First we'll update your local master branch. Go to your local project and check out the branch you want to merge into (your local master branch)
```bash
$ git checkout master
```

Fetch the remote, bringing the branches and their commits from the remote repository.
You can use the -p, --prune option to delete any remote-tracking references that no longer exist in the remote. Commits to master will be stored in a local branch, remotes/origin/master
```bash
$ git fetch -p origin
```

Merge the changes from origin/master into your local master branch. This brings your master branch in sync with the remote repository, without losing your local changes. If your local branch didn't have any unique commits, Git will instead perform a "fast-forward".
```bash
$ git merge origin/master
```

Check out the branch you want to merge into
```bash
$ git checkout <feature-branch>
```

Merge your (now updated) master branch into your feature branch to update it with the latest changes from your team.
```bash
$ git merge master
```

Depending on your git configuration this may open vim. Enter a commit message, save, and quit vim: 
1. Press `a` to enter insert mode and append text following the current cursor position.
2. Press the **esc** key to enter command mode.
3. Type `:wq` to write the file to disk and quit.

This only updates your local feature branch. To update it on GitHub, push your changes.
```bash
$ git push origin <feature-branch>
```
