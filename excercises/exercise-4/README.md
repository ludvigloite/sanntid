# Exercise 4: From Prototype to Production

1. Don't overengineer.
2. Always design properly.
3. Minor detail change will ruin your perfect model.
4. Always prototype first.
5. You will only know what you didn't know on the second iteration.
6. There is only budget/time for the first version.
7. The prototype will become production.


## Network Module

This exercise aims to produce the "v0.1" of the network module used in your project and at the same time prepare you for the network part on the design review.


You should start by taking a look back at [the beginning of Exercise 1](https://github.com/TTK4145/Exercise1/blob/master/Part1/README.md), and reevaluate them in the light of what you have learned about network programming (and - if applicable - concurrency control). At the same time you might want to look into what kind of libraries that already exist for your chosen language.

 - [C network module](https://github.com/TTK4145/Network-c)
 - [D network module](https://github.com/TTK4145/Network-D)
 - [Go network module](https://github.com/TTK4145/Network-go)
 - [Rust network module](https://github.com/edvardsp/network-rust)
 - [Distributed Erlang](http://erlang.org/doc/reference_manual/distributed.html)

By the end of this exercise, you should be able to send some data structure (struct, record, etc) from one machine to another. How you achieve this (in terms of network topology, protocol, serialization) does not matter. The key here is *abstraction*.  

Don't forget that this module should *simplify* the interface between machines: Creating and handling sockets in all routines that need to communicate with the outside world is possible, but is likely to be unwieldy and unmaintainable. We want to encapsulate all the necessary functionality in a single module, so we have a single decoupled component where we can say "This module sends our data over the network". This will almost always be preferable, but above all else: *Think about what best suits your particular design*.








### Design questions

To get you started with designing a network module and/or the application that uses it, try to find answers to the questions below:

 - Guarantees about elevators:
   - What should happen if one of the nodes loses network?
      - Some other nodes should have all orders/all necessary information.
      - All orders delegated to the non-communicating node should be delegated to other nodes.
      - The node that loses network should give up all of its orders.
      - Should try to connect continously to the non-communicating node
      - if Master: another node should be Master.

   - What should happen if one of the computers loses power for a brief moment?
      - initialize the elevator. Some other node always have all necessary information. Some other node(Master?) give the elevator the orders once again.

   - What should happen if some unforeseen event causes the elevator to never reach its destination but communication remains intact?
      - If it uses more than (3) seconds, give the order to some other elevator. Restart the problem elevator.

 - Guarantees about orders:
   - Do all your nodes need to "agree" on an order for it to be accepted? In that case, how is a faulty node handled?
      - We plan on using a master. The master must send an "agree"-message to the receiving elevator.

   - How can you be sure that at least as many nodes as needs to agree on the order in fact agrees on the order?
      - The nodes must send a "agree"-message

   - Do you share the entire state of the current orders , or just the changes as they occur?
      - All nodes have all current orders. If a slave recieves an order, or finishes an order, they will send only the change to the master. Master will send all current orders to all nodes.

     - For either one: What should happen when an elevator re-joins after having been offline?
        - initialize. Should be given orders.

*Pencil and paper is encouraged! Drawing a diagram/graph of the message pathways between nodes (elevators) will aid in visualizing complexity. Drawing the order of messages through time will let you more easily see what happens when communication fails.*

 - Topology:
   - What kind of network topology do you want to  implement? Peer to peer? Master slave? Circle?
      - Master-Slave.

   - In the case of a master-slave configuration: Do you have only one program, or two (a "master" executable and a "slave")?
      - ONE PROGRAM. In this one program the elevators should decide which elevator is master.
              In case of failure in master, another elevator should take over the responsibilty.

     - Is a slave becoming a master a part of the network module?
        - Yes.

 - Technical implementation:
   - If you are using TCP: How do you know who connects to who?
        - Needs to give a acceptance of connection to master.

     - Do you need an initialization phase to set up all the connections?
        - Yes.

   - Will you be using blocking sockets & many threads, or nonblocking sockets and [`select()`](http://en.wikipedia.org/wiki/Select_%28Unix%29)?
      - Non-blocking and select.

   - Do you want to build the necessary reliability into the module, or handle that at a higher level?
      - We want to build the reliability into the network module

   - How will you pack and unpack (serialize) data?
     - Do you use structs, classes, tuples, lists, ...?
          - Our first plan is to implement a struct with this data:
            - array for orders
            - array for each of the elevators(position etc)
            - ++

     - JSON, XML, or just plain strings?
        - JSON

     - Is serialization a part of the network module?
        - We think that we should send orders to another module(to the Master node)

   - Is detection (and handling) of things like lost messages or lost nodes a part of the network module?
        - Yes.

 - Protocols:
   - TCP gives you a data stream that is guaranteed to arrive in the same order as it was sent in (or not at all)
   - UDP might reorder the packets you send into the network //Dette har ikke negativ effekt for vårt tenkte system
   - TCP will resend packets if they're dropped (at least until the socket times out)
   - UDP may drop packets as it pleases //fix kan være å sende mange packets.
   - TCP requires that you to set up a connection, so you will have to know who connects to who
   - UDP doesn't need a connection, and even allows broadcasting
   - (You're also allowed to use any other network library or language feature you may desire)



### Running from another computer

In order to test a network module, you will have to run your code from multiple machines at once. The best way to do this is to log in remotely. Remember to be nice the people sitting at that computer (don't mess with their files, and so on), and try to avoid using the same ports as them.

 - Logging in:
   - `ssh username@#.#.#.#` where #.#.#.# is the remote IP
   student@10.100.23.210
 - Copying files between machines:
   - `scp source destination`, with optional flag `-r` for recursive copy (folders)
   - Examples:
     - Copying files *to* remote: `scp -r fileOrFolderAtThisMachine username@#.#.#.#:fileOrFolderAtOtherMachine`
     - Copying files *from* remote: `scp -r username@#.#.#.#:fileOrFolderAtOtherMachine fileOrFolderAtThisMachine`
s
      - scp ~/Desktop/G14/exercise-3-gruppe-14/test2.go student@10.100.23.210:~/Desktop/hack


*If you are scripting something to automate any part of this process, remember to **not** include the login password in any files you upload to github (or anywhere else for that matter)*

## Extracurricular

[The Night Watch](https://web.archive.org/web/20140214100538/http://research.microsoft.com/en-us/people/mickens/thenightwatch.pdf)
*"Systems people discover bugs by waking up and discovering that their first-born children are missing and "ETIMEDOUT" has been written in blood on the wall."*

[The case of the 500-mile email](http://www.ibiblio.org/harris/500milemail.html)
*"We can't send mail farther than 500 miles from here," he repeated. "A little bit more, actually. Call it 520 miles. But no farther."*

[21 Nested Callbacks](http://blog.michellebu.com/2013/03/21-nested-callbacks/)
*"I gathered from these exchanges that programmers have a perpetual competition to see who can claim the most things as 'simple.'"*
