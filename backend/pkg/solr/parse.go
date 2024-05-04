package solr

import (
	"fmt"
	"regexp"

	"golang.org/x/text/unicode/norm"
)

var SOLR_SPECIAL_CHARACTERS = regexp.MustCompile(`(\+|\-|&&|\|\||!|\(|\)|\{|\}|\[|\]|\^|"|\~|\*|\?|:|/|AND|OR)`)

func Sanitize(s string) string {
	return SOLR_SPECIAL_CHARACTERS.ReplaceAllString(s, `\${0}`)
}

func Sanitizes(s []string) []string {
	res := make([]string, len(s))
	for i, s := range s {
		res[i] = Sanitize(s)
	}
	return res
}

type SearchWord struct {
	Word     string
	Negative bool
	Phrase   bool
}

func (w SearchWord) String() string {
	if w.Negative {
		if w.Phrase {
			return fmt.Sprintf(`-"%s"`, Sanitize(w.Word))
		} else {
			return fmt.Sprintf("-%s", Sanitize(w.Word))
		}
	} else {
		if w.Phrase {
			return fmt.Sprintf(`"%s"`, Sanitize(w.Word))
		} else {
			return Sanitize(w.Word)
		}
	}
}

const (
	READ_FIRST = 0
	COLLECTING = 1
	NORMAL     = 2
	NEGATIVE   = 4
	PHRASE     = 8
)

const (
	STATE_INITIAL             = READ_FIRST
	STATE_NORMAL              = NORMAL | COLLECTING
	STATE_NEGATIVE_READ_FIRST = NEGATIVE | READ_FIRST
	STATE_NEGATIVE_COLLECTING = NEGATIVE | COLLECTING
	STATE_PHRASE              = PHRASE | COLLECTING
	STATE_NEGATIVE_PHRASE     = NEGATIVE | PHRASE | COLLECTING
)

func Parse(s string) []SearchWord {
	res := make([]SearchWord, 0, 2)
	buffer := make([]rune, 0, 16)

	state := STATE_INITIAL
loop:
	for _, c := range norm.NFKC.String(s) {
		switch state {
		case STATE_INITIAL:
			switch c {
			case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
				continue loop
			case '-':
				state = STATE_NEGATIVE_READ_FIRST
			case '"':
				state = STATE_PHRASE
			default:
				buffer = append(buffer, c)
				state = STATE_NORMAL
			}
		case STATE_NORMAL:
			switch c {
			case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
				if w := string(buffer); w != "" {
					res = append(res, SearchWord{Word: w, Negative: (state & NEGATIVE) == NEGATIVE, Phrase: (state & PHRASE) == PHRASE})
				}
				buffer = buffer[:0]
				state = STATE_INITIAL
				continue loop
			default:
				buffer = append(buffer, c)
			}
		case STATE_NEGATIVE_READ_FIRST:
			switch c {
			case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
				if w := string(buffer); w != "" {
					res = append(res, SearchWord{Word: w, Negative: (state & NEGATIVE) == NEGATIVE, Phrase: (state & PHRASE) == PHRASE})
				}
				buffer = buffer[:0]
				state = STATE_INITIAL
				continue loop
			case '"':
				state = STATE_NEGATIVE_PHRASE
				continue loop
			default:
				buffer = append(buffer, c)
				state = STATE_NEGATIVE_COLLECTING
			}
		case STATE_NEGATIVE_COLLECTING:
			switch c {
			case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
				if w := string(buffer); w != "" {
					res = append(res, SearchWord{Word: w, Negative: (state & NEGATIVE) == NEGATIVE, Phrase: (state & PHRASE) == PHRASE})
				}
				buffer = buffer[:0]
				state = STATE_INITIAL
				continue loop
			default:
				buffer = append(buffer, c)
			}
		case STATE_PHRASE:
			switch c {
			case '"':
				if w := string(buffer); w != "" {
					res = append(res, SearchWord{Word: w, Negative: (state & NEGATIVE) == NEGATIVE, Phrase: (state & PHRASE) == PHRASE})
				}
				buffer = buffer[:0]
				state = STATE_INITIAL
				continue loop
			default:
				buffer = append(buffer, c)
			}
		case STATE_NEGATIVE_PHRASE:
			switch c {
			case '"':
				if w := string(buffer); w != "" {
					res = append(res, SearchWord{Word: w, Negative: (state & NEGATIVE) == NEGATIVE, Phrase: (state & PHRASE) == PHRASE})
				}
				buffer = buffer[:0]
				state = STATE_INITIAL
				continue loop
			default:
				buffer = append(buffer, c)
			}
		}
	}
	if len(buffer) > 0 {
		if w := string(buffer); w != "" {
			res = append(res, SearchWord{Word: w, Negative: (state & NEGATIVE) == NEGATIVE, Phrase: (state & PHRASE) == PHRASE})
		}
	}
	return res
}
