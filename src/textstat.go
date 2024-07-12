package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
)

// Struct to hold text statistics
type TextStats struct {
	WordCount             int
	LetterCount           int
	SentenceCount         int
	ParagraphCount        int
	AverageWordLength     float64
	AverageSentenceLength float64
	LongestWord           string
	MostCommonWord        string
	UniqueWordCount       int
	FleschKincaidGrade    float64
	GunningFogIndex       float64
	SMOGGrade             float64
	EnglishLevel          string
	SMOGInterpretation    string
	FogInterpretation     string
}

// Function to count syllables in a word
func countSyllables(word string) int {
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

// Function to calculate text statistics
func calculateStats(text string) TextStats {
	words := strings.Fields(text)
	wordCount := len(words)
	letterCount := len(regexp.MustCompile(`\S`).FindAllString(text, -1))
	sentenceCount := len(regexp.MustCompile(`[.!?]`).FindAllString(text, -1))
	paragraphCount := len(strings.Split(text, "\n\n"))

	var totalWordLength, totalSyllables, complexWordCount int
	longestWord := ""
	wordFrequency := make(map[string]int)

	for _, word := range words {
		totalWordLength += len(word)
		if len(word) > len(longestWord) {
			longestWord = word
		}
		syllables := countSyllables(word)
		totalSyllables += syllables
		if syllables >= 3 {
			complexWordCount++
		}
		wordFrequency[strings.ToLower(word)]++
	}

	averageWordLength := float64(totalWordLength) / float64(wordCount)
	averageSentenceLength := float64(wordCount) / float64(sentenceCount)

	mostCommonWord := ""
	maxFrequency := 0
	for word, frequency := range wordFrequency {
		if frequency > maxFrequency {
			mostCommonWord = word
			maxFrequency = frequency
		}
	}

	uniqueWordCount := len(wordFrequency)

	fleschKincaidGrade := 0.39*averageSentenceLength + 11.8*float64(totalSyllables)/float64(wordCount) - 15.59
	gunningFogIndex := 0.4 * (averageSentenceLength + 100*float64(complexWordCount)/float64(wordCount))
	smogGrade := 1.043*math.Sqrt(float64(complexWordCount)*(30.0/float64(sentenceCount))) + 3.1291

	var englishLevel, smogInterpretation, fogInterpretation string

	// Determining English level based on Flesch-Kincaid Grade Level
	if fleschKincaidGrade <= 5 {
		englishLevel = "Basic"
	} else if fleschKincaidGrade <= 8 {
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

	return TextStats{
		WordCount:             wordCount,
		LetterCount:           letterCount,
		SentenceCount:         sentenceCount,
		ParagraphCount:        paragraphCount,
		AverageWordLength:     averageWordLength,
		AverageSentenceLength: averageSentenceLength,
		LongestWord:           longestWord,
		MostCommonWord:        mostCommonWord,
		UniqueWordCount:       uniqueWordCount,
		FleschKincaidGrade:    fleschKincaidGrade,
		GunningFogIndex:       gunningFogIndex,
		SMOGGrade:             smogGrade,
		EnglishLevel:          englishLevel,
		SMOGInterpretation:    smogInterpretation,
		FogInterpretation:     fogInterpretation,
	}
}

// Function to print text statistics
func printStats(stats TextStats) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Metric\tValue\tInterpretation\t")
	fmt.Fprintf(w, "Word count\t%d\t\n", stats.WordCount)
	fmt.Fprintf(w, "Letter count\t%d\t\n", stats.LetterCount)
	fmt.Fprintf(w, "Sentence count\t%d\t\n", stats.SentenceCount)
	fmt.Fprintf(w, "Paragraph count\t%d\t\n", stats.ParagraphCount)
	fmt.Fprintf(w, "Average word length\t%.2f\t\n", stats.AverageWordLength)
	fmt.Fprintf(w, "Average sentence length\t%.2f words\t\n", stats.AverageSentenceLength)
	fmt.Fprintf(w, "Longest word\t%s\t\n", stats.LongestWord)
	fmt.Fprintf(w, "Most common word\t%s\t\n", stats.MostCommonWord)
	fmt.Fprintf(w, "Unique word count\t%d\t\n", stats.UniqueWordCount)
	fmt.Fprintf(w, "Flesch-Kincaid Grade Level\t%.2f\t%s\t\n", stats.FleschKincaidGrade, stats.EnglishLevel)
	fmt.Fprintf(w, "Gunning Fog Index\t%.2f\t%s\t\n", stats.GunningFogIndex, stats.FogInterpretation)
	fmt.Fprintf(w, "SMOG Grade\t%.2f\t%s\t\n", stats.SMOGGrade, stats.SMOGInterpretation)
	w.Flush()
}

func main() {
	var filePath string
	flag.StringVar(&filePath, "file", "", "Path to the input file")
	flag.Parse()

	var text string
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

	stats := calculateStats(text)
	printStats(stats)
}
