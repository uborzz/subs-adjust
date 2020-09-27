package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const time_regex = `(?:(?:([01]?\d|2[0-3]):)?([0-5]?\d):)?([0-5]?\d),([0-9]*)`

func parseTime(time_extraction []string) time.Duration {
	t := time_extraction
	parsed := fmt.Sprintf("%sh%sm%ss%sms", t[1], t[2], t[3], t[4])
	result, _ := time.ParseDuration(parsed)
	return result
}

func formatTime(time_duration time.Duration) string {
	t := time_duration

	if t.Milliseconds() <= 0 {
		return "00:00:00,000"
	} else {
		return fmt.
			Sprintf("%02d:%02d:%02d,%03d",
				int(t.Hours()),
				int(t.Minutes())%60,
				int(t.Seconds())%60,
				int(t.Milliseconds())%1000,
			)
	}

}

func modifySubs(in_file string, out_file string, seconds float64) error {

	seconds_time := time.Duration(1000 * 1000 * 1000 * seconds) // to nanos

	data, err := ioutil.ReadFile(in_file)
	if err != nil {
		return err
	}

	open_file, err := os.OpenFile(out_file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(open_file)

	r, _ := regexp.Compile(time_regex + " --> " + time_regex)
	t, _ := regexp.Compile(time_regex)

	ss := strings.Split(string(data), "\n")
	for _, s := range ss {

		if r.MatchString(s) {

			groups := t.FindAllStringSubmatch(s, 2)

			t1 := groups[0]
			parsed_t1 := parseTime(t1)
			new_t1 := formatTime(parsed_t1 + seconds_time)

			t2 := groups[1]
			parsed_t2 := parseTime(t2)
			new_t2 := formatTime(parsed_t2 + seconds_time)

			line := fmt.Sprintln(new_t1, " --> ", new_t2)
			datawriter.WriteString(line)

		} else {
			datawriter.WriteString(s + "\n")
		}

	}

	datawriter.Flush()
	open_file.Close()

	return nil
}

func main() {

	if len(os.Args) != 4 {
		fmt.Println("Usage: subs-adjust INPUT OUTPUT SECONDS")
		os.Exit(1)
	}

	// Params
	in_file := os.Args[1]
	out_file := os.Args[2]
	seconds, _ := strconv.ParseFloat(os.Args[3], 64)

	modifySubs(in_file, out_file, (seconds))

	fmt.Println("Done!")

}
