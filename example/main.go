package main

import (
	"fmt"

	"github.com/aldanasjuan/graph"
)

func main() {
	a := graph.New("a", "a")
	a.AddChild("b", "b").AddChild("c", "c").AddChild("d", "d")
	d := a.Get("b.c.d", false)

	fmt.Println(d)
}
