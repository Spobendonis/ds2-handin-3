# ds2-handin-3

## Docker

Build the server image: `` docker build -t go-server -f DOCKERFILE.server .``

Run the server image: `` docker run --rm -p 50051:50051 go-server``


Build the client image: `` docker build -t go-client -f DOCKERFILE.client .``

Run the client image: `` docker run --rm -i --net=host go-client``
Pass arguments by appending the run command with `` -username "NAME" -sPort "ADDRESS" ``
