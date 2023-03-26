package datastructure

import (
	"errors"
	"fmt"
)

// Graph is a simple data structure implementation without any special considerations.
// Each operation that modifies the graph does integrity checks.
type Graph struct {
	vertices []*Vertex
}

// Vertex ...
type Vertex struct {
	Id       string
	adjacent map[int]*Vertex

	// Having a disabled flag is more performant than actually removing the item
	disabled bool
}

// AllEdges returns all enabled edge Ids
func (vertex *Vertex) AllEdges() []int {
	edges := make([]int, 0)
	for edgeId, edge := range vertex.adjacent {
		if !edge.disabled {
			edges = append(edges, edgeId)
		}
	}

	return edges
}

func (vertex *Vertex) Enabled() bool {
	return !vertex.disabled
}

func (vertex *Vertex) Disable() {
	vertex.disabled = true
}

func (vertex *Vertex) GetAdjacent(edgeId int) *Vertex {
	return vertex.adjacent[edgeId]
}

var (
	ErrVertexDuplicated = errors.New("vertex already exist")
	ErrEdgeDuplicated   = errors.New("edge already exist")
	ErrVertexNotFound   = errors.New("vertex does not exist")
)

func (graph *Graph) AddVertex(id string) (*Vertex, error) {
	// Integrity check
	for _, vertex := range graph.vertices {
		if vertex.Id == id {
			return nil, fmt.Errorf("%w: %q", ErrVertexDuplicated, id)
		}
	}

	newVertex := &Vertex{
		Id:       id,
		adjacent: make(map[int]*Vertex),
	}

	graph.vertices = append(graph.vertices, newVertex)

	return newVertex, nil
}

// GetVertex returns the vertex with the matching ID, if it's not found, it returns nil.
func (graph *Graph) GetVertex(id string) *Vertex {
	for _, vertex := range graph.vertices {
		if vertex.Id == id {
			return vertex
		}
	}

	return nil
}

// AddEdge returns an error in case an edge with the same destination already exists.
func (graph *Graph) AddEdge(id int, from, to string) error {
	var (
		fromVertex = graph.GetVertex(from)
		toVertex   = graph.GetVertex(to)
	)

	// Integrity checks
	{
		if fromVertex == nil {
			return fmt.Errorf("%w: %q", ErrVertexNotFound, from)
		}

		if toVertex == nil {
			return fmt.Errorf("%w: %q", ErrVertexNotFound, to)
		}

		for _, fromVertexAdjacent := range fromVertex.adjacent {
			if fromVertexAdjacent.Id == to {
				return fmt.Errorf("%w: from %q to %q", ErrEdgeDuplicated, from, to)
			}
		}
	}

	fromVertex.adjacent[id] = toVertex
	return nil
}
