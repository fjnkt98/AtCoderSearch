
.PHONY: gen-grpc-backend
gen-grpc-backend:
	@protoc --go_out=../backend/grpc --go_opt=paths=source_relative --go-grpc_out=../backend/grpc --go-grpc_opt=paths=source_relative atcodersearch.proto
