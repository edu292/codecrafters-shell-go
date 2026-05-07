package parser

type ParserState int

const (
	StateNormal ParserState = iota
	InDoubleQuotes
	InSingleQuotes
)

type (
	Op  int
	Cmd struct {
		Name string
		Args []string
		Next Op
	}
)

const (
	OpNil Op = iota
	OpRedirectStdOut
	OpRedirectStdErr
	OpRedirectBoth
)

func ParseInput(input []byte) []*Cmd {
	var commands []*Cmd
	var fields []string
	var buf []byte

	parserState := StateNormal
	isEscaped := false
	for _, b := range input {
		switch parserState {
		case InDoubleQuotes:
			if isEscaped {
				buf = append(buf, b)
				isEscaped = false
				break
			}

			switch b {
			case '"':
				parserState = StateNormal
			case '\\':
				isEscaped = true
			default:
				buf = append(buf, b)
			}
		case InSingleQuotes:
			if b == '\'' {
				parserState = StateNormal
			} else {
				buf = append(buf, b)
			}
		case StateNormal:
			if isEscaped {
				buf = append(buf, b)
				isEscaped = false
				break
			}

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
				isEscaped = true
			case '>':
				op := OpRedirectStdOut
				if len(buf) > 0 {
					matched := true
					switch buf[len(buf)-1] {
					case '1':
						op = OpRedirectStdOut
					case '2':
						op = OpRedirectStdErr
					case '&':
						op = OpRedirectBoth
					default:
						matched = false
					}

					if matched {
						buf = buf[:len(buf)-1]
					}

					if len(buf) > 0 {
						fields = append(fields, string(buf))
					}
				}

				buf = buf[:0]
				args := append([]string(nil), fields[1:]...)
				commands = append(commands, &Cmd{fields[0], args, op})
				clear(fields)
				fields = fields[:0]
			default:
				buf = append(buf, b)
			}
		}
	}

	if len(buf) > 0 {
		fields = append(fields, string(buf))
	}

	if len(fields) > 0 {
		commands = append(commands, &Cmd{fields[0], fields[1:], OpNil})
	}

	return commands
}
