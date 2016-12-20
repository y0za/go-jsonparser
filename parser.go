package jsonparser

import (
	"errors"
	"fmt"
	"strconv"
)

// Parser is struct for parsing json
type Parser struct {
	at   int    // index of current character
	cc   rune   // current character
	text []rune // whole text
}

var (
	escapee = map[rune]rune{
		'"':  '"',
		'\\': '\\',
		'/':  '/',
		'b':  'b',
		'f':  '\f',
		'n':  '\n',
		'r':  '\r',
		't':  '\t',
	}
)

func (p *Parser) next(c rune) (rune, error) {
	if c != 0 && c != p.cc {
		return 0, fmt.Errorf("Expected %q instead of %q", c, p.cc)
	}

	p.cc = p.text[p.at]
	p.at++
	return p.cc, nil
}

func (p *Parser) number() (float64, error) {
	var str []rune

	if p.cc == '-' {
		str = append(str, '-')
		p.next('-')
	}
	for p.cc >= '0' && p.cc <= '9' {
		str = append(str, p.cc)
		p.next(0)
	}
	if p.cc == '.' {
		str = append(str, '.')
		p.next(0)
		for p.cc >= '0' && p.cc <= '9' {
			str = append(str, p.cc)
			p.next(0)
		}
	}
	if p.cc == 'e' || p.cc == 'E' {
		str = append(str, p.cc)
		p.next(0)
		if p.cc == '-' || p.cc == '+' {
			str = append(str, p.cc)
			p.next(0)
		}
		for p.cc >= '0' && p.cc <= '9' {
			str = append(str, p.cc)
			p.next(0)
		}
	}

	return strconv.ParseFloat(string(str), 64)
}

func (p *Parser) string() (string, error) {
	var str []rune

	if p.cc != '"' {
		return "", errors.New("Invalid string")
	}

	p.next(0)
	for p.cc != 0 {
		if p.cc == '"' {
			return string(str), nil
		} else if p.cc == '\\' {
			p.next(0)
			if p.cc == 'u' {
				var uffff int64
				for i := 0; i < 4; i++ {
					p.next(0)
					hex, err := strconv.ParseInt(string(p.cc), 16, 64)
					if err != nil {
						break
					}
					uffff += uffff*16 + hex
				}
				str = append(str, rune(uffff))
			} else if e, ok := escapee[p.cc]; ok {
				str = append(str, e)
			}
		} else {
			str = append(str, p.cc)
		}
		p.next(0)
	}

	return "", errors.New("Invalid string")
}
