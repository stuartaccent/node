package node

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
