package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	filename string = "problems.csv"
	duration int    = 30
)

type problem struct {
	q string
	a string
}

func main() {
	n := flag.String("csv", filename, "a csv file in the format of 'question,answer'")
	d := flag.Int("duration", duration, "the duration of the quiz in seconds")
	m := flag.Bool("mix", false, "get results in random order")
	flag.Parse()

	if err := run(n, d, m); err != nil {
		fmt.Println("Error: ", err)
	}
}

func run(n *string, d *int, m *bool) error {
	f, err := os.Open(*n)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	start := ""
	fmt.Println("Press enter to start the quiz (and timer)")
	fmt.Scanln(&start)

	problems := parseRecords(records)

	if *m {
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	correct := printAndCount(d, problems)

	fmt.Println("Correct answers: ", correct)
	fmt.Println("Total questions: ", len(problems))

	return nil
}

func printAndCount(d *int, problems []problem) int {
	correct := 0
	timer := time.NewTimer(time.Duration(*d) * time.Second)

	for _, p := range problems {
		fmt.Println(p.q)

		c := make(chan int, 1)

		go func() {
			var a string
			fmt.Scanln(&a)

			if strings.TrimSpace(a) == p.a {
				c <- 1
			}

			c <- 0
		}()

		select {
		case <-timer.C:
			fmt.Println("Time's up!")
			return correct
		case i := <-c:
			correct += i
		}
	}

	return correct
}

func parseRecords(records [][]string) []problem {
	p := make([]problem, len(records))

	for i, record := range records {
		p[i] = problem{
			q: record[0],
			a: strings.TrimSpace(record[1]),
		}
	}

	return p
}
