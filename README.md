# Atlassian test task

# Installation

    go get -u github.com/yurii-stakhiv/atlassian

# Command line

    cd $GOPATH/src/github.com/yurii-stakhiv/atlassian/cmd
    go build

Usage

    ./cmd 'Good morning @bob! (megusta) (coffee) http://www.nbcolympics.com'


# API

    cd $GOPATH/src/github.com/yurii-stakhiv/atlassian/server
    go build

Usage

    ./server :8080

    curl -X POST -d '{"input": "Good morning @bob! (megusta) (coffee) http://www.nbcolympics.com"}' http://localhost:8080/tokenize
