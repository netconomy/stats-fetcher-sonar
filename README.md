#stats-fetcher-sonar

Simple Server written in Go for fetching metrics from a configurable SonarQube instance and pushing these metrics to a configurable endpoint via UDP. 

#BUILD

make deps

make

#TEST

make test

#CONFIG

Use ./config.json or provide another configuration using the --config={file} flag.

Configuration:

* metrics to fetch
* the SonarQube instance to fetch from
* the UDP Endpoint to push the data to

#RUN

./statsfetcher serve {port}

access via: http://localhost:{port}/fetch/{projectname}

