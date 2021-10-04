package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents a link (<a href="...">) in an HTML document.
type Link struct {
	Href string
	Text string
}

// Parse will take in an HTML document and will return a slice
// of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	// 1.Find <a> nodes in document
	nodes := findNodes(doc)
	var links []Link
	// 2. for each link node...
	for _, node := range nodes {
		// 2.1 build a Link
		links = append(links, buildLink(node))
	}
	// 3. return the Links
	return links, nil
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = getText(n)
	return ret
}

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += getText(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}

func findNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, findNodes(c)...)
	}
	return ret
}

//func depthFirstSearch(n *html.Node, padding string) {
//	tag := n.Data
//	if n.Type == html.ElementNode {
//		tag = "<" + tag + ">"
//	}
//	fmt.Println(padding, tag)
//	for c := n.FirstChild; c != nil; c = c.NextSibling {
//		depthFirstSearch(c, padding + "  ")
//	}
//}
