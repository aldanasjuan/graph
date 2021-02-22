package main

import (
	"github.com/aldanasjuan/graph"
)

func main() {
	a := graph.New("a", "a")
	a.AddChild("b", "b").AddChild("b2", "b2").AddChild("b3", "b3")
	a.AddChild("c", "c").AddChild("c2", "c2").AddChild("c3", "c3")
	a.AddChild("d", "d").AddChild("d2", "d2").AddChild("d3", "d3")
	b := a.Get("b", false)
	b.AddParent("a2", "a2", false)

	a.PrintAll()
}
