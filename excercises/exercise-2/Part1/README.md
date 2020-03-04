# Mutex and Channel basics

### What is an atomic operation?

> An atomic operation are programming operations that runs independently of other processes. In other words an opoeration is atomic if it completes in a single step reative to other threads.

### What is a semaphore?
> A semaphore is a variable. It is used to control access to a resource by multiple processes in a concurrent system such as multitasking operating systems. Semaphores are of two types: Binary semaphores and counting semaphores.

### What is a mutex?
> A mutex provide mutal exclusion which is a type of concurrency control that has a purpose of preventing race conditions. One thread of execution should not enter its critical section at the same time as another concurrent thread is entering its own section.

### What is the difference between a mutex and a binary semaphore?
> A semaphore is a generalized mutex. This means that the binary semaphore can be signaled by any thread while the mutex only will be released by the thread that have acquired it. A mutex can therefore be seen as a locking machanism while the semaphore is a signaling mechanism.

### What is a critical section?
> A critical section has to be executed as an atomic section as it contains shared variables, and hence multiple processes can't execute this section simultaneously. If the critical section is not executed atomic it can lead to wrong answers.

### What is the difference between race conditions and data races?
 > A race condition happens when two ore more threads tries to access shared data and change it at the same time. A data races is when two instructions from different threads tries to access the same memory location and at least one of them is a write instruction.  

### List some advantages of using message passing over lock-based synchronization primitives.
> Message passing: parallel programming
Synchronization primitives needs to wait for each other while transferring the message. Message passing just sends a message and relies on the reciever process to select and run the code.

### List some advantages of using lock-based synchronization primitives over message passing.
> You will always be certain that your memory is protected. No 2 functions can manipulate the same data.
