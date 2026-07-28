package main

import (
	"bufio"
	"os"
	"strings"
)

func loadBanner(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0, 855)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func getCharLines(lines []string, c rune) []string {
	// Calculate the start line index using the formula
	startLine := int((c-32))*9 + 1

	// Collect the next 8 lines starting from startLine
	endline := startLine + 7

	// Return them as a slice
	return lines[startLine : endline+1]
}

func renderLine(banner []string, text string) string {
	var result strings.Builder

	result.Grow(len(text) * 8 * 10)

	for row := 0; row <= 7; row++ {
		for _, ch := range text {
			charLines := getCharLines(banner, ch)
			result.WriteString(charLines[row])
		}
		result.WriteString("\n")
	}

	return result.String()
}

func render(banner []string, input string) string {
	// 1. إزالة \r الخاص بالمتصفح ونظام الويندوز
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "")

	// 2. تحويل السطر الجديد الحقيقي (\n) القادم من الـ textarea إلى "\\n" ليتقسم صح بـ strings.Split
	input = strings.ReplaceAll(input, "\n", "\\n")

	parts := strings.Split(input, "\\n")
	var result strings.Builder

	result.Grow(len(input) * 8 * 10)

	for i, part := range parts {
		if part == "" {
			if i < len(parts)-1 {
				result.WriteString("\n")
			}
			continue
		}

		for row := 0; row <= 7; row++ {
			for _, ch := range part {
				charLines := getCharLines(banner, ch)
				result.WriteString(charLines[row])
			}
			result.WriteString("\n")
		}
	}

	return result.String()
}

func validateInput(input string, lines []string) (bool, rune) {
	// 1. تنظيف \r قبل فحص الحروف
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "")

	maxChar := (len(lines) / 9) + 32

	for _, ch := range input {
		/*
			if ch == '\n' {
				continue
			}*/

		if ch < 32 || ch > rune(maxChar) {
			return false, ch
		}
	}
	return true, 0
}
