package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to read CSV file.")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)
		answerCh := make(chan string)
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			ans := strings.Trim(scanner.Text(), " ")
			answerCh <- ans
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nCorrect %d out of %d.", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.answer {
				correct++
			}
		}
	}
	fmt.Printf("Correct %d out of %d.", correct, len(problems))
}

type problem struct {
	question string
	answer   string
}

func parseLines(lines [][]string) []problem {
	lst := make([]problem, len(lines))
	for i, line := range lines {
		lst[i] = problem{
			question: line[0],
			answer:   line[1],
		}
	}
	return lst
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
