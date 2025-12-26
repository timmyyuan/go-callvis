package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestDotAttrsFormatting(t *testing.T) {
	attrs := dotAttrs{"color": "red"}

	if got := attrs.String(); got != `color="red"` {
		t.Fatalf("dotAttrs.String() = %q, want %q", got, `color="red"`)
	}

	if got := attrs.Lines(); got != "color=\"red\";" {
		t.Fatalf("dotAttrs.Lines() = %q, want %q", got, `color="red";`)
	}
}

func TestDotGraphWriteDot(t *testing.T) {
	root := NewDotCluster("root")
	nodeA := &dotNode{ID: "A", Attrs: dotAttrs{"label": "FuncA"}}
	nodeB := &dotNode{ID: "B", Attrs: dotAttrs{"label": "FuncB"}}
	root.Nodes = append(root.Nodes, nodeA, nodeB)

	graph := &dotGraph{
		Title:   "Test Graph",
		Cluster: root,
		Edges: []*dotEdge{
			{From: nodeA, To: nodeB, Attrs: dotAttrs{"label": "calls"}},
		},
		Options: map[string]string{
			"rankdir":   "LR",
			"nodesep":   "0.5",
			"nodeshape": "box",
			"nodestyle": "filled",
			"minlen":    "1",
		},
	}

	var buf bytes.Buffer
	if err := graph.WriteDot(&buf); err != nil {
		t.Fatalf("WriteDot returned error: %v", err)
	}
	out := buf.String()

	for _, substr := range []string{
		`digraph gocallvis {`,
		`label="Test Graph";`,
		`rankdir="LR";`,
		`subgraph "cluster_root" {`,
		`"A" [ label="FuncA" ]`,
		`"B" [ label="FuncB" ]`,
		`"A" -> "B" [ label="calls" ]`,
	} {
		if !strings.Contains(out, substr) {
			t.Fatalf("DOT output missing %q\noutput:\n%s", substr, out)
		}
	}
}
