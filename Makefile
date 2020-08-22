gen:
	protoc -Iproto proto/*.proto --go_out=pb --go-grpc_out=pb --go-grpc_opt=requireUnimplementedServers=false