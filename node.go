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
func (n *Node) renderAttr(ctx context.Context, w io.Writer) error {
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
func (n *Node) renderText(ctx context.Context, w io.Writer) error {
	_, err := io.WriteString(w, html.EscapeString(n.Value))
	return err
}

// Div creates a <div> element with the given children and attributes.
// Attribute nodes in the variadic list are separated automatically.
func Div(nodes ...Node) Node {
	children, attrs := separateChildrenAndAttrs(nodes)
	return Node{
		Type:       NodeTypeTag,
		Tag:        "div",
		Children:   children,
		Attributes: attrs,
	}
}

// Span creates a <span> element with the given children and attributes.
func Span(nodes ...Node) Node {
	children, attrs := separateChildrenAndAttrs(nodes)
	return Node{
		Type:       NodeTypeTag,
		Tag:        "span",
		Children:   children,
		Attributes: attrs,
	}
}

// P creates a <p> element with the given children and attributes.
func P(nodes ...Node) Node {
	children, attrs := separateChildrenAndAttrs(nodes)
	return Node{
		Type:       NodeTypeTag,
		Tag:        "p",
		Children:   children,
		Attributes: attrs,
	}
}

// H1 creates an <h1> element with the given children and attributes.
func H1(nodes ...Node) Node {
	children, attrs := separateChildrenAndAttrs(nodes)
	return Node{
		Type:       NodeTypeTag,
		Tag:        "h1",
		Children:   children,
		Attributes: attrs,
	}
}

// A creates an <a> anchor element with the given children and attributes.
// Use Href(...) to set the link destination.
func A(nodes ...Node) Node {
	children, attrs := separateChildrenAndAttrs(nodes)
	return Node{
		Type:       NodeTypeTag,
		Tag:        "a",
		Children:   children,
		Attributes: attrs,
	}
}

// Button creates a <button> element with the given children and attributes.
func Button(nodes ...Node) Node {
	children, attrs := separateChildrenAndAttrs(nodes)
	return Node{
		Type:       NodeTypeTag,
		Tag:        "button",
		Children:   children,
		Attributes: attrs,
	}
}

// Input creates an <input/> element. It is rendered as self-closing.
func Input(nodes ...Node) Node {
	children, attrs := separateChildrenAndAttrs(nodes)
	return Node{
		Type:       NodeTypeTag,
		Tag:        "input",
		Children:   children,
		Attributes: attrs,
		SelfClose:  true,
	}
}

// Img creates an <img/> element. It is rendered as self-closing.
func Img(nodes ...Node) Node {
	children, attrs := separateChildrenAndAttrs(nodes)
	return Node{
		Type:       NodeTypeTag,
		Tag:        "img",
		Children:   children,
		Attributes: attrs,
		SelfClose:  true,
	}
}

// Class sets the class attribute: class="value".
func Class(value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   "class",
		Value: value,
	}
}

// ID sets the id attribute: id="value".
func ID(value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   "id",
		Value: value,
	}
}

// Href sets the href attribute on anchors: href="value".
func Href(value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   "href",
		Value: value,
	}
}

// Src sets the src attribute: src="value".
func Src(value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   "src",
		Value: value,
	}
}

// Alt sets the alt attribute: alt="value".
func Alt(value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   "alt",
		Value: value,
	}
}

// Type sets the type attribute: type="value" (e.g., for inputs and buttons).
func Type(value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   "type",
		Value: value,
	}
}

// Value sets the value attribute: value="value".
func Value(value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   "value",
		Value: value,
	}
}

// Placeholder sets the placeholder attribute: placeholder="value".
func Placeholder(value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   "placeholder",
		Value: value,
	}
}

// Disabled sets the boolean disabled attribute.
func Disabled() Node {
	return Node{
		Type: NodeTypeAttr,
		Key:  "disabled",
	}
}

// Required sets the boolean required attribute.
func Required() Node {
	return Node{
		Type: NodeTypeAttr,
		Key:  "required",
	}
}

// Attr creates a generic attribute node: key="value". If value is empty, a boolean attribute is emitted.
func Attr(key, value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   key,
		Value: value,
	}
}

// Text creates a text node. Content is HTML-escaped during rendering.
func Text(content string) Node {
	return Node{
		Type:  NodeTypeText,
		Value: content,
	}
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
