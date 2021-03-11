package codec

// Frame gives us the ability to define raw data to send over the pipes
type Frame struct {
	Data []byte
}
