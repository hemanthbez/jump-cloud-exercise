# jump-cloud-exercise

# Steps to run the program locally

1. Setup the GOROOT to the **go** install directory using the below command <br/>
`export GOROOT=/usr/local/go`
1. Open a Terminal window, and Navigate to the folder - '/go/src/github.com/hemanthbez/jump-cloud-exercise/'

2. Launch the program using the below command passing the Port number argument. Default Port number used is "8080" <br/>
`go run *.go [Port_number] clean` (e.g. *go run *.go -port=8085*) clean

3. To Kill the program, press CTRL+C on the Terminal window.


# Endpoints Info

1. "/hash" (with POST) - Used for sending a POST request with password data. <br/>
e.g. *curl -d "password=abc123" http://localhost:8090/hash*

2. "/hash" (with GET) - Used for getting the SHA512 password based on the ID<br/>
e.g. *curl http://localhost:8090/hash/3*

3. "/stats" (with GET) - Returns the stats<br/>
e.g. *curl http://localhost:8090/stats*

4. "/shutdown" (with GET) - Shuts down the server once all the Pending requests are processed <br/>
e.g. *curl http://localhost:8090/shutdown*


# Source Code Info

1. **password-hash.go** : Main file that includes all the application endpoints.

2. **rest-handler.go** : Utility *go* file that has the basic HTTP handlers defined. 
