# no-name-server [![CircleCI](https://circleci.com/gh/Yinkozi/no-name.svg?style=svg&circle-token=a18ffbc369b8ddcf8de823bc2a1eeb628509fcb7)](https://circleci.com/gh/Yinkozi/no-name)


## Environment Installation
1. Visit https://golang.org/doc/install
2. Install glide (MAC Osx : brew install glide)
3. Configure your ssh key so your able to pull the private repository of the dependencies manager (ssh-add <your_ssh_github_key>)

## Compile
```
git clone https://github.com/Yinkozi/no-name-server
cd no-name-server
glide install
go build server.go
```

## Run
```
./server
```
Server should be accessible at http://localhost:8080//apidocs.json  

## Test
1 - Launch tests  
```
go test (go list ./... | grep -v /vendor/)
```
