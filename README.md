# scrapligo-play

What better way to learn Go than with practical examples of network automation with scrapligo.

Based around example code found at https://github.com/PacktPublishing/Network-Automation-with-Go

- waitgroup1 - Waitgroup with nothing to limit number of goroutines.
- waitgroup2 - Waitgroup with chunks of data to restrict the number of goroutines.
- waitgroup3 - Waitgroup with channel to restrict number of goroutines.
- waitgroup3a - Added returning results
- waitgroup3b - Restructure into functions
- waitgroup3c - Attempt at a simple playbook/runbook/taskbook in code example, results placed in host struct
- waitgroup3d - Attempt at a simple playbook/runbook/taskbook in code example, results return via channel
- waitgroup3e - Like waitgroup3d but trying interfaces, filters, and a simple package structure (work in progress)
- waitgroup4 - Waitgroup with semaphore to restrict number of goroutines.
- workerpool1 - In/Out buffered channels with a boolean 'done' channel for completion.
- workerpool2 - In/Out buffered channels with a results returned channel.
- workerpool2a - add connection function and host pointers.
- workerpool2b - add yaml inventory, filter, and write config to file.
- workerpool3 - waitgroup/(in)channel workerpool combo, storing results in the host pointer