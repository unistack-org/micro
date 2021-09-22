package codec

// Frame gives us the ability to define raw data to send over the pipes
type Frame struct {
	Data []byte
}

func (m *Frame) MarshalJSON() ([]byte, error) {
	return m.Data, nil
}

func (m *Frame) UnmarshalJSON(data []byte) error {
	m.Data = data
	return nil
}

func (m *Frame) ProtoMessage() {}

func (m *Frame) Reset() {
	*m = Frame{}
}

func (m *Frame) String() string {
	return string(m.Data)
}

func (m *Frame) Marshal() ([]byte, error) {
	return m.Data, nil
}

func (m *Frame) Unmarshal(data []byte) error {
	m.Data = data
	return nil
}
