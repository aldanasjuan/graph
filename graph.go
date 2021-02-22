package graph

import (
	"strings"
)

//Node represents a node in the tree
type Node struct {
	Name     string
	Parent   *Node
	Children map[string]*Node
	Value    interface{}
}

//New root node
func New(name string, value interface{}) *Node {
	return &Node{
		Name:     name,
		Value:    value,
		Children: map[string]*Node{},
	}
}

//AddChild ...
func (n *Node) AddChild(name string, value interface{}) *Node {
	c := &Node{
		Name:     name,
		Parent:   n,
		Value:    value,
		Children: map[string]*Node{},
	}
	n.Children[c.Name] = c
	return c
}

//AddSibling ...
func (n *Node) AddSibling(name string, value interface{}) *Node {
	// fmt.Println(n.Parent.Name)
	if n.Parent != nil {
		c := &Node{
			Name:     name,
			Parent:   n.Parent,
			Value:    value,
			Children: map[string]*Node{},
		}
		n.Parent.Children[c.Name] = c
		return c
	}
	return nil
}

//AddParent adds a new node between the caller node and it's parent.
func (n *Node) AddParent(name string, value interface{}, all bool) *Node {
	parent := n.Parent
	newparent := &Node{
		Name:   name,
		Value:  value,
		Parent: parent,
	}
	n.Parent = newparent
	if all {
		newparent.Children = parent.Children
		parent.Children = map[string]*Node{
			newparent.Name: newparent,
		}
	} else {
		newparent.Children[n.Name] = n
		parent.Children[newparent.Name] = newparent
	}
	return newparent
}

//Root returns the root node from anywhere in the tree
func (n *Node) Root() *Node {
	if n.Parent != nil {
		return n.Parent.Root()
	} else {
		return n
	}
}

//Siblings returns a map with the siblings
func (n *Node) Siblings() map[string]*Node {
	res := map[string]*Node{}
	if n.Parent != nil {
		for k, v := range n.Parent.Children {
			if v != n {
				res[k] = v
			}
		}
	}
	return res
}

//Get a node from the tree using dot notation: "user.profile.config"
func (n *Node) Get(path string, global bool) *Node {
	// root := n.Root()
	addr := strings.Split(path, ".")
	var el *Node
	if global {
		el = n.Root()
		for _, key := range addr {
			if el == nil {
				break
			}
			if _, ok := el.Children[key]; ok {
				el = el.Children[key]
			} else {
				el = nil
			}
		}
	} else {
		el = n
		for _, key := range addr {
			if el == nil {
				break
			}
			if _, ok := el.Children[key]; ok {
				el = el.Children[key]
			} else {
				el = nil
			}
		}
	}

	return el
}
