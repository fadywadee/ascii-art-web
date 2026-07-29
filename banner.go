package main

import (
	"embed"
	"fmt"
	"strings"
)

var bannerFS embed.FS

func loadBanner(filename string) ([]string, error) {
	data, err := bannerFS.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(content, "\n")
	return lines, nil
}

func getCharLines(lines []string, c rune) ([]string, error) {
	startLine := int((c-32))*9 + 1
	endLine := startLine + 8

	if startLine < 0 || endLine > len(lines) {
		return nil, fmt.Errorf("line out of bounds")
	}

	return lines[startLine:endLine], nil
}

func render(banner []string, input string) string {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "")

	if isOnlyNewlines(input) {
		return strings.Repeat("\n", len(input))
	}

	input = strings.ReplaceAll(input, "\\n", "\n")
	parts := strings.Split(input, "\n")
	var result strings.Builder

	for _, part := range parts {
		if part == "" {
			result.WriteString("\n")
			continue
		}

		for row := 0; row < 8; row++ {
			for _, ch := range part {
				charLines, err := getCharLines(banner, ch)
				if err != nil {
					continue
				}
				result.WriteString(charLines[row])
			}
			result.WriteString("\n")
		}
	}

	return result.String()
}

func validateInput(input string, lines []string) (bool, rune) {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "")

	maxChar := (len(lines) / 9) + 32

	for _, ch := range input {
		if ch == '\n' || ch == '\t' {
			continue
		}
		if ch < 32 || ch > rune(maxChar) {
			return false, ch
		}
	}
	return true, 0
}

func isOnlyNewlines(s string) bool {
	for _, ch := range s {
		if ch != '\n' {
			return false
		}
	}
	return true
}

// دالة وسيطة بتفصل الـ Processing عن الـ Web Handler
func ProcessAsciiArt(text, bannerName string) (string, int, error) {
	if text == "" {
		return "", 400, fmt.Errorf("text is required")
	}

	var bannerFile string
	switch bannerName {
	case "standard":
		bannerFile = "standard.txt"
	case "shadow":
		bannerFile = "shadow.txt"
	case "thinkertoy":
		bannerFile = "thinkertoy.txt"
	default:
		return "", 400, fmt.Errorf("invalid banner")
	}

	banner, err := loadBanner(bannerFile)
	if err != nil {
		return "", 404, fmt.Errorf("banner not found")
	}

	valid, badChar := validateInput(text, banner)
	if !valid {
		return "", 400, fmt.Errorf("unsupported character: %c", badChar)
	}

	result := render(banner, text)
	return result, 200, nil
}
