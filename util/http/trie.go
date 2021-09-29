package http

import (
	"regexp"
	"strings"
	"sync"
)

// TrieOptions contains search options
type TrieOptions struct {
	IgnoreCase bool
}

// TrieOption func signature
type TrieOption func(*TrieOptions)

// IgnoreCase says that search must be case insensitive
func IgnoreCase(b bool) TrieOption {
	return func(o *TrieOptions) {
		o.IgnoreCase = b
	}
}

// Tree is a trie tree.
type Trie struct {
	node   *node
	rcache map[string]*regexp.Regexp
	rmu    sync.RWMutex
}

// node is a node of tree
type node struct {
	actions  map[string]interface{} // key is method, val is handler interface
	children map[string]*node       // key is label of next nodes
	label    string
}

const (
	pathRoot          string = "/"
	pathDelimiter     string = "/"
	paramDelimiter    string = ":"
	leftPtnDelimiter  string = "{"
	rightPtnDelimiter string = "}"
	ptnWildcard       string = "(.+)"
)

// NewTree creates a new trie tree.
func NewTrie() *Trie {
	return &Trie{
		node: &node{
			label:    pathRoot,
			actions:  make(map[string]interface{}),
			children: make(map[string]*node),
		},
		rcache: make(map[string]*regexp.Regexp),
	}
}

// Insert inserts a route definition to tree.
func (t *Trie) Insert(methods []string, path string, handler interface{}) {
	curNode := t.node
	if path == pathRoot {
		curNode.label = path
		for _, method := range methods {
			curNode.actions[method] = handler
		}
		return
	}
	ep := splitPath(path)
	for i, p := range ep {
		nextNode, ok := curNode.children[p]
		if ok {
			curNode = nextNode
		}
		// Create a new node.
		if !ok {
			curNode.children[p] = &node{
				label:    p,
				actions:  make(map[string]interface{}),
				children: make(map[string]*node),
			}
			curNode = curNode.children[p]
		}
		// last loop.
		// If there is already registered data, overwrite it.
		if i == len(ep)-1 {
			curNode.label = p
			for _, method := range methods {
				curNode.actions[method] = handler
			}
			break
		}
	}
}

// Search searches a path from a tree.
func (t *Trie) Search(method string, path string, opts ...TrieOption) (interface{}, map[string]string, bool) {
	params := make(map[string]string)

	options := TrieOptions{}
	for _, o := range opts {
		o(&options)
	}

	curNode := t.node

nodeLoop:
	for _, p := range splitPath(path) {
		nextNode, ok := curNode.children[p]
		if ok {
			curNode = nextNode
			continue nodeLoop
		}
		if options.IgnoreCase {
			// additional loop for case insensitive matching
			for k, v := range curNode.children {
				if literalEqual(k, p, true) {
					curNode = v
					continue nodeLoop
				}
			}
		}
		if len(curNode.children) == 0 {
			if !literalEqual(curNode.label, p, options.IgnoreCase) {
				// no matching path was found
				return nil, nil, false
			}
			break
		}
		isParamMatch := false
		for c := range curNode.children {
			if string([]rune(c)[0]) == leftPtnDelimiter {
				ptn := getPattern(c)
				t.rmu.RLock()
				reg, ok := t.rcache[ptn]
				t.rmu.RUnlock()
				if !ok {
					var err error
					reg, err = regexp.Compile(ptn)
					if err != nil {
						return nil, nil, false
					}
					t.rmu.Lock()
					t.rcache[ptn] = reg
					t.rmu.Unlock()
				}
				if reg.Match([]byte(p)) {
					pn := getParamName(c)
					params[pn] = p
					curNode = curNode.children[c]
					isParamMatch = true
					break
				}
				// no matching param was found.
				return nil, nil, false
			}
		}
		if !isParamMatch {
			return nil, nil, false
		}
	}
	if path == pathRoot {
		if len(curNode.actions) == 0 {
			return nil, nil, false
		}
	}

	handler, ok := curNode.actions[method]
	if !ok || handler == nil {
		return nil, nil, false
	}
	return handler, params, true
}

// getPattern gets a pattern from a label
// {id:[^\d+$]} -> ^\d+$
// {id}         -> (.+)
func getPattern(label string) string {
	leftI := strings.Index(label, leftPtnDelimiter)
	rightI := strings.Index(label, paramDelimiter)
	// if label doesn't have any pattern, return wild card pattern as default.
	if leftI == -1 || rightI == -1 {
		return ptnWildcard
	}
	return label[rightI+1 : len(label)-1]
}

// getParamName gets a parameter from a label
// {id:[^\d+$]} -> id
// {id}         -> id
func getParamName(label string) string {
	leftI := strings.Index(label, leftPtnDelimiter)
	rightI := func(l string) int {
		r := []rune(l)

		var n int

	loop:
		for i := 0; i < len(r); i++ {
			n = i
			switch string(r[i]) {
			case paramDelimiter:
				n = i
				break loop
			case rightPtnDelimiter:
				n = i
				break loop
			}

			if i == len(r)-1 {
				n = i + 1
				break loop
			}
		}

		return n
	}(label)

	return label[leftI+1 : rightI]
}

// splitPath removes an empty value in slice.
func splitPath(path string) []string {
	s := strings.Split(path, pathDelimiter)
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func literalEqual(component, literal string, ignoreCase bool) bool {
	if ignoreCase {
		return strings.EqualFold(component, literal)
	}
	return component == literal
}
