package graph

import (
	"fmt"

	"github.com/moorara/algo/ds/list"
	"github.com/moorara/algo/pkg/graphviz"
)

// Undirected represents an undirected graph data type.
type Undirected struct {
	v, e int
	adj  [][]int
}

// NewUndirected creates a new undirected graph.
func NewUndirected(V int, edges ...[2]int) *Undirected {
	adj := make([][]int, V)
	for i := range adj {
		adj[i] = make([]int, 0)
	}

	g := &Undirected{
		v:   V, // no. of vertices
		e:   0, // no. of edges
		adj: adj,
	}

	for _, e := range edges {
		g.AddEdge(e[0], e[1])
	}

	return g
}

// V returns the number of vertices.
func (g *Undirected) V() int {
	return g.v
}

// E returns the number of edges.
func (g *Undirected) E() int {
	return g.e
}

func (g *Undirected) isVertexValid(v int) bool {
	return v >= 0 && v < g.v
}

// Degree returns the degree of a vertext.
func (g *Undirected) Degree(v int) int {
	if !g.isVertexValid(v) {
		return -1
	}
	return len(g.adj[v])
}

// Adj returns the vertices adjacent from vertex.
func (g *Undirected) Adj(v int) []int {
	if !g.isVertexValid(v) {
		return nil
	}
	return g.adj[v]
}

// AddEdge adds a new edge to the graph.
func (g *Undirected) AddEdge(v, w int) {
	if g.isVertexValid(v) && g.isVertexValid(w) {
		g.e++
		g.adj[v] = append(g.adj[v], w)
		g.adj[w] = append(g.adj[w], v)
	}
}

// DFS Traversal (Recursion)
func (g *Undirected) _traverseDFS(visited []bool, v int, order TraversalOrder, vertexVisitor VertexVisitor, edgeVisitor EdgeVisitor) {
	visited[v] = true

	if order == PreOrder && vertexVisitor != nil {
		if !vertexVisitor.VisitVertex(v) {
			return
		}
	}

	for _, w := range g.adj[v] {
		if !visited[w] {
			if order == PreOrder && edgeVisitor != nil {
				if !edgeVisitor.VisitEdge(v, w) {
					return
				}
			}

			g._traverseDFS(visited, w, order, vertexVisitor, edgeVisitor)
		}
	}

	if order == PostOrder && vertexVisitor != nil {
		if !vertexVisitor.VisitVertex(v) {
			return
		}
	}
}

// DFS Traversal (Driver)
func (g *Undirected) traverseDFS(s int, order TraversalOrder, vertexVisitor VertexVisitor, edgeVisitor EdgeVisitor) {
	visited := make([]bool, g.V())
	g._traverseDFS(visited, s, order, vertexVisitor, edgeVisitor)
}

// Iterative DFS Traversal
func (g *Undirected) traverseDFSi(s int, order TraversalOrder, vertexVisitor VertexVisitor, edgeVisitor EdgeVisitor) {
	visited := make([]bool, g.V())
	stack := list.NewStack(listNodeSize)

	visited[s] = true
	stack.Push(s)

	if order == PreOrder && vertexVisitor != nil {
		if !vertexVisitor.VisitVertex(s) {
			return
		}
	}

	for !stack.IsEmpty() {
		v := stack.Pop().(int)

		if order == PostOrder && vertexVisitor != nil {
			if !vertexVisitor.VisitVertex(v) {
				return
			}
		}

		for _, w := range g.adj[v] {
			if !visited[w] {
				visited[w] = true
				stack.Push(w)

				if order == PreOrder && vertexVisitor != nil {
					if !vertexVisitor.VisitVertex(w) {
						return
					}
				}

				if order == PreOrder && edgeVisitor != nil {
					if !edgeVisitor.VisitEdge(v, w) {
						return
					}
				}
			}
		}
	}
}

// BFS Traversal
func (g *Undirected) traverseBFS(s int, order TraversalOrder, vertexVisitor VertexVisitor, edgeVisitor EdgeVisitor) {
	visited := make([]bool, g.V())
	queue := list.NewQueue(listNodeSize)

	visited[s] = true
	queue.Enqueue(s)

	if order == PreOrder && vertexVisitor != nil {
		if !vertexVisitor.VisitVertex(s) {
			return
		}
	}

	for !queue.IsEmpty() {
		v := queue.Dequeue().(int)

		if order == PostOrder && vertexVisitor != nil {
			if !vertexVisitor.VisitVertex(v) {
				return
			}
		}

		for _, w := range g.adj[v] {
			if !visited[w] {
				visited[w] = true
				queue.Enqueue(w)

				if order == PreOrder && vertexVisitor != nil {
					if !vertexVisitor.VisitVertex(w) {
						return
					}
				}

				if order == PreOrder && edgeVisitor != nil {
					if !edgeVisitor.VisitEdge(v, w) {
						return
					}
				}
			}
		}
	}
}

// TraverseVertices is used for visiting all vertices in graph.
func (g *Undirected) TraverseVertices(s int, strategy TraversalStrategy, order TraversalOrder, visitor VertexVisitor) {
	if !g.isVertexValid(s) {
		return
	}

	if order != PreOrder && order != PostOrder {
		return
	}

	switch strategy {
	case DFS:
		g.traverseDFS(s, order, visitor, nil)
	case DFSi:
		g.traverseDFSi(s, order, visitor, nil)
	case BFS:
		g.traverseBFS(s, order, visitor, nil)
	}
}

// TraverseEdges is used for visiting all edges in graph.
func (g *Undirected) TraverseEdges(s int, strategy TraversalStrategy, visitor EdgeVisitor) {
	if !g.isVertexValid(s) {
		return
	}

	switch strategy {
	case DFS:
		g.traverseDFS(s, PreOrder, nil, visitor)
	case DFSi:
		g.traverseDFSi(s, PreOrder, nil, visitor)
	case BFS:
		g.traverseBFS(s, PreOrder, nil, visitor)
	}
}

// Graphviz returns a visualization of the graph in Graphviz format.
func (g *Undirected) Graphviz() string {
	graph := graphviz.NewGraph(true, false, "", "", "", graphviz.StyleSolid, graphviz.ShapeCircle)

	for i := 0; i < g.v; i++ {
		name := fmt.Sprintf("%d", i)
		graph.AddNode(graphviz.NewNode(name, "", name, "", "", "", "", ""))
	}

	for v := range g.adj {
		for _, w := range g.adj[v] {
			from := fmt.Sprintf("%d", v)
			to := fmt.Sprintf("%d", w)
			graph.AddEdge(graphviz.NewEdge(from, to, graphviz.EdgeTypeUndirected, "", "", "", "", ""))
		}
	}

	return graph.DotCode()
}
