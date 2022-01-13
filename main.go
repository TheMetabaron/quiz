package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "A CSV in 'question,answer' format.")
	timeLimit := flag.Int("limit", 30, "the time limit for each question in seconds")
	randomize := flag.Bool("random", false, "default false, if set to true, will randomize order of questions")
	flag.Parse()
	
	// Read From CSV
	file, err := os.Open(*csvFilename)
	if err != nil {

		exit(fmt.Sprintf("Failed to open csv file: %s\n", *csvFilename))
		os.Exit(1)
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file")
	}
	problems := parseLines(lines)

	// Randomize
	if *randomize && len(problems) > 0 {
		rand.Seed(time.Now().Unix())
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	// Pause until user is presses enter
	fmt.Printf("Press enter when ready to start. Time limit is %d seconds", *timeLimit)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answerCh <- strings.ToLower(strings.Trim(answer, " "))
		}()

		select {
		case <- timer.C:
			fmt.Printf("Time Out: You scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct ++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines{
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}