package mardown

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"gopkg.in/russross/blackfriday.v2"

	"github.com/MichaelMure/git-bug/util/colors"
	"github.com/MichaelMure/git-bug/util/text"
)

var _ blackfriday.Renderer = &renderer{}

type renderer struct {
	// maximum line width allowed
	lineWidth int
	// constant left padding to apply
	leftPad int

	// Count the number of line in the rendered output
	// lines int

	headingNumbering headingNumbering

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
		if entering {
			_, _ = fmt.Fprint(w, pad, "  • ")
		}

	case blackfriday.Paragraph:
		if !entering {
			out, _ := text.WrapLeftPadded(r.paragraph.String(), r.lineWidth, r.leftPad)
			_, _ = fmt.Fprint(w, out, "\n")

			if node.Next != nil {
				switch node.Next.Type {
				case blackfriday.Paragraph, blackfriday.Heading, blackfriday.HorizontalRule:
					_, _ = fmt.Fprintln(w)
				}
			}
			r.paragraph.Reset()
		}

	case blackfriday.Heading:
		// the child node of a heading is a blackfriday.Text. We render the whole thing
		// in one go and skip the child.

		// render the full line with the headingNumbering
		r.headingNumbering.Observe(node.Level)
		rendered := fmt.Sprintf("%s%s %s", pad, r.headingNumbering.Render(), string(node.FirstChild.Literal))

		// output the text, truncated if needed, no line break
		truncated := runewidth.Truncate(rendered, r.lineWidth, "…")
		colored := shade(node.Level)(truncated)
		_, _ = fmt.Fprintln(w, colored)

		// render the underline, if any
		if node.Level == 1 {
			_, _ = fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("─", r.lineWidth-r.leftPad))
		}

		_, _ = fmt.Fprintln(w)

		return blackfriday.SkipChildren

	case blackfriday.HorizontalRule:
		_, _ = fmt.Fprintf(w, "%s%s\n\n", pad, strings.Repeat("─", r.lineWidth-r.leftPad))

	case blackfriday.Emph:
		r.paragraph.WriteString(colors.Italic(string(node.FirstChild.Literal)))
		return blackfriday.SkipChildren

	case blackfriday.Strong:
		r.paragraph.WriteString(colors.Bold(string(node.FirstChild.Literal)))
		return blackfriday.SkipChildren

	case blackfriday.Del:
		r.paragraph.WriteString(colors.CrossedOut(string(node.FirstChild.Literal)))
		return blackfriday.SkipChildren

	case blackfriday.Link:

	case blackfriday.Image:

	case blackfriday.Text:
		r.paragraph.Write(node.Literal)

	case blackfriday.HTMLBlock:

	case blackfriday.CodeBlock:
		r.RenderCodeBlock(w, node)

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

func (r *renderer) RenderCodeBlock(w io.Writer, node *blackfriday.Node) {
	code := string(node.Literal)
	var lexer chroma.Lexer
	// try to get the lexer from the language tag if any
	if len(node.CodeBlockData.Info) > 0 {
		lexer = lexers.Get(string(node.CodeBlockData.Info))
	}
	// fallback on detection
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}
	// all failed :-(
	if lexer == nil {
		lexer = lexers.Fallback
	}
	// simplify the lexer output
	lexer = chroma.Coalesce(lexer)

	var formatter chroma.Formatter
	if color.NoColor {
		formatter = formatters.Fallback
	} else {
		formatter = formatters.TTY8
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		// Something failed, falling back to no highlight render
		r.RenderFormattedCodeBlock(w, code)
		return
	}

	buf := &bytes.Buffer{}

	// TODO: a custom style for the terminal would probably be better
	err = formatter.Format(buf, styles.Pygments, iterator)
	if err != nil {
		// Something failed, falling back to no highlight render
		r.RenderFormattedCodeBlock(w, code)
		return
	}

	r.RenderFormattedCodeBlock(w, buf.String())
}

func (r *renderer) RenderFormattedCodeBlock(w io.Writer, code string) {
	fmt.Println(code)

	// remove the trailing line break
	code = strings.TrimRight(code, "\n")

	pad := strings.Repeat(" ", r.leftPad) + colors.GreenBold("｜ ")

	output, _ := text.WrapWithPad(code, r.lineWidth, pad)
	_, _ = fmt.Fprint(w, output)

	_, _ = fmt.Fprintf(w, "\n\n")
}
