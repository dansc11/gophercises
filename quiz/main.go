package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	var config gameConfig

	pathArg := flag.String("p", "problems.csv", "Relative or absolute path to the problems CSV file")

	flag.Parse()

	config.ProblemsPath = *pathArg
	config.Reader = os.Stdin
	config.Writer = os.Stdout

	if err := play(config); err != nil {
		panic(err)
	}
}

func play(config gameConfig) error {
	problems, err := loadProblems(config.ProblemsPath)

	if err != nil {
		return err
	}

	var correctAnswers []*problem

	for _, problem := range problems {
		if problem.Ask(config.Reader, config.Writer) {
			correctAnswers = append(correctAnswers, &problem)
		}
	}

	fmt.Fprintf(config.Writer, "%d/%d problems were answered correctly", len(correctAnswers), len(problems))

	return nil
}

func loadProblems(path string) ([]problem, error) {
	var problems []problem

	file, err := os.Open(path)

	if err != nil {
		return problems, err
	}

	defer file.Close()

	csvReader := csv.NewReader(file)

	rows, err := csvReader.ReadAll()

	if err != nil {
		return problems, err
	}

	for _, row := range rows {
		problems = append(problems, problem{
			Question: row[0],
			Answer:   row[1],
		})
	}

	return problems, nil
}
