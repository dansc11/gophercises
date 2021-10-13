package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestAsk(t *testing.T) {
	// GIVEN
	p := problem{
		Question: "1+1",
		Answer:   "2",
	}

	// Create temporary stdIn file
	tmpStdIn, err := ioutil.TempFile("", "test-stdin")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpStdIn.Name())

	if _, err := tmpStdIn.Write([]byte("2")); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpStdIn.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	// Revert to the normal stdin when done
	oldStdin := os.Stdin
	defer func() {
		os.Stdin = oldStdin
	}()

	os.Stdin = tmpStdIn

	// Create a buffer for the output
	buffer := bytes.Buffer{}

	// WHEN
	// Check function output
	got := p.Ask(os.Stdin, &buffer)

	// THEN
	want := true

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}

	// Check writer output
	writerGot := buffer.String()
	writerWant := p.Question + "="

	if writerGot != writerWant {
		t.Errorf("got %s, want %s", writerGot, writerWant)
	}

	if err := tmpStdIn.Close(); err != nil {
		log.Fatal(err)
	}
}

func TestIsCorrectTrue(t *testing.T) {
	// GIVEN
	p := problem{
		Answer: "5",
	}

	// WHEN
	got := p.IsCorrect("5")

	// THEN
	want := true

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}

func TestIsCorrectFalse(t *testing.T) {
	// GIVEN
	p := problem{
		Answer: "5",
	}

	// WHEN
	got := p.IsCorrect("4")

	// THEN
	want := false

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}
