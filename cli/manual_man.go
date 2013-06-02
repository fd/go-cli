package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/knieriem/markdown"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	mmparser = markdown.NewParser(nil)
)

func (m *Manual) Open() {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("groffer", "--tty")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	m.format_man(w)
	w.Close()
	r.Close()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

func (m *Manual) format_man(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	buf := bufio.NewWriter(w)
	defer buf.Flush()

	fmt.Fprintf(buf, ".TH command 1\n")
	m.name(buf)
	m.synopsis(buf)
	m.flags(buf)
	m.args(buf)
	m.vars(buf)

	for _, s := range m.sections {
		s.Man(buf)
	}
}

func (m *Manual) name(w *bufio.Writer) {
	s := section_t{
		Header: "NAME",
		Body:   m.summary,
	}
	s.Man(w)
}

func (m *Manual) synopsis(w *bufio.Writer) {
	s := section_t{
		Header: "SYNOPSIS",
		Body:   m.usage,
	}
	s.Man(w)
}

func (m *Manual) flags(w *bufio.Writer) {
	var (
		buf  bytes.Buffer
		exec = m.exec
	)

	fmt.Fprintf(w, ".SH OPTIONS\n")

	mmf := markdown.ToGroffMM(&buf)

	for exec != nil {
		for _, f := range exec.Flags {
			s, p := m.options[f.Name]
			if !p {
				continue
			}

			buf.Reset()
			mmparser.Markdown(bytes.NewReader([]byte(s.Body)), mmf)
			body := buf.String()
			body = strings.TrimPrefix(body, ".P\n")

			fmt.Fprintf(w, ".TP\n\\fB%s\\fP\n%s", f.Tag.Get("flag"), body)
		}
		exec = exec.ParentExec
	}
}

func (m *Manual) vars(w *bufio.Writer) {
	var (
		buf bytes.Buffer
	)

	fmt.Fprintf(w, ".SH \"ENVIRONMENT VARIABLES\"\n")

	mmf := markdown.ToGroffMM(&buf)

	for _, f := range m.exec.Variables {
		s, p := m.options[f.Name]
		if !p {
			continue
		}

		if flag := f.Tag.Get("flag"); flag != "" {
			s.Body = fmt.Sprintf("See %s", flag)
		}

		buf.Reset()
		mmparser.Markdown(bytes.NewReader([]byte(s.Body)), mmf)
		body := buf.String()
		body = strings.TrimPrefix(body, ".P\n")

		fmt.Fprintf(w, ".TP\n\\fB%s\\fP\n%s", f.Tag.Get("env"), body)
	}
}

func (m *Manual) args(w *bufio.Writer) {
	var (
		buf bytes.Buffer
	)

	fmt.Fprintf(w, ".SH ARGUMENTS\n")

	mmf := markdown.ToGroffMM(&buf)

	for _, f := range m.exec.Args {
		s, p := m.options[f.Name]
		if !p {
			continue
		}

		buf.Reset()
		mmparser.Markdown(bytes.NewReader([]byte(s.Body)), mmf)
		body := buf.String()
		body = strings.TrimPrefix(body, ".P\n")

		fmt.Fprintf(w, ".TP\n\\fB%s\\fP\n%s", f.Name, body)
	}
}

func (s *section_t) Man(w *bufio.Writer) {
	fmt.Fprintf(w, ".SH %s\n", strconv.Quote(s.Header))

	f := markdown.ToGroffMM(w)
	mmparser.Markdown(bytes.NewReader([]byte(s.Body)), f)
}
