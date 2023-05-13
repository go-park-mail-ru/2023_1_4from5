go test -coverprofile=coverage.out.tmp -coverpkg=./...  ./...
cat coverage.out.tmp | grep -v _mock.go | grep -v _easyjson.go | grep -v .pb.go | grep -v _grpc.go  > coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
