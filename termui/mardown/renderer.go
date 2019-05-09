package mardown

import (
	"fmt"
	"io"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-runewidth"
	"gopkg.in/russross/blackfriday.v2"

	"github.com/MichaelMure/git-bug/util/text"
)

var _ blackfriday.Renderer = &renderer{}

type renderer struct {
	// maximum line width allowed
	lineWidth int
	// constant left padding to apply
	leftPad int

	// Count the number of line in the rendered output
	lines int

	numbering numbering

	paragraph strings.Builder
}

func newRenderer(lineWidth int, leftPad int) *renderer {
	return &renderer{lineWidth: lineWidth, leftPad: leftPad}
}

func (r *renderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	fmt.Println(node, entering)

	pad := strings.Repeat(" ", r.leftPad)

	switch node.Type {
	case blackfriday.Document:
		// Nothing to do

	case blackfriday.BlockQuote:

	case blackfriday.List:

	case blackfriday.Item:

	case blackfriday.Paragraph:
		if entering {
			r.paragraph.Reset()
		} else {
			out, _ := text.WrapLeftPadded(r.paragraph.String(), r.lineWidth, r.leftPad)
			_, _ = fmt.Fprint(w, out, "\n\n")
		}

	case blackfriday.Heading:
		// the child node of a heading is a blackfriday.Text. We render the whole thing
		// in one go and skip the child.

		// render the full line with the numbering
		r.numbering.NextLevel(node.Level)
		rendered := fmt.Sprintf("%s%s %s", pad, r.numbering.Render(), string(node.FirstChild.Literal))

		// output the text, truncated if needed, no line break
		truncated := runewidth.Truncate(rendered, r.lineWidth, "…")
		colored := aurora.Colorize(truncated, shade(node.Level)).String()
		_, _ = fmt.Fprintln(w, colored)

		// render the underline, if any
		if node.Level == 1 {
			_, _ = fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("─", r.lineWidth-r.leftPad))
		}

		_, _ = fmt.Fprintln(w)

		return blackfriday.SkipChildren

	case blackfriday.HorizontalRule:

	case blackfriday.Emph:
		r.paragraph.WriteString(aurora.Italic(string(node.FirstChild.Literal)).String())
		return blackfriday.SkipChildren

	case blackfriday.Strong:
		r.paragraph.WriteString(aurora.Bold(string(node.FirstChild.Literal)).String())
		return blackfriday.SkipChildren

	case blackfriday.Del:
		r.paragraph.WriteString(aurora.CrossedOut(string(node.FirstChild.Literal)).String())
		return blackfriday.SkipChildren

	case blackfriday.Link:

	case blackfriday.Image:

	case blackfriday.Text:
		r.paragraph.Write(node.Literal)

	case blackfriday.HTMLBlock:

	case blackfriday.CodeBlock:

	case blackfriday.Softbreak:

	case blackfriday.Hardbreak:

	case blackfriday.Code:

	case blackfriday.HTMLSpan:

	case blackfriday.Table:

	case blackfriday.TableCell:

	case blackfriday.TableHead:

	case blackfriday.TableBody:

	case blackfriday.TableRow:

	default:
		panic("Unknown node type " + node.Type.String())
	}

	return blackfriday.GoToNext
}

func (*renderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {}

func (*renderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {}
