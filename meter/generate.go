package meter

//go:generate sh -c "protoc -I./handler -I../ -I$(go list -f '{{ .Dir }}' -m github.com/unistack-org/micro-proto) --go-micro_out='components=micro|http|server',standalone=false,debug=true,paths=source_relative:./handler handler/handler.proto"

import (
	_ "github.com/unistack-org/micro-proto/api"
	_ "github.com/unistack-org/micro-proto/openapiv2"
)
