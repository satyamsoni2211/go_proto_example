gen:
	protoc -Iproto proto/*.proto --go_out=pb --go-grpc_out=pb --go-grpc_opt=requireUnimplementedServers=false

client:
	go run cmd/client/client.go

server:
	go run cmd/server/server.go