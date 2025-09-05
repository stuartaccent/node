package node

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"testing"
)

// helperRender renders n to a bytes.Buffer and resets it each iteration.
func helperRender(b *testing.B, n Node) {
	b.ReportAllocs()
	var buf bytes.Buffer
	// Optionally grow to reduce resizes; tune per benchmark if desired.
	buf.Grow(1 << 15)

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := n.Render(ctx, &buf); err != nil {
			b.Fatalf("render error: %v", err)
		}
	}
}

func BenchmarkRenderSmall(b *testing.B) {
	// A small but representative tree (similar to main.go usage)
	h1 := H1(Text("Welcome!"))
	p := P(Text("This is a "), A(Href("/link"), Text("link")), Text(" example."))
	button := Button(Type("button"), Class("btn"), Text("Click me"))

	div := Div(Class("container"))
	div.Add(h1, p, button)

	helperRender(b, div)
}

func BenchmarkRenderBalanced_b3_d4(b *testing.B) {
	// Balanced tree with branching factor 3 and depth 4
	var build func(depth int) Node
	build = func(depth int) Node {
		if depth == 0 {
			return Span(Text("leaf"))
		}
		return Div(
			Class("lvl-"+strconv.Itoa(depth)),
			build(depth-1), build(depth-1), build(depth-1),
		)
	}
	root := build(4)
	helperRender(b, root)
}

func BenchmarkRenderDeep_d500(b *testing.B) {
	// 500 nested divs with a single text at the bottom
	n := Text("end")
	node := Div(n)
	for i := 0; i < 500; i++ {
		node = Div(node)
	}
	helperRender(b, node)
}

func BenchmarkRenderWide_attrs200_texts200(b *testing.B) {
	// One node with many attributes and many text children
	attrs := make([]Node, 0, 200)
	for i := 0; i < 200; i++ {
		attrs = append(attrs, Attr("data-k"+strconv.Itoa(i), "v"+strconv.Itoa(i)))
	}

	texts := make([]Node, 0, 200)
	for i := 0; i < 200; i++ {
		texts = append(texts, Text("t"+strconv.Itoa(i)))
	}

	all := append([]Node{}, attrs...)
	all = append(all, texts...)
	node := Div(all...)

	helperRender(b, node)
}

func BenchmarkRenderEscapeHeavyText(b *testing.B) {
	// A sizable text with many escapable characters
	base := "<>&\"'"
	payload := strings.Repeat(base, 4096/len(base)) // ~4 KB
	node := Text(payload)
	helperRender(b, node)
}
