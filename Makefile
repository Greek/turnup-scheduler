.PHONY: proto

proto-fmt:
	 buf format -w proto/scheduler.proto

proto-compile:
	protoc --go_out=. --go_opt=paths=source_relative \
	    --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/scheduler.proto

gen:
	export PATH="$$PATH:$$(go env GOPATH)/bin" && \
	make proto-fmt && \
	make proto-compile && \
	echo "done"

run:
	make gen && go run .

build:
	make gen && go build .
