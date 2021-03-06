package asm

// An fstring is a string that keeps track of its position within the
// file from which it was read.
type fstring struct {
	row    int    // line number of string being worked on
	column int    // starting column of substring currently being worked on
	str    string // the actual substring of interest
	full   string // the full line as originally read from the file
}

func newFstring(row int, str string) fstring {
	return fstring{row, 0, str, str}
}

func (l fstring) String() string {
	return l.str
}

func (l fstring) consume(n int) fstring {
	return fstring{l.row, l.column + n, l.str[n:], l.full}
}

func (l fstring) trunc(n int) fstring {
	return fstring{l.row, l.column, l.str[:n], l.full}
}

func (l fstring) substr(start, stop int) fstring {
	return fstring{l.row, l.column + start, l.str[start:stop], l.full}
}

func (l fstring) isEmpty() bool {
	return len(l.str) == 0
}

func (l fstring) startsWith(fn func(c byte) bool) bool {
	return len(l.str) > 0 && fn(l.str[0])
}

func (l fstring) startsWithChar(c byte) bool {
	return len(l.str) > 0 && l.str[0] == c
}

func (l fstring) startsWithString(s string) bool {
	return len(l.str) >= len(s) && l.str[:len(s)] == s
}

func (l fstring) consumeWhitespace() fstring {
	return l.consume(l.scanWhile(whitespace))
}

func (l fstring) scanWhile(fn func(c byte) bool) int {
	i := 0
	for ; i < len(l.str) && fn(l.str[i]); i++ {
	}
	return i
}

func (l fstring) scanUntil(fn func(c byte) bool) int {
	i := 0
	for ; i < len(l.str) && !fn(l.str[i]); i++ {
	}
	return i
}

func (l fstring) consumeWhile(fn func(c byte) bool) (consumed, remain fstring) {
	i := l.scanWhile(fn)
	consumed, remain = l.trunc(i), l.consume(i)
	return
}

func (l fstring) consumeUntil(fn func(c byte) bool) (consumed, remain fstring) {
	i := l.scanUntil(fn)
	consumed, remain = l.trunc(i), l.consume(i)
	return
}

func (l fstring) peekNextWord() fstring {
	return l.trunc(l.scanWhile(wordChar))
}

func (l fstring) stripTrailingComment() fstring {
	lastNonWs := 0
	for i := 0; i < len(l.str); i++ {
		if l.str[i] == ';' {
			break
		}
		if l.str[i] != ' ' && l.str[i] != '\t' {
			lastNonWs = i + 1
		}
	}
	return l.trunc(lastNonWs)
}

//
// character helper functions
//

func whitespace(c byte) bool {
	return c == ' ' || c == '\t'
}

func wordChar(c byte) bool {
	return c != ' ' && c != '\t'
}

func alpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func decimal(c byte) bool {
	return (c >= '0' && c <= '9')
}

func hexadecimal(c byte) bool {
	return decimal(c) || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')
}

func binary(c byte) bool {
	return c == '0' || c == '1'
}

func opcodeChar(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func labelStartChar(c byte) bool {
	return alpha(c) || c == '_' || c == '.'
}

func labelChar(c byte) bool {
	return alpha(c) || decimal(c) || c == '_' || c == '.'
}

func identifierStartChar(c byte) bool {
	return alpha(c) || c == '_'
}

func identifierChar(c byte) bool {
	return alpha(c) || decimal(c) || c == '_' || c == '.'
}
