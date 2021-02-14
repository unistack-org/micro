package flow

type node struct {
	name string
}

func (n *node) ID() string {
	return n.name
}

func (n *node) Name() string {
	return n.name
}

func (n *node) String() string {
	return n.name
}

func (n *node) Hashcode() interface{} {
	return n.name
}
