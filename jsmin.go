// Idomatic Port of Douglas Crockford's (http://crockford.com/javascript/jsmin) JSMin.

package jsmin

import (
	"bytes"
	"errors"
	"io"
	"text/scanner"
)

var (
	ErrUnterminatedComment  = errors.New("Unterminated comment.")
	ErrUnterminatedString   = errors.New("Unterminated string literal.")
	ErrUnterminatedRegExSet = errors.New("Unterminated set in Regular Expression literal.")
	ErrUnterminatedRegEx    = errors.New("Unterminated Regular Expression literal.")
)

func Minify(js io.Reader) (*bytes.Buffer, error) {

	c := Compiler{o: new(bytes.Buffer), s: new(scanner.Scanner)}
	c.s.Init(js)
	return c.o, c.Compile()
}

type Compiler struct {
	o *bytes.Buffer //output

	s *scanner.Scanner
	a rune
	b rune

	x rune
	y rune
}

func (c *Compiler) Compile() error {
	c.a = '\n'
	err := c.action3()
	if err != nil {
		return err
	}
	for c.a != scanner.EOF {
		switch c.a {
		case ' ':
			if alphanum(c.b) {
				c.action1()
			} else {
				c.action2()
			}
		case '\n':
			switch c.b {
			case '{', '[', '(', '+', '-', '!', '~':
				c.action1()
			case ' ':
				c.action3()
			default:
				if alphanum(c.b) {
					c.action1()
				} else {
					c.action2()
				}
			}
		default:
			switch c.b {
			case ' ':
				if alphanum(c.a) {
					c.action1()
				} else {
					c.action3()
				}
			case '\n':
				switch c.a {
				case '}', ']', ')', '+', '-', '"', '\'', '`':
					c.action1()
				default:
					if alphanum(c.a) {
						c.action1()
					} else {
						c.action3()
					}
				}
			default:
				c.action1()
			}
		}
	}
	return nil
}

func (c *Compiler) next() (rune, error) {
	next := c.s.Next()
	if next == '/' {
		switch c.s.Peek() {
		case '/':
			for next = c.s.Next(); next != '\n'; next = c.s.Next() {
			}
		case '*':
			for next = c.s.Next(); next != ' '; next = c.s.Next() {
				switch c.s.Next() {
				case '*':
					if c.s.Peek() == '/' {
						c.s.Next()
						next = ' '
					}
				case scanner.EOF:
					return next, ErrUnterminatedComment
				}
			}
		}
	}
	c.y = c.x
	c.x = next

	return next, nil
}

func (c *Compiler) action1() error {
	c.o.WriteRune(c.a)
	if (c.y == '\n' || c.y == ' ') &&
		(c.a == '+' || c.a == '-' || c.a == '*' || c.a == '/') &&
		(c.b == '+' || c.b == '-' || c.b == '*' || c.b == '/') {
		c.o.WriteRune(c.y)
	}
	return c.action2()
}

func (c *Compiler) action2() error {
	var err error
	c.a = c.b

	if c.a == '\'' || c.a == '"' || c.a == '`' {
		for {
			c.o.WriteRune(c.a)
			c.a, err = c.next()
			if err != nil {
				return err
			}
			if c.a == c.b {
				break
			}

			if c.a == '\\' {
				c.o.WriteRune(c.a)
				c.a, err = c.next()
				if err != nil {
					return err
				}
			}

			if c.a == scanner.EOF {
				return ErrUnterminatedString
			}
		}
	}
	return c.action3()
}

func (c *Compiler) action3() error {
	var err error
	c.b, err = c.next()
	if err != nil {
		return err
	}
	if c.b == '/' && (c.a == '(' || c.a == ',' || c.a == '=' || c.a == ':' ||
		c.a == '[' || c.a == '!' || c.a == '&' || c.a == '|' ||
		c.a == '?' || c.a == '+' || c.a == '-' || c.a == '~' ||
		c.a == '*' || c.a == '/' || c.a == '{' || c.a == '\n') {

		c.o.WriteRune(c.a)
		c.o.WriteRune(c.b)

		for {
			c.a = c.s.Next()

			if c.a == '[' {
				for {
					c.o.WriteRune(c.a)
					c.a = c.s.Next()
					if c.a == ']' {
						break
					}

					if c.a == '\\' {
						c.o.WriteRune(c.a)
						c.a = c.s.Next()
					}
					if c.a == scanner.EOF {
						return ErrUnterminatedRegExSet
					}
				}
			} else if c.a == '/' {
				if peek := c.s.Peek(); peek == '/' || peek == '*' {
					return ErrUnterminatedRegExSet
				}
				break

			} else if c.a == '\\' {
				c.o.WriteRune(c.a)
				c.a = c.s.Next()
			}
			if c.a == scanner.EOF {
				return ErrUnterminatedRegExSet
			}
			c.o.WriteRune(c.a)
		}
		c.b, err = c.next()
		if err != nil {
			return err
		}
	}
	return nil
}

func alphanum(c rune) bool {
	return ((c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_' ||
		c == '$' ||
		c == '\\' ||
		c > 126)
}
