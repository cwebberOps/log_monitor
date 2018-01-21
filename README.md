# log_monitor

This is a simple go program to monitor a W3C common formatted log being written to on the local system. It will run until you use ^C to exit.

## Assumptions

In this version of the application we make a ton of hardcoded (read "bad") assumptions to get something basic working.

1. There is a webserver running on the local system that is writing logs to /private/var/log/apache2/access_log
2. We are looking at the last 10s of traffic for the sections of the site
3. We are looking at the last 2m of traffic for overall traffic alerts
4. We are alerting when the traffic is in excess of an average of 5 requests in the last 2m
5. We can write to /tmp/log.db
6. The end user is familiar with setting up a working go environment to get a working binary


## Running

1. `go get` to fetch dependencies
2. `go run main.go` to start the program
3. (Optionally) Run `./gen_traffic.sh` to send traffic to the webserver. (Assumes the webserver is listing on localhost:80)

## Tests

1. `go get -t` to test dependencies and test dependencies
2. `go test -v` to run the tests

## Improvements and TODO

1. Move all the constants into command line args so they can be specified.
2. Add more useful/interesting output about the overall traffic
3. Move to logrus to get better output handling.
4. Factor out parsing and general db handling into functions to make main() more readable.
5. Add general usage info when `--help` is called
6. Not super happy with the shape of the manageAvgTraffic and getAvgTraffic functions
7. The parsing does not account for HTTP Status and calculates all requests the same. In the future HTTP status should be broken out, etc.
