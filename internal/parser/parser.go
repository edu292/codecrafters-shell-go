package parser

import (
	"bufio"
	"bytes"
)

type ParserState int

const (
	StateNormal ParserState = iota
	InDoubleQuotes
	InSingleQuotes
	CheckToken
)

type (
	Op  uint8
	Cmd struct {
		Name string
		Args []string
		Op   Op
	}
)

const Nil Op = 0
const (
	// Streams
	Stdout Op = 1 << iota // 00001
	Stderr                // 00010

	// Actions
	Redir // 00100
	Pipe  // 01000

	// Mods
	Append // 10000
)

const (
	Both           = Stdout | Stderr
	RedirOut       = Redir | Stdout
	RedirErr       = Redir | Stderr
	RedirOutAppend = RedirOut | Append
	RedirErrAppend = RedirErr | Append
)

func Is(o Op, p Op) bool {
	return o&p == p
}

func tokenToOp(t []byte) Op {
	switch {
	case bytes.Equal(t, []byte("1>>")):
		return RedirOutAppend
	case bytes.Equal(t, []byte(">>")):
		return RedirOutAppend
	case bytes.Equal(t, []byte("2>>")):
		return RedirErrAppend
	case bytes.Equal(t, []byte("1>")):
		return RedirOut
	case bytes.Equal(t, []byte(">")):
		return RedirOut
	case bytes.Equal(t, []byte("2>")):
		return RedirErr
	default:
		return Nil
	}
}

func ParseInput(input []byte) []*Cmd {
	rb := bytes.NewReader(input)
	r := bufio.NewReaderSize(rb, len(input))

	var commands []*Cmd
	var fields []string
	var buf []byte

	parserState := StateNormal
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}

		switch parserState {
		case InDoubleQuotes:
			switch b {
			case '"':
				parserState = StateNormal
			case '\\':
				escByte, err := r.ReadByte()
				if err != nil {
					buf = append(buf, b)
					break
				}

				switch escByte {
				case '"', '\\', '$', '`', 'n':
					buf = append(buf, escByte)
				default:
					buf = append(buf, b)
					buf = append(buf, escByte)
				}
			default:
				buf = append(buf, b)
			}
		case InSingleQuotes:
			if b == '\'' {
				parserState = StateNormal
				break
			}

			buf = append(buf, b)
		case CheckToken:
			if b == '>' {
				buf = append(buf, b)
			} else {
				r.UnreadByte()
			}

			op := tokenToOp(buf)
			buf = buf[:0]

			args := append([]string(nil), fields[1:]...)
			commands = append(commands, &Cmd{fields[0], args, op})

			clear(fields)
			fields = fields[:0]
			parserState = StateNormal
		case StateNormal:
			switch b {
			case ' ':
				if len(buf) > 0 {
					fields = append(fields, string(buf))
				}

				buf = buf[:0]
			case '\'':
				parserState = InSingleQuotes
			case '"':
				parserState = InDoubleQuotes
			case '\\':
				escByte, err := r.ReadByte()
				if err == nil {
					buf = append(buf, escByte)
				}
			case '1', '2':
				n, err := r.ReadByte()

				if n != '>' {
					buf = append(buf, b)
					if err == nil {
						r.UnreadByte()
					}
					break
				}

				parserState = CheckToken
				if len(buf) > 0 {
					fields = append(fields, string(buf))
				}

				buf = buf[:0]
				buf = append(buf, b)
				buf = append(buf, n)
			case '>':
				parserState = CheckToken
				if len(buf) > 0 {
					fields = append(fields, string(buf))
				}

				buf = buf[:0]
				buf = append(buf, b)
			default:
				buf = append(buf, b)
			}
		}
	}

	if len(buf) > 0 {
		fields = append(fields, string(buf))
	}

	if len(fields) > 0 {
		commands = append(commands, &Cmd{fields[0], fields[1:], Nil})
	}

	return commands
}
