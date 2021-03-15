package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

//
const (
	SUCCESS   = 0
	FAILURE   = 1
	FILENAME  = "^(.*)\\.(sql)$"
	INSERT    = "(INSERT INTO).*(\n.*)*(VALUES)"
	COMMENT   = "(^--.*)|([^']--.*)|(\\/\\*.*\\*\\/)"
	EMPTYLINE = "^( )*$"
)

// ErrWrongQuerry is
var ErrWrongQuerry = errors.New("wrong format")

func generateSpace(size int) (out string) {
	if size <= 0 {
		return
	}

	for i := 0; i < size; i++ {
		out = out + " "
	}

	return out
}

func getQuerry(reader *bufio.Reader) (query string, comment string, err error) {
	query, err = reader.ReadString(';')
	if err == io.EOF {
		return
	} else if err != nil {
		err = fmt.Errorf("unable to read string: %v", err)
		return
	}

	insertReg := regexp.MustCompile(INSERT)
	if insertReg.Match([]byte(query)) {
		commReg := regexp.MustCompile(COMMENT)
		comment = string(commReg.Find([]byte(query)))
		comment = strings.ReplaceAll(comment, "\n", "")
		comment = strings.ReplaceAll(comment, "\r", "")

		query = commReg.ReplaceAllString(query, "")
		query = strings.ReplaceAll(query, "\n", "")
		query = strings.ReplaceAll(query, "\r", "")
	} else {
		query = strings.TrimSpace(query)
		query = strings.Trim(query, "\n")
		query = strings.Trim(query, "\r")
	}

	return
}

func connectWronglySplitted(splits []string) (out string) {
	for _, s := range splits {
		if len(out) > 0 {
			out = fmt.Sprintf("%s,%s", out, s)
		} else {
			out = fmt.Sprintf("%s", s)
		}
	}

	return
}

func tryMerge(values []string) (merged []string) {
	for i := 0; i < len(values); i++ {
		if strings.Count(values[i], "'") == 1 {
			for j := i + 1; j < len(values); j++ {
				if strings.Contains(values[j], "'") {
					merged = append(merged, connectWronglySplitted(values[i:j+1]))
					i = j
					break
				}
			}
		} else {
			merged = append(merged, values[i])
		}
	}

	return
}

func formatQuerry(query string) (formatted string, err error) {
	format := "INSERT INTO %v\n" +
		"     VALUES %v\n"

	split := strings.Split(query, "VALUES")
	if len(split) < 2 {
		err = ErrWrongQuerry
		formatted = query
		return
	}

	split[0] = strings.Replace(split[0], "INSERT INTO", "", 1)
	splitCols := strings.Split(split[0], ",")
	splitVals := strings.Split(split[1], ",")

	if len(splitCols) != len(splitVals) {
		splitVals = tryMerge(splitVals)

		if len(splitCols) != len(splitVals) {
			err = ErrWrongQuerry
			formatted = "-- CHECK FORMAT!\n" + query + "\n-- CHECK FORMAT!"
			return
		}
	}

	var (
		columns string
		values  string
	)

	headerCols := strings.Split(splitCols[0], "(")

	columns += strings.TrimSpace(headerCols[0]) + "(" + strings.TrimSpace(headerCols[1])
	values += strings.TrimSpace(splitVals[0])
	columns += generateSpace(len(values) - len(columns))
	values = generateSpace(len(columns)-len(values)) + values

	for i := 1; i < len(splitCols); i++ {

		columns += ", " + strings.TrimSpace(splitCols[i])
		values += ", " + strings.TrimSpace(splitVals[i])

		columns += generateSpace(len(values) - len(columns))
		values += generateSpace(len(columns) - len(values))
	}

	formatted = fmt.Sprintf(format, columns, values)

	return
}

func main() {
	regFileName := regexp.MustCompile(FILENAME)
	fileName := os.Args[1]
	if !regFileName.MatchString(fileName) {
		fmt.Fprintf(os.Stderr, "wrong file name\n")
		os.Exit(FAILURE)
	}

	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open file: %v\n", err)
		os.Exit(FAILURE)
	}
	defer file.Close()
	defer file.Sync()

	scanner := bufio.NewReader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create scanner on file: %v\n", err)
		os.Exit(FAILURE)
	}

	var out []string
	for {
		query, comment, err := getQuerry(scanner)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "unable to get query: %v\n", err)
			os.Exit(FAILURE)
		}

		if len(comment) > 2 {
			out = append(out, fmt.Sprintf("%v\n", comment))
		}

		query, err = formatQuerry(query)
		if err == ErrWrongQuerry {
			out = append(out, fmt.Sprintf("%v\n", query))
			continue
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "unable to format query: %v\n", err)
			os.Exit(FAILURE)
		}

		out = append(out, fmt.Sprintf("%v", query))
	}

	err = file.Truncate(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to truncate file: %v\n", err)
		os.Exit(FAILURE)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to seek begining of file: %v\n", err)
		os.Exit(FAILURE)
	}

	writer := bufio.NewWriter(file)
	for _, o := range out {
		_, err = writer.WriteString(o)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to write to a file: %v\n", err)
			os.Exit(FAILURE)
		}
		writer.Flush()
	}

	os.Exit(SUCCESS)
}
