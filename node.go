// Package node contains a minimal HTML node builder and renderer.
//
// The Node type and helpers in this file allow you to declaratively construct
// small HTML fragments in Go and render them directly to an io.Writer with
// proper escaping for attribute values and text content.
package node

import (
	"context"
	"fmt"
	"html"
	"io"
)

// NodeType identifies the kind of Node.
type NodeType int

// NodeType values.
const (
	// NodeTypeTag represents an element/tag node, e.g., <div>.
	NodeTypeTag NodeType = iota
	// NodeTypeAttr represents an attribute on a tag, e.g., class="...".
	NodeTypeAttr
	// NodeTypeText represents text content within a tag.
	NodeTypeText
)

// Precomputed byte slices for faster rendering.
var (
	lt         = []byte("<")
	gt         = []byte(">")
	ltSlash    = []byte("</")
	gtSlash    = []byte("/>")
	space      = []byte(" ")
	equalQuote = []byte(`="`)
	quote      = []byte(`"`)
)

// Node represents a minimal HTML node.
//
// It can be one of three kinds (Type):
//   - NodeTypeTag:    An element/tag node like <div> or <img> with optional attributes and children.
//   - NodeTypeAttr:   An attribute node like class="..." or disabled.
//   - NodeTypeText:   Text content. Text and attribute values are HTML-escaped when rendered.
//
// For NodeTypeTag nodes:
//   - Tag is the element name (e.g., "div").
//   - Attributes holds attribute nodes (NodeTypeAttr).
//   - Children holds child nodes (tags or text).
//   - SelfClose indicates whether the tag is self-closing (e.g., <img/>).
//
// For NodeTypeAttr nodes:
//   - Key is the attribute name (e.g., "class").
//   - Value is the attribute value; if empty, a boolean attribute is emitted (e.g., disabled).
//
// For NodeTypeText nodes:
//   - Value is the text content.
type Node struct {
	Type       NodeType
	Tag        string // HTML tag name (for NodeTypeTag)
	Key        string // Attribute name (for NodeTypeAttr)
	Value      string // Attribute value or text content
	Children   []Node // Children (for NodeTypeTag)
	Attributes []Node // Attributes (for NodeTypeTag)
	SelfClose  bool   // Whether Tag is self-closing
}

// Add appends the provided nodes to n.
// Attribute nodes are added to Attributes, everything else to Children.
func (n *Node) Add(nodes ...Node) {
	children, attrs := separateChildrenAndAttrs(nodes)
	if len(children) > 0 {
		n.Children = append(n.Children, children...)
	}
	if len(attrs) > 0 {
		n.Attributes = append(n.Attributes, attrs...)
	}
}

// Render writes the HTML representation of the node and its descendants to w.
// Text and attribute values are HTML-escaped.
func (n *Node) Render(ctx context.Context, w io.Writer) error {
	switch n.Type {
	case NodeTypeTag:
		return n.renderTag(ctx, w)
	case NodeTypeAttr:
		return n.renderAttr(ctx, w)
	case NodeTypeText:
		return n.renderText(ctx, w)
	default:
		return fmt.Errorf("unknown node type: %d", n.Type)
	}
}

// renderTag renders a NodeTypeTag node: <tag [attrs]>[children]</tag> or self-closing.
func (n *Node) renderTag(ctx context.Context, w io.Writer) error {
	if _, err := w.Write(lt); err != nil {
		return err
	}
	if _, err := io.WriteString(w, n.Tag); err != nil {
		return err
	}
	for i := range n.Attributes {
		if _, err := w.Write(space); err != nil {
			return err
		}
		if err := n.Attributes[i].Render(ctx, w); err != nil {
			return err
		}
	}
	if n.SelfClose {
		if _, err := w.Write(gtSlash); err != nil {
			return err
		}
		return nil
	}
	if _, err := w.Write(gt); err != nil {
		return err
	}
	for i := range n.Children {
		if err := n.Children[i].Render(ctx, w); err != nil {
			return err
		}
	}
	if _, err := w.Write(ltSlash); err != nil {
		return err
	}
	if _, err := io.WriteString(w, n.Tag); err != nil {
		return err
	}
	_, err := w.Write(gt)
	return err
}

// renderAttr renders a NodeTypeAttr node: key[="value"]. Values are escaped; empty value emits a boolean attribute.
func (n *Node) renderAttr(_ context.Context, w io.Writer) error {
	if _, err := io.WriteString(w, n.Key); err != nil {
		return err
	}
	if n.Value != "" {
		if _, err := w.Write(equalQuote); err != nil {
			return err
		}
		if _, err := io.WriteString(w, html.EscapeString(n.Value)); err != nil {
			return err
		}
		if _, err := w.Write(quote); err != nil {
			return err
		}
	}
	return nil
}

// renderText renders a NodeTypeText node, escaping HTML entities.
func (n *Node) renderText(_ context.Context, w io.Writer) error {
	_, err := io.WriteString(w, html.EscapeString(n.Value))
	return err
}

// separateChildrenAndAttrs splits a mixed list of nodes into children (non-attr) and attrs.
// Helper used by element constructors and Add.
func separateChildrenAndAttrs(nodes []Node) (children []Node, attrs []Node) {
	for _, node := range nodes {
		if node.Type == NodeTypeAttr {
			attrs = append(attrs, node)
		} else {
			children = append(children, node)
		}
	}
	return children, attrs
}
