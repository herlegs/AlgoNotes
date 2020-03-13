package main

import (
	"fmt"
	"sync"
)

type Closable interface {
	Shutdown() error
	Done() chan struct{}
}

type Node struct {
	ID             string
	Parents        map[*Node]bool
	Children       map[*Node]bool
	DataCh         chan interface{}
	wg             sync.WaitGroup
	externalSource Closable
	doneCh         chan struct{}
	closeOnce      sync.Once
}

func NewNodeFromID(id string) *Node {
	return &Node{
		ID:       id,
		Parents:  map[*Node]bool{},
		Children: map[*Node]bool{},
		DataCh:   make(chan interface{}),
		doneCh:   make(chan struct{}),
	}
}

// ClosePipeline here we define the function as: close the data flow graph
func (n *Node) ClosePipeline() {
	Topological(n, func(node *Node) {
		node.close()
	})
}

// close is the internal function to close resources
// currently we don't see any necessity of allow user to manipulate close of this stream collection only
// hence make it private
func (n *Node) close() {
	n.closeOnce.Do(func() {
		fmt.Printf("closing %v\n", n.ID)
		//closing logic of all the resources of this stream
		//in case it is relying on has external data source, eg. a stream reader
		if n.externalSource != nil {
			_ = n.externalSource.Shutdown()
			<-n.externalSource.Done()
		}
		//here we should assume all parents/external source that will insert to data channel are closed when reached here
		//that means the following check is not necessary. eg:
		/*
			for p := range n.Parents {
				<-p.done()
			}
		*/
		close(n.DataCh)
		//when data channel is closed, wait all routine/job replying on it to auto-close
		//eg. Map/Filter function's conversion job routine
		n.wg.Wait()
		close(n.doneCh)
		fmt.Printf("closed %v\n", n.ID)
	})
}

//DonePipeline here we define as checking the entire graph is closed or not
//currently we don't see any necessity of a blocking function for waiting this stream collection to close
//if there is we can use another function
func (n *Node) PipelineDone() chan struct{} {
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		Topological(n, func(node *Node) {
			<-node.done()
		})
	}()
	return doneCh
}

func (n *Node) done() chan struct{} {
	return n.doneCh
}

//Topological is the function for graph traversal
func Topological(any *Node, fn func(*Node)) {
	visited := map[*Node]bool{}
	getVisitGraph(any, visited)
	for node, v := range visited {
		if !v {
			dfs(node, visited, fn)
		}
	}
}

func getVisitGraph(n *Node, graph map[*Node]bool) {
	if _, exist := graph[n]; exist {
		return
	}
	graph[n] = false
	for p := range n.Parents {
		getVisitGraph(p, graph)
	}
	for c := range n.Children {
		getVisitGraph(c, graph)
	}
}

func dfs(n *Node, visited map[*Node]bool, fn func(*Node)) {
	if visited[n] {
		return
	}
	visited[n] = true
	for p := range n.Parents {
		dfs(p, visited, fn)
	}
	fn(n)
}

func main() {
	/*construct nodes

	A->B
	   ↓
	C->D->E->F

	expected (any):
	(A, B, C),
	(C, A, B),
	(A, C, B), D, E, F

	*/
	nodeMap := buildGraph([][2]string{
		{"A", "B"}, {"B", "D"}, {"C", "D"}, {"D", "E"}, {"E", "F"},
	})

	node := nodeMap["A"]

	go func() {
		node.ClosePipeline()
	}()

	<-node.PipelineDone()
}

func buildGraph(edges [][2]string) map[string]*Node {
	nodeMap := map[string]*Node{}

	for _, e := range edges {
		from, to := e[0], e[1]
		fromNode, exist := nodeMap[from]
		if !exist {
			fromNode = NewNodeFromID(from)
			nodeMap[from] = fromNode
		}
		toNode, exist := nodeMap[to]
		if !exist {
			toNode = NewNodeFromID(to)
			nodeMap[to] = toNode
		}
		fromNode.Children[toNode] = true
		toNode.Parents[fromNode] = true
		go func() {

		}()
	}

	return nodeMap
}
