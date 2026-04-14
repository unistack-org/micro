package broker

// IsValidHandler func signature
func IsValidHandler(sub any) error {
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
