# Exercise 9 - Scheduling

## Properties

### Task 1:
 1. Why do we assign priorities to tasks?
 	We want that the closest elevator should take the order. This means we need some kind of cost function. This is for effectiveness.

 2. What features must a scheduler have for it to be usable for real-time systems?
 	It must be able to take new inputs and place it in the middle of the queue.
 	It must be able to handle errors, like one "user" falling out etc.

 

## Inversion and inheritance


| Task | Priority   | Execution sequence | Release time |
|------|------------|--------------------|--------------|
| a    | 1          | `E Q V E`          | 4            |
| b    | 2          | `E V V E E E`      | 2            |
| c    | 3 (lowest) | `E Q Q Q E`        | 0            |

 - `E` : Executing
 - `Q` : Executing with resource Q locked
 - `V` : Executing with resource V locked


### Task 2: Draw Gantt charts to show how the former task set:
 1. Without priority inheritance
 	Done. PDF

 2. With priority inheritance
 	Done. PDF

### Task 3: Explain:
 1. What is priority inversion? What is unbounded priority inversion?
 	That a high priority task is indirectly interrupted by a lower priority task, and therefore "has higher priority". Often because of shared resources. Like in our example, when c har "higher rank" than a, when executing Q.
 2. Does priority inheritance avoid deadlocks?
 	no.




## Utilization and response time

### Task set 2:

| Task | Period (T) | Exec. Time (C) |
|------|------------|----------------|
| a    | 50         | 15             |
| b    | 30         | 10             |
| c    | 20         | 5              |

### Task 4:
 1. There are a number of assumptions/conditions that must be true for the utilization and response time tests to be usable (The "simple task model"). What are these assumptions? Comment on how realistic they are.
 	Fixed set of periodic, independent tasks with known periods, constant worst-case execution times, deadlines equal to their periods. Runs on single processor. Overhead runs in zero time. 
 	We think these are fairly realistic.

 2. Perform the utilization test for the task set. Is the task set schedulable?
 	No it is not schedulable. 0,88 > 0,78.

 3. Perform response-time analysis for the task set. Is the task set schedulable? If you got different results than in 2), explain why.
 	We chose c to have highest priority
 	Rc = 5

 	Rb = 15 

 	Ra = 30

 4. (Optional) Draw a Gantt chart to show how the task set executes using rate monotonic priority assignment, and verify that your conclusions are correct.

## Formulas

Utilization:  
![U = \sum_{i=1}^{n} \frac{C_i}{T_i} \leq n(2^{\frac{1}{n}}-1)](eqn-utilization.png)

Response-time:  
![w_{i}^{n+1} = C_i + \sum_{j \in hp(i)} \bigg \lceil {\frac{w_i^n}{T_j}} \bigg \rceil C_j](eqn-responsetime.png)


