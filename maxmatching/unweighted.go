package maxmatching

import "fmt"

//Given an unweighted bipartite graph, find the maximum matching
//(with most possible points connected and connected only once)

/*
A       0
B		1
C		2
D		3

A->0,1,2,3
B->2,3
C->0
D->2

0->1
1->3
2->0
3->2

*/

/*
edges[i]: all points connected to i
m: total number of points on right side (total on left side is len(edges))
*/
func MaxMatching(edges [][]int, m int) int {
	from := make([]int, m)
	links := 0
	for i := range from {
		from[i] = -1
	}
	for i := range edges {
		if dfs(edges, i, from, map[int]bool{}) {
			links++
		}
		fmt.Printf("from:%v\n", from)
	}
	for i := 0; i < len(from); i++ {
		fmt.Printf("[%v]->[%v]\n", from[i], i)
	}
	fmt.Printf("links:%v\n", links)
	return links
}

func dfs(edges [][]int, left int, from []int, visiting map[int]bool) bool {
	for _, right := range edges[left] {
		if visiting[right] {
			continue
		}
		visiting[right] = true
		if from[right] == -1 || dfs(edges, from[right], from, visiting) {
			from[right] = left
			return true
		}
	}
	return false
}
