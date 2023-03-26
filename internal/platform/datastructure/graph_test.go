package datastructure

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVertexAllEdges(t *testing.T) {
	vertex := &Vertex{
		adjacent: map[int]*Vertex{
			1: {Id: "B", adjacent: make(map[int]*Vertex), disabled: true},
			2: {Id: "C", adjacent: make(map[int]*Vertex), disabled: false},
			3: {Id: "D", adjacent: make(map[int]*Vertex), disabled: false},
			4: {Id: "E", adjacent: make(map[int]*Vertex), disabled: true},
		},
	}

	expected := []int{2, 3}
	if got := vertex.AllEdges(); !reflect.DeepEqual(got, expected) {
		t.Errorf("AllEdges() = %v, expected %v", got, expected)
	}
}

func TestVertexEnabled(t *testing.T) {
	vertex := &Vertex{disabled: false}
	if got := vertex.Enabled(); !got {
		t.Errorf("Enabled() = %v, expected true", got)
	}

	vertex.disabled = true
	if got := vertex.Enabled(); got {
		t.Errorf("Enabled() = %v, expected false", got)
	}
}

func TestVertexDisable(t *testing.T) {
	vertex := &Vertex{}
	vertex.Disable()

	if !vertex.disabled {
		t.Errorf("Disable() failed, vertex should be disabled")
	}
}

func TestVertexGetAdjacent(t *testing.T) {
	edgeId := 1
	adjacentVertex := &Vertex{}
	vertex := &Vertex{
		adjacent: map[int]*Vertex{edgeId: adjacentVertex},
	}

	if got := vertex.GetAdjacent(edgeId); got != adjacentVertex {
		t.Errorf("GetAdjacent() = %v, expected %v", got, adjacentVertex)
	}
}

func TestGraphAddVertex(t *testing.T) {
	graph := &Graph{}
	vertexId := "A"

	// Adding a new vertex
	if vertex, err := graph.AddVertex(vertexId); err != nil {
		t.Errorf("AddVertex(%q) failed with error %v", vertexId, err)
	} else {
		if vertex.Id != vertexId {
			t.Errorf("AddVertex(%q) returned vertex with Id %q, expected %q", vertexId, vertex.Id, vertexId)
		}

		if len(graph.vertices) != 1 {
			t.Errorf("AddVertex(%q) did not add vertex to graph", vertexId)
		}
	}

	// Adding a duplicate vertex
	if _, err := graph.AddVertex(vertexId); !errors.Is(err, ErrVertexDuplicated) {
		t.Errorf("AddVertex(%q) should return %v, but returned %v", vertexId, ErrVertexDuplicated, err)
	}
}

func TestGraph_GetVertex(t *testing.T) {
	graph := &Graph{}
	vertex1, _ := graph.AddVertex("1")

	// Test case: vertex found
	foundVertex := graph.GetVertex("1")
	assert.Equal(t, vertex1, foundVertex)

	// Test case: vertex not found
	foundVertex = graph.GetVertex("3")
	assert.Nil(t, foundVertex)
}

func TestGraph_AddEdge(t *testing.T) {
	graph := &Graph{}
	_, err := graph.AddVertex("1")
	assert.NoError(t, err)
	_, err = graph.AddVertex("2")
	assert.NoError(t, err)

	// Test case: successful add
	err = graph.AddEdge(1, "1", "2")
	assert.NoError(t, err)

	// Test case: duplicate edge
	err = graph.AddEdge(2, "1", "2")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEdgeDuplicated)

	// Test case: vertex not found
	err = graph.AddEdge(3, "1", "3")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrVertexNotFound)
}
