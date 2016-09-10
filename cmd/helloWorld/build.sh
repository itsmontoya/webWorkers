#!/bin/bash
go build # Build main.go 
sudo setcap 'cap_net_bind_service=+ep' helloWorld # Allow helloWorld binary to listen to port(s) 80 and/or 443
