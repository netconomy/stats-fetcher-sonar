#Sonar Stats Fetcher

Simple Server written in Go for fetching metrics from a configurable SonarQube instance and pushing these metrics to a configurable endpoint via UDP. 

#BUILD IT

make deps
make

#TEST IT

make test

#CONFIG IT

Use ./config.json or provide another configuration using the --config={file} flag.
Configuration:
* metrics to fetch
* the SonarQube instance to fetch from
* the UDP Endpoint to push the data to

#RUN IT

./statsfetcher serve {port}

access via: http://localhost:{port}/fetch/{projectname}

#TODOS

* multithreading / channels, that UDPConn isn't shared

