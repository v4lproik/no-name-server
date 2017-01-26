# no-name-server [![CircleCI](https://circleci.com/gh/Yinkozi/no-name-server.svg?style=svg&circle-token=c1d5d6c435da140a7eb83e80b75ff21b61886196)](https://circleci.com/gh/Yinkozi/no-name-server)


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
