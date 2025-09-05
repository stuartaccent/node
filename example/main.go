package main

import (
	"fmt"
	"os"

	. "github.com/stuartaccent/node"
)

var (
	h1          = H1(Text("Welcome!"))
	p           = P(Text("This is a "), A(Href("/link"), Text("link")), Text(" example."))
	button      = Button(Type("button"), Class("btn"), Text("Click me"))
	inputText   = Input(ID("id-text"), Type("text"), Placeholder("Enter text..."), Required(), Value("hello"))
	inputHidden = Input(ID("id-hidden"), Type("hidden"), Value("hidden value"), Disabled())
	img         = Img(Src("/image.jpg"), Alt("Example image"))
)

func main() {
	for _, el := range []Node{
		h1,
		p,
		button,
		inputText,
		inputHidden,
		img,
	} {
		div := Div(Class("container"))
		div.Add(el)
		_ = div.Render(nil, os.Stdout)
		_, _ = fmt.Fprint(os.Stdout, "\n")
	}
}
