package node

// Attr creates a generic attribute node: key="value". If value is empty, a boolean attribute is emitted.
func Attr(key, value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   key,
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

// Class sets the class attribute: class="value".
func Class(value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   "class",
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

// Href sets the href attribute on anchors: href="value".
func Href(value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   "href",
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

// Placeholder sets the placeholder attribute: placeholder="value".
func Placeholder(value string) Node {
	return Node{
		Type:  NodeTypeAttr,
		Key:   "placeholder",
		Value: value,
	}
}

// Required sets the boolean required attribute.
func Required() Node {
	return Node{
		Type: NodeTypeAttr,
		Key:  "required",
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

// Text creates a text node. Content is HTML-escaped during rendering.
func Text(content string) Node {
	return Node{
		Type:  NodeTypeText,
		Value: content,
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
