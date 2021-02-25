package main

import (
	"encoding/json"
	"fmt"

	"github.com/aldanasjuan/graph"
)

func main() {
	a := graph.New("a", "a")
	a.AddChild("b", "b").AddSibling("b2", map[string]string{"this": "rules"}).AddSibling("b3", "b3")
	a.AddChild("c", "c").AddSibling("c2", "c2").AddSibling("c3", "c3")
	a.AddChild("d", "d").AddSibling("d2", "d2").AddSibling("d3", "d3")
	b := a.Get("b", false)
	if b != nil {
		b.AddParent("a2", "a2", false)
	}
	j, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		fmt.Println("error:", err)

	}
	fmt.Println(string(j))

	//result from Marshaling
	t := `{
		"name": "a",
		"value": [1,2,3],
		"children": [
		  {
			"name": "a2",
			"value": "a2",
			"children": [
			  {
				"name": "b",
				"value": "b",
				"children": []
			  }
			]
		  },
		  {
			"name": "b2",
			"value": {
			  "this": "rules"
			},
			"children": []
		  },
		  {
			"name": "b3",
			"value": "b3",
			"children": []
		  },
		  {
			"name": "c",
			"value": "c",
			"children": []
		  },
		  {
			"name": "c2",
			"value": "c2",
			"children": []
		  },
		  {
			"name": "c3",
			"value": "c3",
			"children": []
		  },
		  {
			"name": "d",
			"value": "d",
			"children": []
		  },
		  {
			"name": "d2",
			"value": "d2",
			"children": []
		  },
		  {
			"name": "d3",
			"value": "d3",
			"children": []
		  }
		]
	  }`

	//recreate the tree from json
	z, err := graph.FromJSON([]byte(t))

	if err != nil {
		fmt.Println(err)
	}
	b = z.Get("a2.b", false)
	fmt.Println(b.Parent.Name)

	z.PrintAll()
	// fmt.Println("-------")
	// a.PrintAll()
}
