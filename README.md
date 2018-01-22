# log_monitor

This is a simple go program to monitor a W3C common formatted log being written to on the local system. It will run until you use ^C to exit.

## Assumptions

1. There is a webserver running on the local system that is writing W3C common formatted logs that can be read.
2. There is a place where we can write out a temporary sqlite database. (Defaults to /tmp/log.db)
3. The end user is familiar with setting up a working go environment to get a working binary

## Running

1. `go get` to fetch dependencies
2. `go run main.go` to start the program
3. (Optionally) Run `./scripts/gen_traffic.sh` to send traffic to the webserver. (Assumes the webserver is listing on localhost:80)

## Tests

1. `go get -t` to test dependencies and test dependencies
2. `go test -v` to run the tests

## Building
1. `go get` to fetch dependencies
2. `go build`

## Improvements and TODO

1. Add more useful/interesting output about the overall traffic
2. Move to logrus to get better output handling.
3. Factor out parsing and general db handling into functions to make main() more readable.
4. Not super happy with the shape of the manageAvgTraffic and getAvgTraffic functions
5. The parsing does not account for HTTP Status and calculates all requests the same. In the future HTTP status should be broken out, etc.
6. The parsing is janky at best. Porting https://github.com/xojoc/logparse would probably be best.

