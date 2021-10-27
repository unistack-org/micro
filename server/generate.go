package server

//go:generate sh -c "protoc -I./health -I../ -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out='components=micro|http|server',standalone=false,debug=true,paths=source_relative:./health health/health.proto"

import (

	// import required packages
	_ "go.unistack.org/micro-proto/v3/api"

	// import required packages
	_ "go.unistack.org/micro-proto/v3/openapiv3"
)
