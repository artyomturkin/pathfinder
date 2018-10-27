package pathfinder

import (
	"fmt"
	"strings"
)

// Node route container
type Node struct {
	parent       *Node
	static       string
	children     map[string]*Node
	other        *Node
	parametrized bool
	params       []string

	path    string
	payload interface{}
}

// Add new subroute to node
func (n *Node) Add(path string, payload interface{}) error {
	ps := strings.Split(path, "/")
	if n.children == nil {
		n.children = map[string]*Node{}
	}
	params := []string{}

	current := n
	for pos, p := range ps {
		var next *Node
		set := false

		// Check and set parametrized
		if strings.HasPrefix(p, ":") {
			if current.other == nil {
				current.other = &Node{
					parent:       current,
					children:     map[string]*Node{},
					parametrized: true,
				}
			}
			next = current.other
			set = true
			params = append(params, p)
		}

		// Check and set if already exists
		if f, ok := current.children[p]; ok && !set {
			next = f
			set = true
		}

		// Create new static
		if !set {
			next = &Node{
				static:   p,
				parent:   current,
				children: map[string]*Node{},
			}
			current.children[p] = next
		}

		// Continue if not last one
		if pos != len(ps)-1 {
			current = next
			continue
		}

		if next.path != "" && next.path != path {
			return fmt.Errorf("path '%s' conflicts with registered '%s'", path, next.path)
		}
		if next.path == path {
			next.payload = payload
			break
		}
		next.path = path
		next.payload = payload
		next.params = params
		break
	}
	return nil
}

// Lookup path in node and subnodes
func (n *Node) Lookup(path string) (interface{}, map[string]string, error) {
	ps := strings.Split(path, "/")
	params := []string{}

	current := n
	for pos, p := range ps {

		var next *Node
		set := false

		if f, ok := current.children[p]; ok {
			next = f
			set = true
		}
		if !set && current.other != nil {
			next = current.other
			params = append(params, p)
			set = true
		}
		if set && pos != len(ps)-1 {
			current = next
			continue
		}

		if !set || next.path == "" {
			return nil, nil, fmt.Errorf("Not found")
		}
		current = next
	}
	p := map[string]string{}
	for i, pa := range params {
		p[current.params[i]] = pa
	}
	return current.payload, p, nil
}
