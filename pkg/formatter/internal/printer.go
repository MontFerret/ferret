package internal

import (
	"io"
	"strings"
)

type printer struct {
	out             io.Writer
	opts            *Options
	indent          int
	atLineStart     bool
	lastWasSpace    bool
	forceSingleLine bool
	sawHardNewline  bool
	err             error
}

func newPrinter(out io.Writer, opts *Options) *printer {
	return &printer{
		out:         out,
		opts:        opts,
		atLineStart: true,
	}
}

func (p *printer) Err() error {
	return p.err
}

func (p *printer) writeIndent() {
	if p.err != nil || !p.atLineStart {
		return
	}
	if p.indent <= 0 || p.opts.tabWidth == 0 {
		return
	}

	indent := strings.Repeat(" ", int(p.opts.tabWidth)*p.indent)
	_, err := io.WriteString(p.out, indent)

	if err != nil {
		p.err = err
		return
	}
}

func (p *printer) write(s string) {
	if p.err != nil || s == "" {
		return
	}

	if p.atLineStart {
		p.writeIndent()
	}

	_, err := io.WriteString(p.out, s)
	if err != nil {
		p.err = err
		return
	}

	p.atLineStart = false
	p.lastWasSpace = false
}

func (p *printer) writeRaw(s string) {
	if p.err != nil || s == "" {
		return
	}

	for _, r := range s {
		if r == '\n' {
			p.sawHardNewline = true

			if p.forceSingleLine {
				p.space()

				continue
			}

			_, err := io.WriteString(p.out, "\n")
			if err != nil {
				p.err = err

				return
			}

			p.atLineStart = true
			p.lastWasSpace = false

			continue
		}

		if p.atLineStart {
			p.atLineStart = false
		}

		_, err := io.WriteString(p.out, string(r))
		if err != nil {
			p.err = err

			return
		}

		p.lastWasSpace = r == ' '
	}
}

func (p *printer) space() {
	if p.err != nil || p.atLineStart || p.lastWasSpace {
		return
	}

	_, err := io.WriteString(p.out, " ")
	if err != nil {
		p.err = err

		return
	}

	p.lastWasSpace = true
}

func (p *printer) newline() {
	if p.err != nil {
		return
	}

	if p.forceSingleLine {
		p.space()

		return
	}

	_, err := io.WriteString(p.out, "\n")
	if err != nil {
		p.err = err

		return
	}

	p.atLineStart = true
	p.lastWasSpace = false
}

func (p *printer) withIndent(fn func()) {
	p.indent++
	fn()
	p.indent--
}
