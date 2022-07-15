# Scrapligo and "go test" for testing operational state of the network

I could have had the tests and collectors with the same package name.  I decided to try in separate packages. This meant I needed to end the files containing the functions that return the anonymous test function with _test and have the directory equal to the package name.

Speaking of the functions that return the anonymous test functions, although it didn't need to, I've named them like you would in a proper test function, but with a lower case 't' so that it can never be automatically run by 'go test'.

This uses sub-tests so I think you need at least Go1.17.  The first level is the test function itself, the second is a loop through the hosts, the third is a loop through the each task, and you can have further levels within loops of the task, like in testBgp.

It needs a better way of packaging, but this is primarily a proof of concept only.

And why did I suddenly needed to use export CGO_ENABLED=0?

```
PS D:\Users\user\Documents\workspace\scrapligo-play\tester1> go test . -count=1 -v -name 192.168.204.104 -name no.suchdomain -task version -task bgp
2022/03/28 20:10:43 Names: 192.168.204.104, no.suchdomain
2022/03/28 20:10:43 Tasks: version, bgp
2022/03/28 20:10:43 Gathering ...
2022/03/28 20:10:43 failed to open driver: dial tcp: lookup no.suchdomain: no such host
2022/03/28 20:10:44 Runner took 488.9382ms
=== RUN   TestHost
=== RUN   TestHost/no.suchdomain
=== RUN   TestHost/no.suchdomain/connection
    a_connection_test.go:21: failed to open driver: dial tcp: lookup no.suchdomain: no such host
=== RUN   TestHost/no.suchdomain/version
    a_version_test.go:17: Skipped due to no results
=== RUN   TestHost/no.suchdomain/bgp
    a_bgp_test.go:16: Skipped due to no results
=== RUN   TestHost/192.168.204.104
=== RUN   TestHost/192.168.204.104/connection
=== RUN   TestHost/192.168.204.104/version
=== RUN   TestHost/192.168.204.104/bgp
=== RUN   TestHost/192.168.204.104/bgp/10.0.24.2
=== RUN   TestHost/192.168.204.104/bgp/10.0.46.6
    a_bgp_test.go:44: BGP in Idle/Active state
--- FAIL: TestHost (0.00s)
    --- FAIL: TestHost/no.suchdomain (0.00s)
        --- FAIL: TestHost/no.suchdomain/connection (0.00s)
        --- SKIP: TestHost/no.suchdomain/version (0.00s)
        --- SKIP: TestHost/no.suchdomain/bgp (0.00s)
    --- FAIL: TestHost/192.168.204.104 (0.00s)
        --- PASS: TestHost/192.168.204.104/connection (0.00s)
        --- PASS: TestHost/192.168.204.104/version (0.00s)
        --- FAIL: TestHost/192.168.204.104/bgp (0.00s)
            --- PASS: TestHost/192.168.204.104/bgp/10.0.24.2 (0.00s)
            --- FAIL: TestHost/192.168.204.104/bgp/10.0.46.6 (0.00s)
FAIL
FAIL    tester1 0.703s
FAIL

PS D:\Users\user\Documents\workspace\scrapligo-play\tester1> go test . -count=1 -v
2022/03/28 20:11:26 Names: 
2022/03/28 20:11:26 Tasks: 
2022/03/28 20:11:26 Gathering ...
2022/03/28 20:11:26 failed to open driver: dial tcp: lookup no.suchdomain: no such host
2022/03/28 20:11:27 Runner took 534.7844ms
=== RUN   TestHost
=== RUN   TestHost/192.168.204.102
=== RUN   TestHost/192.168.204.102/connection
=== RUN   TestHost/192.168.204.102/version
=== RUN   TestHost/192.168.204.102/bgp
=== RUN   TestHost/192.168.204.102/bgp/10.0.12.1
    a_bgp_test.go:42: BGP Estatblished but zero learned routes
=== RUN   TestHost/192.168.204.102/bgp/10.0.24.4
    a_bgp_test.go:42: BGP Estatblished but zero learned routes
=== RUN   TestHost/192.168.204.103
=== RUN   TestHost/192.168.204.103/connection
=== RUN   TestHost/192.168.204.103/version
=== RUN   TestHost/192.168.204.103/bgp
=== RUN   TestHost/192.168.204.103/bgp/10.0.13.1
=== RUN   TestHost/192.168.204.103/bgp/10.0.35.5
    a_bgp_test.go:44: BGP in Idle/Active state
=== RUN   TestHost/192.168.204.104
=== RUN   TestHost/192.168.204.104/connection
=== RUN   TestHost/192.168.204.104/version
=== RUN   TestHost/192.168.204.104/bgp
=== RUN   TestHost/192.168.204.104/bgp/10.0.24.2
=== RUN   TestHost/192.168.204.104/bgp/10.0.46.6
    a_bgp_test.go:44: BGP in Idle/Active state
=== RUN   TestHost/no.suchdomain
=== RUN   TestHost/no.suchdomain/connection
    a_connection_test.go:21: failed to open driver: dial tcp: lookup no.suchdomain: no such host
=== RUN   TestHost/no.suchdomain/version
    a_version_test.go:17: Skipped due to no results
=== RUN   TestHost/no.suchdomain/bgp
    a_bgp_test.go:16: Skipped due to no results
=== RUN   TestHost/192.168.204.101
=== RUN   TestHost/192.168.204.101/connection
=== RUN   TestHost/192.168.204.101/version
=== RUN   TestHost/192.168.204.101/bgp
=== RUN   TestHost/192.168.204.101/bgp/10.0.12.2
=== RUN   TestHost/192.168.204.101/bgp/10.0.13.3
    a_bgp_test.go:42: BGP Estatblished but zero learned routes
--- FAIL: TestHost (0.00s)
    --- FAIL: TestHost/192.168.204.102 (0.00s)
        --- PASS: TestHost/192.168.204.102/connection (0.00s)
        --- PASS: TestHost/192.168.204.102/version (0.00s)
        --- FAIL: TestHost/192.168.204.102/bgp (0.00s)
            --- FAIL: TestHost/192.168.204.102/bgp/10.0.12.1 (0.00s)
            --- FAIL: TestHost/192.168.204.102/bgp/10.0.24.4 (0.00s)
    --- FAIL: TestHost/192.168.204.103 (0.00s)
        --- PASS: TestHost/192.168.204.103/connection (0.00s)
        --- PASS: TestHost/192.168.204.103/version (0.00s)
        --- FAIL: TestHost/192.168.204.103/bgp (0.00s)
            --- PASS: TestHost/192.168.204.103/bgp/10.0.13.1 (0.00s)
            --- FAIL: TestHost/192.168.204.103/bgp/10.0.35.5 (0.00s)
    --- FAIL: TestHost/192.168.204.104 (0.00s)
        --- PASS: TestHost/192.168.204.104/connection (0.00s)
        --- PASS: TestHost/192.168.204.104/version (0.00s)
        --- FAIL: TestHost/192.168.204.104/bgp (0.00s)
            --- PASS: TestHost/192.168.204.104/bgp/10.0.24.2 (0.00s)
            --- FAIL: TestHost/192.168.204.104/bgp/10.0.46.6 (0.00s)
    --- FAIL: TestHost/no.suchdomain (0.00s)
        --- FAIL: TestHost/no.suchdomain/connection (0.00s)
        --- SKIP: TestHost/no.suchdomain/version (0.00s)
        --- SKIP: TestHost/no.suchdomain/bgp (0.00s)
    --- FAIL: TestHost/192.168.204.101 (0.00s)
        --- PASS: TestHost/192.168.204.101/connection (0.00s)
        --- PASS: TestHost/192.168.204.101/version (0.00s)
        --- FAIL: TestHost/192.168.204.101/bgp (0.00s)
            --- PASS: TestHost/192.168.204.101/bgp/10.0.12.2 (0.00s)
            --- FAIL: TestHost/192.168.204.101/bgp/10.0.13.3 (0.00s)
FAIL
FAIL    tester1 0.714s
FAIL

```