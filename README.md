# sanntid

STATUS

  Lurer paa om systemet med hvordan rank endres naar noder faller ikkke fungerer. I tillegg detter noder inn og ut hele tiden. Kan hende det er at jeg ikke tillatter aa sende UDP paa PCen min, maa evt sjekke dette hvis dere andre klarer aa kjore koden uproblematisk paa deres maskiner.

  Slik programmet er naa maa man alltid ha en med rank 1 for at det skal fungere. Altsaa maa heis med ID = 1 alltid begynne for at heisen skal fungere.  

NYTTIG

    Kill eldste prosess fra sanntids-PC

$ pkill -o -u ludvig sshd

    Fra Mac:

$ osascript sanntid_terminal_opener.scpt

    Fra remote:

$ ./open_files.sh

For å kjøre ElevatorDriver

    Gå i terminal
    skriv "pwd", "/home/student/" bør da komme opp. Hvis ikke, skriv "cd".
    Skriv "cd .cargo/bin"
    Skriv "./ElevatorServer"
    Man kan nå kjøre main.go så vil heisen kjøre.

Kjør programmet med ELEV_ID

    Naviger til mappen
    NUMMER er enten 1, 2 eller 3

$ go build main.go
$ ./main -elevID='NUMMER'

Linker for å kjøre Ludvigs mac

    KJØR Simulator

$ ./Desktop/Local\ Storage/heisSimulator/Simulator-v2/SimElevatorServer --port 10001

    Åpne prosjektet

$ cd Desktop/Local\ Storage/Sanntid_prosjekt/sanntid

Åpne filer:
$ subl project

Build(inne i /project):
$ go build main.go

Kjør:
$ ./main -elevID=1 -port=10001
$ ./Desktop/Local\ Storage/Sanntid_prosjekt/sanntid/project/main -elevID=1 -port=10001

Hvordan update branch til Master

https://gist.github.com/santisbon/a1a60db1fb8eecd1beeacd986ae5d3ca

First we'll update your local master branch. Go to your local project and check out the branch you want to merge into (your local master branch)

$ git checkout master

Fetch the remote, bringing the branches and their commits from the remote repository. You can use the -p, --prune option to delete any remote-tracking references that no longer exist in the remote. Commits to master will be stored in a local branch, remotes/origin/master

$ git fetch -p origin

Merge the changes from origin/master into your local master branch. This brings your master branch in sync with the remote repository, without losing your local changes. If your local branch didn't have any unique commits, Git will instead perform a "fast-forward".

$ git merge origin/master

Check out the branch you want to merge into

$ git checkout <feature-branch>

Merge your (now updated) master branch into your feature branch to update it with the latest changes from your team.

$ git merge master

Depending on your git configuration this may open vim. Enter a commit message, save, and quit vim:

    Press a to enter insert mode and append text following the current cursor position.
    Press the esc key to enter command mode.
    Type :wq to write the file to disk and quit.

This only updates your local feature branch. To update it on GitHub, push your changes.

$ git push origin <feature-branch>
