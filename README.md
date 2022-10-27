# ds2-handin-3

The server and client can be run with or without docker. 
By default, the client username is 'anonymous' and the address of the server is assumed to be localhost:50051.

## With Docker

Build the server image: `` docker build -t go-server -f DOCKERFILE.server .``

Run the server image: `` docker run --rm -p 50051:50051 go-server``


Build the client image: `` docker build -t go-client -f DOCKERFILE.client .``

Run the client image: `` docker run --rm -i --net=host go-client``
Pass arguments by appending the run command with `` -username "NAME" -sPort "ADDRESS OF SERVER" ``

## Without Docker

Run the server: `` go run server/server.go ``

Run the client: `` go run client/client.go -username "NAME" -sPort "ADDRESS OF SERVER"``
