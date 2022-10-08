# scrapligo-play

What better way to learn Go than with practical examples of network automation with scrapligo.

Initially based around example code found at https://github.com/PacktPublishing/Network-Automation-with-Go and has grown from there

- waitgroup1 - Waitgroup with nothing to limit number of goroutines.
- waitgroup2 - Waitgroup with chunks of data to restrict the number of goroutines.
- waitgroup3 - Waitgroup with channel to restrict number of goroutines.
- waitgroup3a - Added returning results
- waitgroup3b - Restructure into functions
- waitgroup3c - Attempt at a simple playbook/runbook/taskbook in code example, results placed in host struct
- waitgroup3d - Attempt at a simple playbook/runbook/taskbook in code example, results returned from function
- waitgroup3e - Like waitgroup3d but 
  - interface decorators for task interfaces
    - think https://refactoring.guru/design-patterns/decorator/go/example
    - solves the connector interface casting problem for different connection methods
  - basic inventory
  - filter
  - simple package structure (after many iterations this now looks like gornir)
  - connection interfaces
    - scrapligo (ssh/netconf)
    - gomiko
    - netrasp
    - go ssh
    - go expect
      - which breaks Windows, so moved to using WSL/Ubuntu
      - also, no interact() which would have been cool
  - restconf example
  - still work in progress
- waitgroup4 - Waitgroup with semaphore to restrict number of goroutines.
- workerpool1 - In/Out buffered channels with a boolean 'done' channel for completion.
- workerpool2 - In/Out buffered channels with a results returned channel.
- workerpool2a - add connection function and host pointers.
- workerpool2b - add yaml inventory, filter, and write config to file.
- workerpool3 - waitgroup/(in)channel workerpool combo, storing results in the host pointer
- tester1 - Scrapligo and "go test" for testing operational state of the network
- template1 - templates, yaml, and tcl to create file on disk via scrapligo
- template2 - tcl ping created from template, plus textfsm and an attempt at tabwriter
- database1 - sqlite - create.go, gather.go, query.go example.
- minigenieparser - practising text parsing like Cisco do with their python Genie project
- miniconfparser - practising (very basic) config parsing similar to the python ciscoconfparse by mpenning

Note: gomiko for waitgroup3e
"git mod tidy" or "go get -u" without the @master didn't get latest, you may need to use
go get -u github.com/Ali-aqrabawi/gomiko/pkg@master

Handy Links
https://yourbasic.org/golang/gotcha-data-race-closure/
