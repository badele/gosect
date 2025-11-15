package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Section represents a found section in the content
type Section struct {
	Name     string
	StartIdx int
	EndIdx   int
	SrcFile  string
	Content  string
}

// initial regex patterns
var (
	reBegin = regexp.MustCompile(`(?m)BEGIN SECTION ([A-Za-z0-9_-]+)(?: file=([^ >]+))?`) // captures name + optional file
	reEnd   = regexp.MustCompile(`(?m)END SECTION ([A-Za-z0-9_-]+)`)                      // captures name
)

func makeRegex(begin, end string) (*regexp.Regexp, *regexp.Regexp) {
	b := regexp.MustCompile("(?m)" + regexp.QuoteMeta(begin) + ` ([A-Za-z0-9_-]+)(?: file=([^ >]+))?`)
	e := regexp.MustCompile("(?m)" + regexp.QuoteMeta(end) + ` ([A-Za-z0-9_-]+)`)

	return b, e
}

// /////////////////////////////////////////////////////////////////////////////
// find all sections in content
// /////////////////////////////////////////////////////////////////////////////
func findSections(content string, reBegin, reEnd *regexp.Regexp) ([]Section, error) {

	begins := reBegin.FindAllStringSubmatchIndex(content, -1)
	ends := reEnd.FindAllStringSubmatchIndex(content, -1)

	var sections []Section

	for _, b := range begins {
		name := content[b[2]:b[3]]
		file := ""
		if b[4] != -1 && b[5] != -1 {
			file = content[b[4]:b[5]]
		}

		// find corresponding END
		endIdx := -1
		for _, e := range ends {
			endName := content[e[2]:e[3]]
			if endName == name && e[0] > b[1] {
				endIdx = e[0]
				break
			}
		}

		if endIdx == -1 {
			return nil, fmt.Errorf("no END SECTION for %s", name)
		}

		sections = append(sections, Section{
			Name:     name,
			StartIdx: b[0],
			EndIdx:   endIdx,
			SrcFile:  file,
		})
	}

	return sections, nil
}

func replaceSections(content string, sections []Section, verbose bool, reBegin, reEnd *regexp.Regexp) (string, error) {

	out := content
	offset := 0

	for _, s := range sections {
		if s.SrcFile == "" {
			return "", fmt.Errorf("section %s has no file= source", s.Name)
		}

		b, err := os.ReadFile(s.SrcFile)
		if err != nil {
			return "", err
		}
		src := strings.TrimSpace(string(b))
		if verbose {
			fmt.Fprintf(os.Stderr, "[gosect] section=%s source=%s\n", s.Name, s.SrcFile)
		}

		// reconstruct - find end of BEGIN line and start of END line
		beginPos := s.StartIdx + offset
		endPos := s.EndIdx + offset

		// Trouver la fin de la ligne BEGIN (jusqu'au \n)
		endOfBeginLine := strings.Index(out[beginPos:], "\n")
		if endOfBeginLine == -1 {
			return "", fmt.Errorf("malformed BEGIN line for section %s", s.Name)
		}
		endOfBeginLine += beginPos

		// Trouver le début de la ligne END (depuis le dernier \n)
		startOfEndLine := strings.LastIndex(out[:endPos], "\n")
		if startOfEndLine == -1 {
			startOfEndLine = 0
		} else {
			startOfEndLine++ // garder le \n
		}

		before := out[:endOfBeginLine+1] // +1 pour inclure le \n
		after := out[startOfEndLine:]

		newBlock := before + "\n" + src + "\n\n" + after

		delta := len(newBlock) - len(out)
		out = newBlock
		offset += delta
	}

	return out, nil
}

// entry point
func main() {
	// Get command-line flags
	beginFlag := flag.String("begin", "BEGIN SECTION", "begin marker prefix")
	endFlag := flag.String("end", "END SECTION", "end marker prefix")
	filePath := flag.String("file", "", "input file path")
	stdout := flag.Bool("stdout", false, "print to stdout instead of writing file")
	verbose := flag.Bool("verbose", false, "log details about processed sections")

	flag.Parse()

	// Validate required flags
	if *filePath == "" {
		fmt.Fprintln(os.Stderr, "-file required")
		os.Exit(1)
	}

	// Read input file
	inputBytes, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	input := string(inputBytes)

	// Create regex patterns based on flags
	reBegin, reEnd := makeRegex(*beginFlag, *endFlag)

	// Find all sections
	sections, err := findSections(input, reBegin, reEnd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Replace all sections
	result, err := replaceSections(input, sections, *verbose, reBegin, reEnd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Output result to stdout
	if *stdout {
		fmt.Print(result)
		return
	}

	f, err := os.Create(*filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Write result to file
	w := bufio.NewWriter(f)
	w.WriteString(result)
	w.Flush()
}
