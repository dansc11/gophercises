package main

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestPlay(t *testing.T) {
	// GIVEN
	filename := "test.csv"

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

	outputBuffer := bytes.Buffer{}

	config := gameConfig{
		ProblemsPath: filename,
		Reader:       os.Stdin,
		Writer:       &outputBuffer,
	}

	csvString := "1+1,2\n"
	csvData := []byte(csvString)

	if err := os.WriteFile(filename, csvData, fs.ModePerm); err != nil {
		log.Fatal(err)
	}

	defer os.Remove(filename)

	// WHEN
	err = play(config)

	// THEN
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	writerGot := outputBuffer.String()
	writerWant := "1+1=1/1 problems were answered correctly\n"
	if writerGot != writerWant {
		t.Fatalf("Wanted %s, got %s", writerGot, writerWant)
	}
}

func TestPlayLoadProblemsError(t *testing.T) {
	// GIVEN
	filename := "test.csv"

	config := gameConfig{
		ProblemsPath: filename,
	}

	// WHEN
	err := play(config)

	// THEN
	if err == nil {
		t.Fatal("Expected error reading missing CSV")
	}
}

func TestLoadProblems(t *testing.T) {
	// GIVEN
	filename := "test.csv"

	problemData := [][]string{
		{"5+5", "10"},
		{"7+3", "10"},
		{"1+1", "2"},
	}

	var csvString string

	for _, problem := range problemData {
		problemString := strings.Join(problem, ",") + "\n"
		csvString += problemString
	}

	csvData := []byte(csvString)

	if err := os.WriteFile(filename, csvData, fs.ModePerm); err != nil {
		log.Fatal(err)
	}

	defer os.Remove(filename)

	// WHEN
	problems, err := loadProblems(filename)

	// THEN
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	// Check all problems have been loaded
	gotProblemsLen := len(problems)
	wantProblemsLen := 3

	if gotProblemsLen != wantProblemsLen {
		t.Fatalf("Expected %d problems, got %d", gotProblemsLen, wantProblemsLen)
	}

	// Check the correct question and answer have been loaded for each problem, and they are in order
	for i, problem := range problemData {
		wantQuestion := problem[0]
		wantAnswer := problem[1]

		gotQuestion := problems[i].Question
		gotAnswer := problems[i].Answer

		if wantQuestion != gotQuestion {
			t.Fatalf("Expected question %s, got %s for problem %d", wantQuestion, gotQuestion, i)
		}

		if wantAnswer != gotAnswer {
			t.Fatalf("Expected answer %s, got %s for problem %d", wantAnswer, gotAnswer, i)
		}
	}
}

func TestLoadProblemsCsvOpenError(t *testing.T) {
	// GIVEN
	filename := "test.csv"

	// WHEN
	problems, err := loadProblems(filename)

	// THEN
	if err == nil {
		t.Fatal("Error not raised for missing CSV")
	}

	if len(problems) != 0 {
		t.Fatal("Unexpected problems returned")
	}
}

func TestLoadProblemsCsvDataReaderError(t *testing.T) {
	// GIVEN
	// Prepare test CSV
	filename := "test.csv"

	// Invalid CSV data
	csvString := "a\""
	csvData := []byte(csvString)

	if err := os.WriteFile(filename, csvData, fs.ModePerm); err != nil {
		log.Fatal(err)
	}

	defer os.Remove(filename)

	// WHEN
	problems, err := loadProblems(filename)

	// THEN
	if err == nil {
		t.Fatal("Error not raised when reading invalid CSV data")
	}

	if len(problems) != 0 {
		t.Fatal("Unexpected problems returned")
	}
}
