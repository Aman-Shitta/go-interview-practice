package main

import (
	"fmt"
	"slices"
	"sync"
)

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.

	if numWorkers <= 0 {
		return nil
	}

	semaphores := make(chan bool, numWorkers)
	bfsResChan := make(chan struct {
		query  int
		result []int
	}, len(queries))

	bfsRes := make(map[int][]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, q := range queries {
		semaphores <- true

		wg.Add(1)
		go func(query int) {
			defer func() {
				<-semaphores
				wg.Done()
			}()

			res := bfs(query, graph)

			bfsResChan <- struct {
				query  int
				result []int
			}{query: query, result: res}
		}(q)
	}

	close(semaphores)
	wg.Wait()
	close(bfsResChan)

	for res := range bfsResChan {
		mu.Lock()
		bfsRes[res.query] = res.result
		mu.Unlock()
	}

	// fmt.Println("bfsRes :: ", bfsRes)
	return bfsRes

}

func main() {
	// You can insert optional local tests here if desired.
	graph := map[int][]int{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {4},
		4: {},
	}
	queries := []int{0, 1, 2}
	numWorkers := 2

	results := ConcurrentBFSQueries(graph, queries, numWorkers)

	fmt.Println(results)
}

func bfs(node int, graph map[int][]int) []int {

	var visitedOrder []int
	var stack, processed []int

	stack = append(stack, node)

	for len(stack) > 0 {

		currentNode := stack[0]
		stack = stack[1:]

		if slices.Contains(processed, currentNode) {
			continue
		}

		visitedOrder = append(visitedOrder, currentNode)
		processed = append(processed, currentNode)
		edges := graph[currentNode]

		for _, edge := range edges {
			if !slices.Contains(processed, edge) {
				stack = append(stack, edge)
			}
		}

	}
	return visitedOrder
}
