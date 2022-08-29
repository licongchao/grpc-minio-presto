
# Golang
<!-- 
    go get google.golang.org/grpc
    go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
 -->
$ protoc --go_out=./go_gen --go-grpc_out=./go_gen ./*.proto

# Python

% pip install grpclib protobuf
% pip install grpcio  grpcio-tools -i http://pypi.douban.com/simple/ --default-timeout=99999 --trusted-host pypi.douban.com

$ python -m grpc_tools.protoc --python_out=./py_gen/ --grpc_python_out=./py_gen/ -I. modelopr.proto
