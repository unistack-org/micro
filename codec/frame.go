package codec

// Frame gives us the ability to define raw data to send over the pipes
type Frame struct {
	Data []byte
}

// NewFrame returns new frame with data
func NewFrame(data []byte) *Frame {
	return &Frame{Data: data}
}

// MarshalJSON returns frame data
func (m *Frame) MarshalJSON() ([]byte, error) {
	return m.Marshal()
}

// UnmarshalJSON set frame data
func (m *Frame) UnmarshalJSON(data []byte) error {
	return m.Unmarshal(data)
}

// MarshalYAML returns frame data
func (m *Frame) MarshalYAML() ([]byte, error) {
	return m.Marshal()
}

// UnmarshalYAML set frame data
func (m *Frame) UnmarshalYAML(data []byte) error {
	m.Data = append((m.Data)[0:0], data...)
	return nil
}

// ProtoMessage noop func
func (m *Frame) ProtoMessage() {}

// Reset resets frame
func (m *Frame) Reset() {
	*m = Frame{}
}

// String returns frame as string
func (m *Frame) String() string {
	return string(m.Data)
}

// Marshal returns frame data
func (m *Frame) Marshal() ([]byte, error) {
	return m.Data, nil
}

// Unmarshal set frame data
func (m *Frame) Unmarshal(data []byte) error {
	m.Data = data
	return nil
}
