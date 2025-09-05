# node

A tiny, allocation-conscious HTML node builder and renderer for Go. It provides a small set of primitives (Node, elements, attributes, and text) and renders them to io.Writer without reflection. The design aims for clarity and speed, with simple, chainable helpers.

## Features

- Lightweight tree of nodes (elements, attributes, text)
- Composable helpers like Div(...), Span(...), P(...), A(...), Text("..."), Attr(key, val)
- Context-aware rendering to io.Writer
- Benchmarks included (see node_bench_test.go)

## Install

```bash
go get ./...
```

If this module is used externally, you can use:

```bash
go get github.com/stuartaccent/node
```

Adjust the module path to your actual repository if publishing.

## Quick Start

Below is a minimal example showing how to build and render a simple HTML fragment using the helpers. This mirrors the usage in main.go.

```go
package main

import (
	"bytes"
	"context"
	"fmt"

	. "github.com/stuartaccent/node"
)

func main() {
	// Build a small tree
	h1 := H1(Text("Welcome!"))
	p := P(Text("This is a "), A(Href("/link"), Text("link")), Text(" example."))
	button := Button(Type("button"), Class("btn"), Text("Click me"))

	div := Div(Class("container"))
	div.Add(h1, p, button)

	// Render to a buffer
	var buf bytes.Buffer
	if err := div.Render(context.Background(), &buf); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}
```

Typical output (minified style):

```html
<div class="container"><h1>Welcome!</h1><p>This is a <a href="/link">link</a> example.</p><button type="button" class="btn">Click me</button></div>
```

## API Sketch

- Node: interface implemented by element, attribute, and text nodes
- Elements: Div(...), Span(...), P(...), A(...), H1(...), Button(...), etc.
- Text: Text("...") escapes HTML as needed
- Attributes: Attr(key, val) or specific helpers like Class("..."), Href("..."), Type("...")
- Composition:
  - Pass children to element helpers: Div(child1, child2)
  - Add more later via `Add(children...)`
- Rendering: `node.Render(ctx, w io.Writer) error`

Refer to node.go for the full set of helpers and details.

## Benchmarks

The repository includes a set of benchmarks in node_bench_test.go. To run them:

```bash
go test -bench=. -benchmem
```

You can focus on specific cases, e.g.:

```bash
go test -bench=RenderSmall -benchmem
```

## Development

- Go 1.21+ recommended (match your local go.mod)
- Format and vet code:

```bash
go fmt ./...
go vet ./...
```

## License

MIT.
