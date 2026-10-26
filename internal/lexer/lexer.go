package lexer

import (
	"fmt"
	"strings"
)

type Lexer struct {
	input  []rune
	pos    int
	line   int
	col    int
	tokens []Token
}

func New(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		pos:   0,
		line:  1,
		col:   1,
	}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekAt(offset int) rune {
	if l.pos+offset >= len(l.input) {
		return 0
	}
	return l.input[l.pos+offset]
}

func (l *Lexer) advance() rune {
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && (l.peek() == ' ' || l.peek() == '\t' || l.peek() == '\r') {
		l.advance()
	}
}

func (l *Lexer) skipComment() {
	if l.peek() == '/' && l.peekAt(1) == '/' {
		for l.pos < len(l.input) && l.peek() != '\n' {
			l.advance()
		}
	} else if l.peek() == '/' && l.peekAt(1) == '*' {
		l.advance()
		l.advance()
		for l.pos < len(l.input)-1 {
			if l.peek() == '*' && l.peekAt(1) == '/' {
				l.advance()
				l.advance()
				return
			}
			l.advance()
		}
	}
}

func (l *Lexer) readString() Token {
	line, col := l.line, l.col
	quote := l.advance()
	var sb strings.Builder
	for l.pos < len(l.input) && l.peek() != quote {
		ch := l.advance()
		if ch == '\\' && l.pos < len(l.input) {
			esc := l.advance()
			switch esc {
			case 'n':
				sb.WriteRune('\n')
			case 't':
				sb.WriteRune('\t')
			case '"', '\'', '\\':
				sb.WriteRune(esc)
			default:
				sb.WriteRune('\\')
				sb.WriteRune(esc)
			}
		} else {
			sb.WriteRune(ch)
		}
	}
	if l.pos < len(l.input) {
		l.advance()
	}
	return Token{Type: STRING_LIT, Literal: sb.String(), Line: line, Col: col}
}

func (l *Lexer) readNumber() Token {
	line, col := l.line, l.col
	var sb strings.Builder
	isFloat := false
	for l.pos < len(l.input) && (isDigit(l.peek()) || l.peek() == '.') {
		if l.peek() == '.' {
			if isFloat {
				break
			}
			isFloat = true
		}
		sb.WriteRune(l.advance())
	}
	if isFloat {
		return Token{Type: FLOAT_LIT, Literal: sb.String(), Line: line, Col: col}
	}
	return Token{Type: INT_LIT, Literal: sb.String(), Line: line, Col: col}
}

func (l *Lexer) readIdent() Token {
	line, col := l.line, l.col
	var sb strings.Builder
	for l.pos < len(l.input) && isAlphaNum(l.peek()) {
		sb.WriteRune(l.advance())
	}
	lit := sb.String()
	typ := LookupIdent(lit)
	return Token{Type: typ, Literal: lit, Line: line, Col: col}
}

func (l *Lexer) tok(t TokenType, lit string) Token {
	return Token{Type: t, Literal: lit, Line: l.line, Col: l.col}
}

func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	for l.pos < len(l.input) {
		l.skipWhitespace()
		if l.pos >= len(l.input) {
			break
		}

		ch := l.peek()

		// comments
		if ch == '/' && (l.peekAt(1) == '/' || l.peekAt(1) == '*') {
			l.skipComment()
			continue
		}

		// newlines as statement separators
		if ch == '\n' {
			line, col := l.line, l.col
			l.advance()
			tokens = append(tokens, Token{Type: NEWLINE, Literal: "\n", Line: line, Col: col})
			continue
		}

		if ch == '"' || ch == '\'' {
			tokens = append(tokens, l.readString())
			continue
		}

		if isDigit(ch) {
			tokens = append(tokens, l.readNumber())
			continue
		}

		if isAlpha(ch) {
			tokens = append(tokens, l.readIdent())
			continue
		}

		line, col := l.line, l.col
		l.advance()

		switch ch {
		case '=':
			if l.peek() == '=' {
				l.advance()
				tokens = append(tokens, Token{EQ, "==", line, col})
			} else {
				tokens = append(tokens, Token{ASSIGN, "=", line, col})
			}
		case '!':
			if l.peek() == '=' {
				l.advance()
				tokens = append(tokens, Token{NEQ, "!=", line, col})
			} else {
				tokens = append(tokens, Token{BANG, "!", line, col})
			}
		case '<':
			if l.peek() == '=' {
				l.advance()
				tokens = append(tokens, Token{LTE, "<=", line, col})
			} else {
				tokens = append(tokens, Token{LT, "<", line, col})
			}
		case '>':
			if l.peek() == '=' {
				l.advance()
				tokens = append(tokens, Token{GTE, ">=", line, col})
			} else {
				tokens = append(tokens, Token{GT, ">", line, col})
			}
		case '&':
			if l.peek() == '&' {
				l.advance()
				tokens = append(tokens, Token{AND, "&&", line, col})
			} else {
				return nil, fmt.Errorf("line %d: unexpected character '&'", line)
			}
		case '|':
			if l.peek() == '|' {
				l.advance()
				tokens = append(tokens, Token{OR, "||", line, col})
			} else {
				return nil, fmt.Errorf("line %d: unexpected character '|'", line)
			}
		case '+':
			if l.peek() == '=' {
				l.advance()
				tokens = append(tokens, Token{PLUS_ASSIGN, "+=", line, col})
			} else {
				tokens = append(tokens, Token{PLUS, "+", line, col})
			}
		case '-':
			if l.peek() == '>' {
				l.advance()
				tokens = append(tokens, Token{ARROW, "->", line, col})
			} else if l.peek() == '=' {
				l.advance()
				tokens = append(tokens, Token{MINUS_ASSIGN, "-=", line, col})
			} else {
				tokens = append(tokens, Token{MINUS, "-", line, col})
			}
		case '*':
			tokens = append(tokens, Token{STAR, "*", line, col})
		case '/':
			tokens = append(tokens, Token{SLASH, "/", line, col})
		case '%':
			tokens = append(tokens, Token{PERCENT, "%", line, col})
		case '(':
			tokens = append(tokens, Token{LPAREN, "(", line, col})
		case ')':
			tokens = append(tokens, Token{RPAREN, ")", line, col})
		case '{':
			tokens = append(tokens, Token{LBRACE, "{", line, col})
		case '}':
			tokens = append(tokens, Token{RBRACE, "}", line, col})
		case '[':
			tokens = append(tokens, Token{LBRACKET, "[", line, col})
		case ']':
			tokens = append(tokens, Token{RBRACKET, "]", line, col})
		case ',':
			tokens = append(tokens, Token{COMMA, ",", line, col})
		case ';':
			tokens = append(tokens, Token{SEMICOLON, ";", line, col})
		case ':':
			tokens = append(tokens, Token{COLON, ":", line, col})
		case '.':
			tokens = append(tokens, Token{DOT, ".", line, col})
		default:
			return nil, fmt.Errorf("line %d: unexpected character '%c'", line, ch)
		}
	}
	tokens = append(tokens, Token{Type: EOF, Literal: "", Line: l.line, Col: l.col})
	return tokens, nil
}

func isDigit(ch rune) bool  { return ch >= '0' && ch <= '9' }
func isAlpha(ch rune) bool  { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' }
func isAlphaNum(ch rune) bool { return isAlpha(ch) || isDigit(ch) }
