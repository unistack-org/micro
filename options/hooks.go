package options // import "go.unistack.org/micro/v3/options"

// Hook func interface
type Hook interface{}

// Hooks func slice
type Hooks []Hook

// Append is used to add hooks
func (hs *Hooks) Append(h ...Hook) {
	*hs = append(*hs, h...)
}

// Replace is used to set hooks
func (hs *Hooks) Replace(h ...Hook) {
	*hs = h
}

// EachNext is used to itearate over hooks forward
func (hs *Hooks) EachNext(fn func(Hook)) {
	for idx := 0; idx < len(*hs); idx++ {
		fn((*hs)[idx])
	}
}

// EachPrev is used to iterate over hooks backward
func (hs *Hooks) EachPrev(fn func(Hook)) {
	for idx := len(*hs) - 1; idx >= 0; idx-- {
		fn((*hs)[idx])
	}
}
