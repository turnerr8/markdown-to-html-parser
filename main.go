package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ------------- Error Checks -----------------

// checks errors
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// check whether file exists
func fileExists(fileName string) {
	_, err := os.Stat(fileName)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("File does not exist.")
		os.Exit(1)
	}
}

// make sure that file is a markdown file
func checkFileType(fileName string) {
	fi, err := os.Stat(fileName)
	check(err)
	extention := fi.Name()[len(fi.Name())-2:]
	if !strings.EqualFold("md", extention) {
		println("Please only use markdown files")
		os.Exit(1)
	}
}

// -------------- Inline Level Markdown ----------------

// checks for bold markdown and returns new string with html formatted bold
func handleBold(line string) string {
	splitLine := strings.Split(line, "**")
  newLine := splitLine[0]
  for i:=1; i<len(splitLine); i++ {
    fmt.Println(newLine)
    if i%2 != 0{
      newLine += "<b>"+splitLine[i]+"</b>"     
    } else {
      newLine += splitLine[i]
    }
  }
  fmt.Println(newLine)
  return newLine

}



// checks for italic markdown and returns new string with html formatted italic
func handleItalic(line string) string {
	splitLine := strings.Split(line, "*")
  newLine := splitLine[0]
  for i:=1; i<len(splitLine); i++ {
    if i%2 != 0{
      newLine += "<i>"+splitLine[i]+"</i>"     
    } else {
      newLine += splitLine[i]
    }
  }
  return newLine
}
// --------------- Line Level Markdown ----------------



// handle headers
func handleHeader(line string) string {

	header := strings.Split(line, " ")
	return "<h" +
		strconv.Itoa(len(header[0])) +
		">" + strings.Join(header[1:], " ") +
		"</h" + strconv.Itoa(len(header[0])) +
		">"
}



// -------------- Block Level Markdown ----------------


func getFileContents(fileName string) {
	fi, err := os.Open(fileName)
	check(err)
	defer fi.Close()
	fileNameClean := fileName[:len(fileName)-3]
	fo, err := os.Create(fileNameClean + ".html")
	check(err)
	defer fo.Close()
	scanner := bufio.NewScanner(fi)
	writer := bufio.NewWriter(fo)
	lineNum := 1
	writeHeader(fileNameClean, writer)

	//go through input file and write each converted line to output file
	//then flush out writer
	for scanner.Scan() {
		handleFileContents(scanner.Text(), writer)
		lineNum++
	}
	writeFooter(writer)

	writer.Flush()
}

// write opening and header elements
func writeHeader(fileName string, writer *bufio.Writer) {
	header := "<!DOCTYPE html>\n<html lang='en'>\n<head>\n<title>\n" +
		fileName +
		"\n</title>\n</head>\n<body>"
	_, err := writer.WriteString(header)
	check(err)
}

// write closing elements
func writeFooter(writer *bufio.Writer) {
	_, err := writer.WriteString("</body>\n</html>")
	check(err)

}

// handle all file contents
func handleFileContents(line string, writer *bufio.Writer) {
	//if empty line, dont print anything to new file
	if line == "" {
		return
	}
	//inline level markdown
	line = handleBold(line) 
  line = handleItalic(line)
	//line level markdown
	switch line[0] {
	case '#':
	  line =	handleHeader(line)
	}
	//block level markdown
  writer.WriteString(line + "\n")
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("filename needed")
		os.Exit(1)
	}
	filename := args[1]

	//file checks
	fileExists(filename)
	checkFileType(filename)

	getFileContents(filename)
}
