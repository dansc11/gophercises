package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type problem struct {
	Question string
	Answer   string
}

func (p *problem) Ask(ioReader io.Reader, ioWriter io.Writer) bool {
	reader := bufio.NewReader(ioReader)

	fmt.Fprintf(ioWriter, "%s=", p.Question)

	text, _ := reader.ReadString('\n')
	answer := strings.Replace(text, "\n", "", -1)

	return p.IsCorrect(answer)
}

func (p *problem) IsCorrect(answer string) bool {
	return p.Answer == answer
}
