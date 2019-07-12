package main

import "fmt"

type Node struct {
	ID       string
	Parents  []*Node
	Children []*Node
	ch       chan struct{}
}

func (n *Node) Close() {
	fmt.Printf("closing %v\n", n.ID)
	close(n.ch)
}

func (n *Node) Done() chan struct{} {
	return n.ch
}

func DFS(any *Node) {

}

func main() {
	/*construct nodes

	A->B
	   ↓
	C->D->E->F

	expected (any):
	(A, B, C), D, E, F
	(C, A, B),
	(A, C, B),

	*/
	A := &Node{ID: "A", ch: make(chan struct{})}
	B := &Node{ID: "B", ch: make(chan struct{})}
	C := &Node{ID: "C", ch: make(chan struct{})}
	D := &Node{ID: "D", ch: make(chan struct{})}
	E := &Node{ID: "E", ch: make(chan struct{})}
	F := &Node{ID: "F", ch: make(chan struct{})}

	A.Children = []*Node{B}
	B.Parents = []*Node{A}
	B.Children = []*Node{D}
	C.Children = []*Node{D}
	D.Parents = []*Node{B, C}
	D.Children = []*Node{E}
	E.Parents = []*Node{D}
	E.Children = []*Node{F}
	F.Parents = []*Node{E}

	DFS(F)
}
