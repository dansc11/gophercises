package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var config gameConfig

	pathArg := flag.String("p", "problems.csv", "Relative or absolute path to the problems CSV file")
	timeArg := flag.Int("t", 30, "Time limit in seconds for the quiz")

	flag.Parse()

	config.ProblemsPath = *pathArg
	config.Reader = os.Stdin
	config.Writer = os.Stdout
	config.TimeLimit = *timeArg

	fmt.Print("Press 'Enter' to start...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	fmt.Printf("You have %d seconds... GO!\n", config.TimeLimit)

	timer := time.NewTimer(time.Duration(config.TimeLimit) * time.Second)

	go func() {
		<-timer.C
		fmt.Println("Time's up!")
		os.Exit(0)
	}()

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

	fmt.Fprintf(config.Writer, "%d/%d problems were answered correctly\n", len(correctAnswers), len(problems))

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
