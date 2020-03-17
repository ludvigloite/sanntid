# TO DO
1. Nye ordre kommer ikke inn når døra er åpen
2. Bruke goroutines og channels istedet for basically C-kode
3. Gjøre koden mye penere. Få all order_managment til å bare skje inne i order_managment modulen. Altså ikke bruk order.func(order.GetX(),order.GetY())
4. Fikse nettverk



## For å kjøre ElevatorDriver
1. Gå i terminal
2. skriv "pwd",  "/home/student/" bør da komme opp. Hvis ikke, skriv "cd".
3. Skriv "cd .cargo/bin"
4. Skriv "./ElevatorServer"
5. Man kan nå kjøre main.go så vil heisen kjøre.


## Kjør programmet med ELEV_ID
1. Naviger til mappen
2. NUMMER er enten 1, 2 eller 3 
```bash
$ go build main.go
$ ./main -elevID='NUMMER'
```

## Linker for å kjøre Ludvigs mac

1. KJØR Simulator
```bash
$ ./Desktop/Local\ Storage/heisSimulator/Simulator-v2/SimElevatorServer --port 10001
```
2. Åpne prosjektet
```bash
$ cd Desktop/Local\ Storage/Sanntid_prosjekt/sanntid

Åpne filer:
$ subl project 

Build(inne i /project):
$ go build main.go

Kjør:
$ ./main -elevID=1 -port=10001
$ ./Desktop/Local\ Storage/Sanntid_prosjekt/sanntid/project/main -elevID=1 -port=10001
```









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
