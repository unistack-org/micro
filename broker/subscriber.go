package broker

// IsValidHandler func signature
func IsValidHandler(sub interface{}) error {
	switch sub.(type) {
	default:
		return ErrInvalidHandler
	case func(Message) error:
		break
	case func([]Message) error:
		break
	}
	return nil
}
