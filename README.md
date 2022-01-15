# scrapligo-play

What better way to learn Go than with practical examples of network automation with scrapligo.

Based around example code found at https://github.com/PacktPublishing/Network-Automation-with-Go

- waitgroup1 - Waitgroup with nothing to limit number of goroutines.
- waitgroup2 - Waitgroup with chunks of data to restrict the number of goroutines.
- waitgroup3 - Waitgroup with channel to restrict number of goroutines.
- waitgroup4 - Waitgroup with semaphore to restrict number of goroutines.
- workerpool1 - In/Out buffered channels with a boolean 'done' channel for completion.
- workerpool2 - In/Out buffered channels with a results returned channel.
- workerpool3 - add connection function and host pointers.
- workerpool4 - add yaml inventory, filter, and write config to file. (work in progress)