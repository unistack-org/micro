package server

//go:generate protoc -I./health -I../ -I/home/vtolstov/.cache/go-path/pkg/mod/github.com/unistack-org/micro-proto@v0.0.1 --micro_out=components=micro|http|server,standalone=false,debug=true,paths=source_relative:./health health/health.proto
