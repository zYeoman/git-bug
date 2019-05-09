package mardown

import "strconv"

type numbering struct {
	levels [6]int
}

func (n *numbering) NextLevel(level int) {
	if level <= 0 {
		panic("level start at 1, ask blackfriday why")
	}
	if level > 6 {
		panic("Markdown is limited to 6 levels of heading")
	}

	n.levels[level-1]++
	for i := level; i < 6; i++ {
		n.levels[i] = 0
	}
}

func (n *numbering) Render() string {
	slice := n.levels[:]

	// pop the last zero levels
	for i := 5; i >= 0; i-- {
		if n.levels[i] != 0 {
			break
		}
		slice = slice[:len(slice)-1]
	}

	var result string

	for i := range slice {
		if i > 0 {
			result += "."
		}
		result += strconv.Itoa(slice[i])
	}

	return result
}
