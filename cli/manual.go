package cli

import (
	"fmt"
	"strings"
	"unicode"
)

type Manual struct {
	Usage      string
	Summary    string
	Paragraphs []Paragraph
}

type Paragraph struct {
	Header string
	Body   string
}

func (m *Manual) Parse(source string, handles map[string][]string) {
	var (
		indent = determine_indent(source)
		lines  = strings.Split(source, "\n")
	)

	remove_indent(lines, indent)

	var (
		in_section   bool
		section_name string
		section_body string
	)
	for _, line := range lines {
		name, body, empty := parse_line(line)

		if empty {
			if in_section {
				section_body += "\n"
			}
			continue
		}

		if name != "" {
			if in_section {
				m.parse_section(section_name, section_body, handles)
				section_name = ""
				section_body = ""
			}

			in_section = true
			section_name = name
		}

		section_body += body + "\n"
	}

	if in_section {
		m.parse_section(section_name, section_body, handles)
	}
}

func (m *Manual) parse_section(name, body string, handles map[string][]string) {
	{
		var (
			indent = determine_indent(body)
			lines  = strings.Split(body, "\n")
		)

		remove_indent(lines, indent)
		body = strings.TrimSpace(strings.Join(lines, "\n"))
	}

	switch name {

	case "Usage":
		m.Usage = body

	case "Summary":
		m.Summary = body

	default:

		if strings.HasPrefix(name, ".") {
			m.parse_option(name[1:], body, handles)

		} else {
			p := Paragraph{Header: name, Body: body}
			m.Paragraphs = append(m.Paragraphs, p)
		}

	}
}

func (m *Manual) parse_option(name, body string, handles map[string][]string) {
	p := Paragraph{Body: body}
	p.Header = strings.Join(handles[name], " ")
	m.Paragraphs = append(m.Paragraphs, p)
}

func parse_line(line string) (section, body string, empty bool) {
	if len(line) == 0 || is_space_only(line) {
		return "", "", true
	}

	if unicode.IsSpace([]rune(line)[0]) {
		return "", line, false
	}

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		panic(fmt.Sprintf("Invalid indentation for line: `%s`", line))
	}

	return strings.TrimSpace(parts[0]), strings.TrimLeft(parts[1], " \t"), false
}

func remove_indent(lines []string, n int) {
	for i, line := range lines {
		lines[i] = skip_at_most_n_spaces(line, n)
	}
}

func determine_indent(source string) int {
	var (
		indent int
	)

	for _, c := range source {
		if c == '\n' {
			indent = 0
		} else if unicode.IsSpace(c) {
			indent += 1
		} else {
			break
		}
	}

	return indent
}

func skip_at_most_n_spaces(line string, n int) string {
	var (
		prefix string
		suffix string
	)

	if len(line) < n {
		prefix = line
	} else {
		prefix = line[:n]
		suffix = line[n:]
	}

	if is_space_only(prefix) {
		return suffix
	}

	panic(fmt.Sprintf("Invalid indentation for line: `%s`", line))
}

func is_space_only(s string) bool {
	for _, c := range s {
		if !unicode.IsSpace(c) {
			return false
		}
	}
	return true
}
