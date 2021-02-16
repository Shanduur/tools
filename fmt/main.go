package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	SUCCESS = 0
	FAILURE = 1
)

func generateSpace(size int) (out string) {
	if size <= 0 {
		return
	}

	for i := 0; i < size; i++ {
		out = out + " "
	}

	return out
}

func main() {
	template := "^(.*)\\.(sql)$"

	regFileName, err := regexp.Compile(template)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to compile regexp: %v", err)
		os.Exit(FAILURE)
	}

	fileName := os.Args[1]
	if !regFileName.MatchString(fileName) {
		fmt.Fprintf(os.Stderr, "wrong file name")
		os.Exit(FAILURE)
	}

	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open file: %v", err)
		os.Exit(FAILURE)
	}
	defer file.Close()
	defer file.Sync()

	var out []string

	reader := bufio.NewReader(file)
	for {
		var (
			columns string
			values  string
		)

		line1, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "unable to read contents of file: %v", err)
			os.Exit(FAILURE)
		}

		if !strings.Contains(line1, "INSERT INTO") {
			continue
		}

		line2, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "unable to read contents of file: %v", err)
			os.Exit(FAILURE)
		}
		if !strings.Contains(line2, "VALUES") {
			continue
		}

		split1 := strings.Split(line1, ",")
		split2 := strings.Split(line2, ",")

		for i := 0; i < len(split1); i++ {
			if len(columns) < 11 {
				columns += split1[i]
			} else {
				columns += ", " + split1[i]
			}

			if len(values) < 6 {
				values += split2[i]
			} else {
				values += ", " + split2[i]
			}

			values += generateSpace(len(columns) - len(values))
			columns += generateSpace(len(values) - len(columns))
		}

		columns = strings.ReplaceAll(columns, "\n", "")
		values = strings.ReplaceAll(values, "\n", "")
		out = append(out, fmt.Sprintf("%v\n%v\n", columns, values))
	}

	err = file.Truncate(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to truncate file: %v", err)
		os.Exit(FAILURE)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to seek begining of file: %v", err)
		os.Exit(FAILURE)
	}

	writer := bufio.NewWriter(file)
	for _, o := range out {
		_, err = writer.WriteString(o)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to write to a file: %v", err)
			os.Exit(FAILURE)
		}
		writer.Flush()
	}

	os.Exit(SUCCESS)
}
