#protoc --go_out=plugins=grpc:./ --go_opt=./ health_check.proto
#protoc --go_out=plugins=grpc:./ --go_opt=paths=./ --go-grpc_out=./ --go-grpc_opt=paths=./ ./health_check.proto

#protoc --go_out=plugins=grpc:./ --go_opt=paths=./ ./health_check.proto
#protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./message.proto
protoc --go_out=plugins=grpc:. *.proto
