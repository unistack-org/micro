// Package metadata is a way of defining message headers
package metadata

var (
	// HeaderTopic is the header name that contains topic name.
	HeaderTopic = "Micro-Topic"
	// HeaderContentType specifies content type of message.
	HeaderContentType = "Content-Type"
	// HeaderEndpoint specifies endpoint in service.
	HeaderEndpoint = "Micro-Endpoint"
	// HeaderService specifies service.
	HeaderService = "Micro-Service"
	// HeaderTimeout specifies timeout of operation.
	HeaderTimeout = "Micro-Timeout"
	// HeaderAuthorization specifies Authorization header.
	HeaderAuthorization = "Authorization"
	// HeaderXRequestID specifies request id.
	HeaderXRequestID = "X-Request-Id"
)
