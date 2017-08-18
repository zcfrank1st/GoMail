package prompt

import (
	"log"
	"strings"
)

const (
	shortenSuffix    = "..."
	leftPrefix       = " "
	leftSuffix       = " "
	rightPrefix      = " "
	rightSuffix      = " "
	leftMargin       = len(leftPrefix + leftSuffix)
	rightMargin      = len(rightPrefix + rightSuffix)
	completionMargin = leftMargin + rightMargin
)

type Suggest struct {
	Text        string
	Description string
}

type CompletionManager struct {
	selected  int // -1 means nothing one is selected.
	tmp       []Suggest
	max       uint16
	completer Completer
}

func (c *CompletionManager) GetSelectedSuggestion() (s Suggest, ok bool) {
	if c.selected == -1 {
		return Suggest{}, false
	} else if c.selected < -1 {
		log.Printf("[ERROR] shoud be reached here, selected=%d", c.selected)
		c.selected = -1
		return Suggest{}, false
	}
	return c.tmp[c.selected], true
}

func (c *CompletionManager) GetSuggestions() []Suggest {
	return c.tmp
}

func (c *CompletionManager) Reset() {
	c.selected = -1
	c.Update(*NewDocument())
	return
}

func (c *CompletionManager) Update(in Document) {
	c.tmp = c.completer(in)
	return
}

func (c *CompletionManager) Previous() {
	c.selected--
	c.update()
	return
}

func (c *CompletionManager) Next() {
	c.selected++
	c.update()
	return
}

func (c *CompletionManager) Completing() bool {
	return c.selected != -1
}

func (c *CompletionManager) update() {
	max := int(c.max)
	if len(c.tmp) < max {
		max = len(c.tmp)
	}
	if c.selected >= max {
		c.Reset()
	} else if c.selected < -1 {
		c.selected = max - 1
	}
}

func formatTexts(o []string, max int, prefix, suffix string) (new []string, width int) {
	l := len(o)
	n := make([]string, l)

	lenPrefix := len([]rune(prefix))
	lenSuffix := len([]rune(suffix))
	lenShorten := len(shortenSuffix)
	min := lenPrefix + lenSuffix + lenShorten
	for i := 0; i < l; i++ {
		if width < len([]rune(o[i])) {
			width = len([]rune(o[i]))
		}
	}

	if width == 0 {
		return n, 0
	}
	if min >= max {
		log.Println("[WARN] formatTexts: max is lower than length of prefix and suffix.")
		return n, 0
	}
	if lenPrefix+width+lenSuffix > max {
		width = max - lenPrefix - lenSuffix
	}

	for i := 0; i < l; i++ {
		r := []rune(o[i])
		x := len(r)
		if x <= width {
			spaces := strings.Repeat(" ", width-x)
			n[i] = prefix + o[i] + spaces + suffix
		} else if x > width {
			n[i] = prefix + string(r[:width-lenShorten]) + shortenSuffix + suffix
		}
	}
	return n, lenPrefix + width + lenSuffix
}

func formatSuggestions(suggests []Suggest, max int) (new []Suggest, width int) {
	num := len(suggests)
	new = make([]Suggest, num)

	left := make([]string, num)
	for i := 0; i < num; i++ {
		left[i] = suggests[i].Text
	}
	right := make([]string, num)
	for i := 0; i < num; i++ {
		right[i] = suggests[i].Description
	}

	left, leftWidth := formatTexts(left, max, leftPrefix, leftSuffix)
	if leftWidth == 0 {
		return []Suggest{}, 0
	}
	right, rightWidth := formatTexts(right, max-leftWidth, rightPrefix, rightSuffix)

	for i := 0; i < num; i++ {
		new[i] = Suggest{Text: left[i], Description: right[i]}
	}
	return new, leftWidth + rightWidth
}

func NewCompletionManager(completer Completer, max uint16) *CompletionManager {
	return &CompletionManager{
		selected:  -1,
		max:       max,
		completer: completer,
	}
}
