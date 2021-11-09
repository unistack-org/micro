package http

// Radix tree implementation below is a based on the original work by
// Armon Dadgar in https://github.com/armon/go-radix/blob/master/radix.go
// (MIT licensed). It's been heavily modified for use as a HTTP routing tree.
// Copied from chi mux tree.go https://raw.githubusercontent.com/go-chi/chi/master/tree.go

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type methodTyp uint

const (
	mSTUB methodTyp = 1 << iota
	mCONNECT
	mDELETE
	mGET
	mHEAD
	mOPTIONS
	mPATCH
	mPOST
	mPUT
	mTRACE
)

var mALL = mCONNECT | mDELETE | mGET | mHEAD |
	mOPTIONS | mPATCH | mPOST | mPUT | mTRACE

var methodMap = map[string]methodTyp{
	http.MethodConnect: mCONNECT,
	http.MethodDelete:  mDELETE,
	http.MethodGet:     mGET,
	http.MethodHead:    mHEAD,
	http.MethodOptions: mOPTIONS,
	http.MethodPatch:   mPATCH,
	http.MethodPost:    mPOST,
	http.MethodPut:     mPUT,
	http.MethodTrace:   mTRACE,
}

// RegisterMethod adds support for custom HTTP method handlers, available
// via Router#Method and Router#MethodFunc
func RegisterMethod(method string) error {
	if method == "" {
		return nil
	}
	method = strings.ToUpper(method)
	if _, ok := methodMap[method]; ok {
		return nil
	}
	n := len(methodMap)
	if n > strconv.IntSize-2 {
		return fmt.Errorf("max number of methods reached (%d)", strconv.IntSize)
	}
	mt := methodTyp(2 << n)
	methodMap[method] = mt
	mALL |= mt

	return nil
}

type nodeTyp uint8

const (
	ntStatic   nodeTyp = iota // /home
	ntRegexp                  // /{id:[0-9]+}
	ntParam                   // /{user}
	ntCatchAll                // /api/v1/*
)

func NewTrie() *Node {
	return &Node{typ: ntStatic}
}

type Node struct {
	// regexp matcher for regexp nodes
	rex *regexp.Regexp

	// HTTP handler endpoints on the leaf node
	endpoints endpoints

	// prefix is the common prefix we ignore
	prefix string

	// child nodes should be stored in-order for iteration,
	// in groups of the node type.
	children [ntCatchAll + 1]nodes

	// first byte of the child prefix
	tail byte

	// node type: static, regexp, param, catchAll
	typ nodeTyp

	// first byte of the prefix
	label byte
}

// endpoints is a mapping of http method constants to handlers
// for a given route.
type endpoints map[methodTyp]*endpoint

type endpoint struct {
	// parameters keys recorded on handler nodes
	paramKeys []string
	// endpoint handler
	handler interface{}
	// pattern is the routing pattern for handler nodes
	pattern string
}

func (s endpoints) Value(method methodTyp) *endpoint {
	mh, ok := s[method]
	if !ok {
		mh = &endpoint{}
		s[method] = mh
	}

	return mh
}

func (n *Node) Insert(methods []string, pattern string, handler interface{}) error {
	var err error
	for _, method := range methods {
		if err = n.insert(methodMap[method], pattern, handler); err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) insert(method methodTyp, pattern string, handler interface{}) error {
	var parent *Node
	search := pattern

	for {
		// Handle key exhaustion
		if len(search) == 0 {
			// Insert or update the node's leaf handler
			return n.setEndpoint(method, handler, pattern)
		}

		// We're going to be searching for a wild node next,
		// in this case, we need to get the tail
		label := search[0]
		var segTail byte
		var segEndIdx int
		var segTyp nodeTyp
		var segRexpat string
		var err error
		if label == '{' || label == '*' {
			segTyp, _, segRexpat, segTail, _, segEndIdx, err = patNextSegment(search)
		}
		if err != nil {
			return err
		}

		var prefix string
		if segTyp == ntRegexp {
			prefix = segRexpat
		}

		// Look for the edge to attach to
		parent = n
		n = n.getEdge(segTyp, label, segTail, prefix)

		// No edge, create one
		if n == nil {
			child := &Node{typ: ntStatic, label: label, tail: segTail, prefix: search}
			var hn *Node
			hn, err = parent.addChild(child, search)
			if err != nil {
				return err
			}
			return hn.setEndpoint(method, handler, pattern)
		}

		// Found an edge to match the pattern

		if n.typ > ntStatic {
			// We found a param node, trim the param from the search path and continue.
			// This param/wild pattern segment would already be on the tree from a previous
			// call to addChild when creating a new node.
			search = search[segEndIdx:]
			continue
		}

		// Static nodes fall below here.
		// Determine longest prefix of the search key on match.
		commonPrefix := longestPrefix(search, n.prefix)
		if commonPrefix == len(n.prefix) {
			// the common prefix is as long as the current node's prefix we're attempting to insert.
			// keep the search going.
			search = search[commonPrefix:]
			continue
		}

		// Split the node
		child := &Node{
			typ:    ntStatic,
			prefix: search[:commonPrefix],
		}
		if err = parent.replaceChild(search[0], segTail, child); err != nil {
			return err
		}

		// Restore the existing node
		n.label = n.prefix[commonPrefix]
		n.prefix = n.prefix[commonPrefix:]
		if _, err = child.addChild(n, n.prefix); err != nil {
			return err
		}

		// If the new key is a subset, set the method/handler on this node and finish.
		search = search[commonPrefix:]
		if len(search) == 0 {
			return child.setEndpoint(method, handler, pattern)
		}

		// Create a new edge for the node
		subchild := &Node{
			typ:    ntStatic,
			label:  search[0],
			prefix: search,
		}
		var hn *Node
		hn, err = child.addChild(subchild, search)
		if err != nil {
			return err
		}
		return hn.setEndpoint(method, handler, pattern)
	}
}

// addChild appends the new `child` node to the tree using the `pattern` as the trie key.
// For a URL router like chi's, we split the static, param, regexp and wildcard segments
// into different nodes. In addition, addChild will recursively call itself until every
// pattern segment is added to the url pattern tree as individual nodes, depending on type.
func (n *Node) addChild(child *Node, prefix string) (*Node, error) {
	search := prefix

	// handler leaf node added to the tree is the child.
	// this may be overridden later down the flow
	hn := child

	// Parse next segment
	segTyp, _, segRexpat, segTail, segStartIdx, segEndIdx, err := patNextSegment(search)
	if err != nil {
		return nil, err
	}
	// Add child depending on next up segment
	switch segTyp {

	case ntStatic:
		// Search prefix is all static (that is, has no params in path)
		// noop

	default:
		// Search prefix contains a param, regexp or wildcard

		if segTyp == ntRegexp {
			rex, err := regexp.Compile(segRexpat)
			if err != nil {
				return nil, fmt.Errorf("invalid regexp pattern '%s' in route param", segRexpat)
			}
			child.prefix = segRexpat
			child.rex = rex
		}

		if segStartIdx == 0 {
			// Route starts with a param
			child.typ = segTyp

			if segTyp == ntCatchAll {
				segStartIdx = -1
			} else {
				segStartIdx = segEndIdx
			}
			if segStartIdx < 0 {
				segStartIdx = len(search)
			}
			child.tail = segTail // for params, we set the tail

			if segStartIdx != len(search) {
				// add static edge for the remaining part, split the end.
				// its not possible to have adjacent param nodes, so its certainly
				// going to be a static node next.

				search = search[segStartIdx:] // advance search position

				nn := &Node{
					typ:    ntStatic,
					label:  search[0],
					prefix: search,
				}
				var err error
				hn, err = child.addChild(nn, search)
				if err != nil {
					return nil, err
				}
			}

		} else if segStartIdx > 0 {
			// Route has some param

			// starts with a static segment
			child.typ = ntStatic
			child.prefix = search[:segStartIdx]
			child.rex = nil

			// add the param edge node
			search = search[segStartIdx:]

			nn := &Node{
				typ:   segTyp,
				label: search[0],
				tail:  segTail,
			}
			var err error
			hn, err = child.addChild(nn, search)
			if err != nil {
				return nil, err
			}
		}
	}

	n.children[child.typ] = append(n.children[child.typ], child)
	n.children[child.typ].Sort()
	return hn, nil
}

func (n *Node) replaceChild(label, tail byte, child *Node) error {
	for i := 0; i < len(n.children[child.typ]); i++ {
		if n.children[child.typ][i].label == label && n.children[child.typ][i].tail == tail {
			n.children[child.typ][i] = child
			n.children[child.typ][i].label = label
			n.children[child.typ][i].tail = tail
			return nil
		}
	}
	return fmt.Errorf("replacing missing child")
}

func (n *Node) getEdge(ntyp nodeTyp, label, tail byte, prefix string) *Node {
	nds := n.children[ntyp]
	for i := 0; i < len(nds); i++ {
		if nds[i].label == label && nds[i].tail == tail {
			if ntyp == ntRegexp && nds[i].prefix != prefix {
				continue
			}
			return nds[i]
		}
	}
	return nil
}

func (n *Node) setEndpoint(method methodTyp, handler interface{}, pattern string) error {
	// Set the handler for the method type on the node
	if n.endpoints == nil {
		n.endpoints = make(endpoints)
	}

	paramKeys, err := patParamKeys(pattern)
	if err != nil {
		return err
	}
	if method&mSTUB == mSTUB {
		n.endpoints.Value(mSTUB).handler = handler
	}
	if method&mALL == mALL {
		h := n.endpoints.Value(mALL)
		h.handler = handler
		h.pattern = pattern
		h.paramKeys = paramKeys
		for _, m := range methodMap {
			h := n.endpoints.Value(m)
			h.handler = handler
			h.pattern = pattern
			h.paramKeys = paramKeys
		}
	} else {
		h := n.endpoints.Value(method)
		h.handler = handler
		h.pattern = pattern
		h.paramKeys = paramKeys
	}

	return nil
}

func (n *Node) Search(method string, path string) (interface{}, map[string]string, bool) {
	params := &routeParams{}
	// Find the routing handlers for the path
	rn := n.findRoute(params, methodMap[method], path)
	if rn == nil {
		return nil, nil, false
	}
	ep, ok := rn.endpoints[methodMap[method]]
	if !ok {
		return nil, nil, false
	}

	eparams := make(map[string]string, len(params.keys))
	for idx, key := range params.keys {
		eparams[key] = params.vals[idx]
	}

	return ep.handler, eparams, true
}

type routeParams struct {
	keys []string
	vals []string
}

// Recursive edge traversal by checking all nodeTyp groups along the way.
// It's like searching through a multi-dimensional radix trie.
func (n *Node) findRoute(params *routeParams, method methodTyp, path string) *Node {
	nn := n
	search := path

	for t, nds := range nn.children {
		ntyp := nodeTyp(t)
		if len(nds) == 0 {
			continue
		}

		var xn *Node
		xsearch := search

		var label byte
		if search != "" {
			label = search[0]
		}

		switch ntyp {
		case ntStatic:
			xn = nds.findEdge(label)
			if xn == nil || !strings.HasPrefix(xsearch, xn.prefix) {
				continue
			}
			xsearch = xsearch[len(xn.prefix):]

		case ntParam, ntRegexp:
			// short-circuit and return no matching route for empty param values
			if xsearch == "" {
				continue
			}

			// serially loop through each node grouped by the tail delimiter
			for idx := 0; idx < len(nds); idx++ {
				xn = nds[idx]

				// label for param nodes is the delimiter byte
				p := strings.IndexByte(xsearch, xn.tail)

				if p < 0 {
					if xn.tail == '/' {
						p = len(xsearch)
					} else {
						continue
					}
				} else if ntyp == ntRegexp && p == 0 {
					continue
				}

				if ntyp == ntRegexp && xn.rex != nil {
					if !xn.rex.MatchString(xsearch[:p]) {
						continue
					}
				} else if strings.IndexByte(xsearch[:p], '/') != -1 {
					// avoid a match across path segments
					continue
				}

				prevlen := len(params.vals)
				params.vals = append(params.vals, xsearch[:p])
				xsearch = xsearch[p:]

				if len(xsearch) == 0 {
					if xn.isLeaf() {
						h := xn.endpoints[method]
						if h != nil && h.handler != nil {
							params.keys = append(params.keys, h.paramKeys...)
							return xn
						}
					}
				}

				// recursively find the next node on this branch
				fin := xn.findRoute(params, method, xsearch)
				if fin != nil {
					return fin
				}

				// not found on this branch, reset vars
				params.vals = params.vals[:prevlen]
				xsearch = search
			}

			params.vals = append(params.vals, "")

		default:
			// catch-all nodes
			params.vals = append(params.vals, search)
			xn = nds[0]
			xsearch = ""
		}

		if xn == nil {
			continue
		}

		// did we find it yet?
		if len(xsearch) == 0 {
			if xn.isLeaf() {
				h := xn.endpoints[method]
				if h != nil && h.handler != nil {
					params.keys = append(params.keys, h.paramKeys...)
					return xn
				}
			}
		}

		// recursively find the next node..
		fin := xn.findRoute(params, method, xsearch)
		if fin != nil {
			return fin
		}

		// Did not find final handler, let's remove the param here if it was set
		if xn.typ > ntStatic {
			if len(params.vals) > 0 {
				params.vals = params.vals[:len(params.vals)-1]
			}
		}

	}

	return nil
}

func (n *Node) isLeaf() bool {
	return n.endpoints != nil
}

// patNextSegment returns the next segment details from a pattern:
// node type, param key, regexp string, param tail byte, param starting index, param ending index
func patNextSegment(pattern string) (nodeTyp, string, string, byte, int, int, error) {
	ps := strings.Index(pattern, "{")
	ws := strings.Index(pattern, "*")

	if ps < 0 && ws < 0 {
		return ntStatic, "", "", 0, 0, len(pattern), nil // we return the entire thing
	}

	// Sanity check
	if ps >= 0 && ws >= 0 && ws < ps {
		return ntStatic, "", "", 0, 0, 0, fmt.Errorf("wildcard '*' must be the last pattern in a route, otherwise use a '{param}'")
	}

	var tail byte = '/' // Default endpoint tail to / byte

	if ps < 0 {
		// Wildcard pattern as finale
		if ws < len(pattern)-1 {
			return ntStatic, "", "", 0, 0, 0, fmt.Errorf("wildcard '*' must be the last value in a route. trim trailing text or use a '{param}' instead")
		}

		return ntCatchAll, "*", "", 0, ws, len(pattern), nil
	}

	// Param/Regexp pattern is next
	nt := ntParam

	// Read to closing } taking into account opens and closes in curl count (cc)
	cc := 0
	pe := ps
	for i, c := range pattern[ps:] {
		if c == '{' {
			cc++
		} else if c == '}' {
			cc--
			if cc == 0 {
				pe = ps + i
				break
			}
		}
	}

	if pe == ps {
		return ntStatic, "", "", 0, 0, 0, fmt.Errorf("route param closing delimiter '}' is missing")
	}

	key := pattern[ps+1 : pe]
	pe++ // set end to next position

	if pe < len(pattern) {
		tail = pattern[pe]
	}

	var rexpat string
	if idx := strings.Index(key, ":"); idx >= 0 {
		nt = ntRegexp
		rexpat = key[idx+1:]
		key = key[:idx]
	}

	if len(rexpat) > 0 {
		if rexpat[0] != '^' {
			rexpat = "^" + rexpat
		}
		if rexpat[len(rexpat)-1] != '$' {
			rexpat += "$"
		}
	}

	return nt, key, rexpat, tail, ps, pe, nil
}

func patParamKeys(pattern string) ([]string, error) {
	pat := pattern
	paramKeys := []string{}
	for {
		ptyp, paramKey, _, _, _, e, err := patNextSegment(pat)
		if err != nil {
			return nil, err
		}
		if ptyp == ntStatic {
			return paramKeys, nil
		}
		for i := 0; i < len(paramKeys); i++ {
			if paramKeys[i] == paramKey {
				return nil, fmt.Errorf("routing pattern '%s' contains duplicate param key, '%s'", pattern, paramKey)
			}
		}
		paramKeys = append(paramKeys, paramKey)
		pat = pat[e:]
	}
}

// longestPrefix finds the length of the shared prefix
// of two strings
func longestPrefix(k1, k2 string) int {
	max := len(k1)
	if l := len(k2); l < max {
		max = l
	}
	var i int
	for i = 0; i < max; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return i
}

type nodes []*Node

// Sort the list of nodes by label
func (ns nodes) Sort()              { sort.Sort(ns); ns.tailSort() }
func (ns nodes) Len() int           { return len(ns) }
func (ns nodes) Swap(i, j int)      { ns[i], ns[j] = ns[j], ns[i] }
func (ns nodes) Less(i, j int) bool { return ns[i].label < ns[j].label }

// tailSort pushes nodes with '/' as the tail to the end of the list for param nodes.
// The list order determines the traversal order.
func (ns nodes) tailSort() {
	for i := len(ns) - 1; i >= 0; i-- {
		if ns[i].typ > ntStatic && ns[i].tail == '/' {
			ns.Swap(i, len(ns)-1)
			return
		}
	}
}

func (ns nodes) findEdge(label byte) *Node {
	num := len(ns)
	idx := 0
	i, j := 0, num-1
	for i <= j {
		idx = i + (j-i)/2
		if label > ns[idx].label {
			i = idx + 1
		} else if label < ns[idx].label {
			j = idx - 1
		} else {
			i = num // breaks cond
		}
	}

	if ns[idx].label != label {
		return nil
	}

	return ns[idx]
}
