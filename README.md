# OSwithGO

# Process Management Simulator

## Brief Project Description:
The process scheduling simulator is a tool that enables users to create, manage, and simulate process scheduling in an operating system. It allows users to manually create processes or generate random ones, select from different scheduling algorithms (e.g., Round Robin), define the quantum time for the algorithm, and execute the scheduling simulation.

### Modeling Processes as Goroutines:
- Processes are represented as goroutines in Go.
- Each goroutine represents a process that can be either CPU-bound or I/O-bound.

### CPU and I/O Access Control:
##### Starvation:
- Starvation is a common issue in process scheduling systems.
- It occurs when a process is perpetually prevented from executing due to the prioritization of other processes.
- This can lead to processes waiting indefinitely, causing delays and impacting overall system performance.

##### Avoiding Starvation:
- Control is implemented using a custom Round Robin algorithm.
- The Round Robin algorithm is adapted to handle both I/O-bound and CPU-bound processes simultaneously.
- It initializes the first I/O-bound and CPU-bound processes in the ready queue.
- Each process is allocated execution time based on the defined quantum.
- If a process does not complete execution within its quantum, it is returned to the end of the ready queue, and the next process in the list is scheduled for execution.
- This mechanism prevents starvation by ensuring an equal share of CPU and I/O time for all processes, guaranteeing that all processes have an opportunity to execute. It avoids resource monopolization and ensures fair and balanced CPU and I/O time distribution among processes.

### Shared Resources:
- Shared resources include `cpuResource` and `ioResource` channels, used to control access to the CPU and I/O resources, respectively.
- Processes share these resources to coordinate CPU and I/O operations.
- These resources are implemented in the `cpu.go` and `io.go` files:
  - They create a text file for CPU or I/O.
  - During simulation, they log the process name, ID, and usage time in the corresponding file.
  - Once finished, the resource is released for other processes to use.

### Concurrency Control:
- Concurrency is managed using channels for synchronization between the goroutines representing the processes.
- Channels coordinate access to shared resources, ensuring that only one process can access the CPU or I/O device at a time.

### Parameterized Implementation:
- The scheduling algorithm and quantum time are configurable via a user menu.
- Shared resources and synchronization are globally managed in the code, ensuring consistency and proper concurrency control.

## Detailed Description:

1. **Process Creation:**
   - Users can manually create processes by inputting details such as name, priority, whether the process is I/O-bound, and total CPU time.
   - Alternatively, users can generate a specified number of random processes, with randomly assigned attributes like name, priority, I/O-bound status, and total CPU time.

2. **Scheduling Algorithm Selection:**
   - Users can choose the available scheduling algorithm: Round Robin.
   - Scheduling actions are defined based on the selected algorithm's rules.

3. **Quantum Time Selection:**
   - Users can define the quantum time, which is the maximum time interval a process can occupy the CPU or I/O before being preempted.

4. **Simulation Execution:**
   - The simulation begins when users choose to execute the scheduling.
   - During the simulation, the ready queue is displayed, and the selected scheduling algorithm is applied to processes in the queue.
   - Processes execute according to the algorithm's rules until all processes are completed.
   - During simulation:
     - I/O-bound processes initiate their I/O operations while waiting for the CPU.
     - CPU-bound processes utilize the CPU.
     - If an I/O-bound process needs the CPU, it waits until the CPU is available.
     - If a CPU-bound process needs I/O, it waits until the I/O device is free.
   - This management is handled in the `executeProcess` function in the `scheduler.go` file:
     - The function uses `cpuResource = make(chan bool, 1)` to control CPU resources and `ioResource = make(chan bool, 1)` to manage I/O resources.
     - Channels in Go facilitate communication and synchronization between goroutines (lightweight threads managed by the Go runtime).
     - The `executeProcess` function employs goroutines to synchronously manage CPU and I/O resource usage, ensuring proper behavior during simulation.
   - Additionally, the `wg` variable of type `sync.WaitGroup` is used to coordinate and wait for the completion of goroutines, ensuring the simulation concludes correctly.

5. **Average Waiting Time Calculation:**
   - After the simulation concludes, the average waiting time of processes is calculated and displayed.
   - The average waiting time represents the mean time processes wait in the ready queue before execution.

6. **Interactive Menu Options:**
   - An interactive menu enables users to navigate through the simulator's features, such as creating processes, choosing algorithms, defining quantum time, starting simulations, generating random processes, and exiting the program.