package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
)

func main() {
	var filePath, text string
	flag.StringVar(&filePath, "file", "", "Path to the input file")
	flag.Parse()

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text += scanner.Text() + "\n"
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text += scanner.Text() + "\n"
		}
	}

	words := strings.Fields(text)
	wordCount := len(words)
	letterCount := len(regexp.MustCompile(`\S`).FindAllString(text, -1))

	// Sentence count
	sentenceCount := len(regexp.MustCompile(`[.!?]`).FindAllString(text, -1))

	// Paragraph count
	paragraphCount := len(strings.Split(text, "\n\n"))

	// Average word length
	var totalWordLength int
	for _, word := range words {
		totalWordLength += len(word)
	}
	averageWordLength := float64(totalWordLength) / float64(wordCount)

	// Average sentence length
	averageSentenceLength := float64(wordCount) / float64(sentenceCount)

	// Longest word
	longestWord := ""
	for _, word := range words {
		if len(word) > len(longestWord) {
			longestWord = word
		}
	}

	// Most common word
	wordFrequency := make(map[string]int)
	for _, word := range words {
		wordFrequency[strings.ToLower(word)]++
	}
	mostCommonWord := ""
	maxFrequency := 0
	for word, frequency := range wordFrequency {
		if frequency > maxFrequency {
			mostCommonWord = word
			maxFrequency = frequency
		}
	}

	// Unique word count
	uniqueWordCount := len(wordFrequency)

	// Syllable count function
	countSyllables := func(word string) int {
		word = strings.ToLower(word)
		vowelRegex := regexp.MustCompile(`[aeiouy]+`)
		diphthongRegex := regexp.MustCompile(`[aeiou]{2}`)
		tripthongRegex := regexp.MustCompile(`[aeiou]{3}`)
		leadingTrailingRegex := regexp.MustCompile(`^[^aeiouy]+|[^aeiouy]+$`)

		// Removing leading and trailing non-vowels
		word = leadingTrailingRegex.ReplaceAllString(word, "")

		// Replacing tripthongs with single vowels
		word = tripthongRegex.ReplaceAllString(word, "a")
		// Replacing diphthongs with single vowels
		word = diphthongRegex.ReplaceAllString(word, "a")

		// Counting vowel groups
		syllables := vowelRegex.FindAllString(word, -1)
		return len(syllables)
	}

	// Total syllable count and complex word count
	totalSyllables := 0
	complexWordCount := 0
	for _, word := range words {
		syllables := countSyllables(word)
		totalSyllables += syllables
		if syllables >= 3 {
			complexWordCount++
		}
	}

	fleschKincaidGradeLevel := 0.39*averageSentenceLength + 11.8*float64(totalSyllables)/float64(wordCount) - 15.59
	gunningFogIndex := 0.4 * (averageSentenceLength + 100*float64(complexWordCount)/float64(wordCount))
	smogGrade := 1.043*math.Sqrt(float64(complexWordCount)*(30.0/float64(sentenceCount))) + 3.1291

	var englishLevel, smogInterpretation, fogInterpretation string

	// Determining English level based on Flesch-Kincaid Grade Level
	if fleschKincaidGradeLevel <= 5 {
		englishLevel = "Basic"
	} else if fleschKincaidGradeLevel <= 8 {
		englishLevel = "Intermediate"
	} else {
		englishLevel = "Advanced"
	}

	// Interpretation of SMOG Grade
	switch {
	case smogGrade <= 6:
		smogInterpretation = "Basic English, easily understood by a wide audience, including children and those with basic reading skills."
	case smogGrade <= 9:
		smogInterpretation = "Intermediate English, suitable for a general audience, including young adults and the average reader."
	case smogGrade <= 12:
		smogInterpretation = "Upper Intermediate to Advanced English, suitable for high school students and adults with good reading skills."
	case smogGrade <= 16:
		smogInterpretation = "Advanced English, suitable for college students and readers with strong comprehension skills."
	default:
		smogInterpretation = "Very Advanced English, suitable for readers with higher education or specialized knowledge."
	}

	// Interpretation of Gunning Fog Index
	switch {
	case gunningFogIndex <= 8:
		fogInterpretation = "Basic English, easily understood by children and those with basic reading skills."
	case gunningFogIndex <= 12:
		fogInterpretation = "Intermediate English, suitable for high school students."
	case gunningFogIndex <= 16:
		fogInterpretation = "Upper Intermediate to Advanced English, suitable for college students."
	default:
		fogInterpretation = "Very Advanced English, suitable for postgraduate students and professionals."
	}

	fmt.Printf("Word count: %d\n", wordCount)
	fmt.Printf("Letter count: %d\n", letterCount)
	fmt.Printf("Sentence count: %d\n", sentenceCount)
	fmt.Printf("Paragraph count: %d\n", paragraphCount)
	fmt.Printf("Average word length: %.2f\n", averageWordLength)
	fmt.Printf("Average sentence length: %.2f words\n", averageSentenceLength)
	fmt.Printf("Longest word: %s\n", longestWord)
	fmt.Printf("Most common word: %s\n", mostCommonWord)
	fmt.Printf("Unique word count: %d\n", uniqueWordCount)
	fmt.Printf("Flesch-Kincaid Grade Level and interpretation: %.2f - %s\n", fleschKincaidGradeLevel, englishLevel)
	fmt.Printf("Gunning Fog index and interpretation: %.2f - %s\n", gunningFogIndex, fogInterpretation)
	fmt.Printf("SMOG grade and interpretation: %.2f - %s\n", smogGrade, smogInterpretation)
}
