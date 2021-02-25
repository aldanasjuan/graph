package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fasttemplate"
)

//Node represents a node in the tree
type Node struct {
	Name     string      `json:"name,omitempty"`
	Parent   *Node       `json:"-"`
	Children []*Node     `json:"children,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

//New root node
func New(name string, value interface{}) *Node {
	return &Node{
		Name:     name,
		Value:    value,
		Children: []*Node{},
	}
}

func (n *Node) appendParent(p *Node) {
	n.Parent = p
	for _, c := range n.Children {
		c.appendParent(n)
	}
}

//AddChild ...
func (n *Node) AddChild(name string, value interface{}) *Node {
	c := &Node{
		Name:     name,
		Parent:   n,
		Value:    value,
		Children: []*Node{},
	}
	n.Children = append(n.Children, c)
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
			Children: []*Node{},
		}
		n.Parent.Children = append(n.Parent.Children, c)
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
		parent.Children = []*Node{newparent}
	} else {
		newparent.Children = []*Node{n}
		// parent.Children = append(parent.Children, newparent)
		children := []*Node{newparent}
		for _, c := range parent.Children {
			if c != n {
				children = append(children, c)
			}
		}
		parent.Children = children
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
func (n *Node) Siblings() []*Node {
	res := []*Node{}
	if n.Parent != nil {
		for _, c := range n.Parent.Children {
			if c != n {
				res = append(res, c)
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
			index := -1
			for i, c := range el.Children {
				if c.Name == key {
					index = i
					break
				}
			}
			if index == -1 {
				el = nil
			} else {
				el = el.Children[index]
			}
		}
	} else {
		el = n
		for _, key := range addr {
			if el == nil {
				break
			}
			index := -1
			for i, c := range el.Children {
				if c.Name == key {
					index = i
					break
				}
			}
			if index == -1 {
				el = nil
			} else {
				el = el.Children[index]
			}
		}
	}

	return el
}

//PrintAll prints the tree
func (n *Node) PrintAll() {
	r := n.Root()
	r.print(0)
}

//Print from the node down
func (n *Node) Print() {
	n.print(0)
}

func (n *Node) print(level int) {
	ident := strings.Join(make([]string, level+1), "  ")
	fmt.Printf("%v- name: '%v' value: %v\n", ident, n.Name, n.Value)
	for _, v := range n.Children {
		v.print(level + 1)
	}
}

//MarshalJSON ...
func (n *Node) MarshalJSON() ([]byte, error) {
	tmp := `{"name": "<name>", "value": <value>, "children": [<children>]}`
	w, _ := fasttemplate.NewTemplate(tmp, "<", ">")
	b := bytes.Buffer{}
	w.ExecuteFunc(&b, n.writeTag)
	// fmt.Println(b.String())
	return b.Bytes(), nil
}

func (n *Node) writeTag(w io.Writer, tag string) (int, error) {
	switch tag {
	case "name":
		return w.Write([]byte(n.Name))
	case "value":
		j, err := json.Marshal(n.Value)
		if err != nil {
			return 0, err
		}
		return w.Write(j)
	case "children":
		res := [][]byte{}
		for _, c := range n.Children {
			bts, _ := c.MarshalJSON()
			res = append(res, bts)
		}
		return w.Write(bytes.Join(res, []byte(",\n")))
	}
	return 0, fmt.Errorf("failed to write tag, got %v", tag)
}

// FromJSON unmarshals from json, links the nodes together and returns the root node
func FromJSON(b []byte) (*Node, error) {
	n := &Node{}
	err := json.Unmarshal(b, n)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Children {
		c.appendParent(n)
	}
	return n, nil
}
