package textstat

import (
	"math"
	"regexp"
	"strings"
	"unicode"
)

// Pre-compiled regex patterns for performance
var (
	// Patterns for syllable counting
	vowelGroupRegex      = regexp.MustCompile(`[aeiouy]+`)
	diphthongRegex       = regexp.MustCompile(`[aeiou]{2}`)
	tripthongRegex       = regexp.MustCompile(`[aeiou]{3}`)
	leadingTrailingRegex = regexp.MustCompile(`^[^aeiouy]+|[^aeiouy]+$`)

	// Patterns for text analysis
	sentenceEndRegex = regexp.MustCompile(`[.!?]+`)
	wordRegex        = regexp.MustCompile(`\b\w+\b`)
	nonLetterRegex   = regexp.MustCompile(`[^a-zA-Z]`)
)

// TextStats holds all calculated text statistics
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

// Analyzer provides text analysis functionality
type Analyzer struct{}

// NewAnalyzer creates a new text analyzer instance
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// countSyllables counts the number of syllables in a word using an improved algorithm
func (a *Analyzer) countSyllables(word string) int {
	if word == "" {
		return 0
	}

	word = strings.ToLower(strings.TrimSpace(word))

	// Remove non-alphabetic characters
	word = nonLetterRegex.ReplaceAllString(word, "")
	if word == "" {
		return 0
	}

	// Special cases for common English words
	switch word {
	case "a", "i":
		return 1
	case "":
		return 0
	}

	// Remove leading and trailing non-vowels
	word = leadingTrailingRegex.ReplaceAllString(word, "")

	// Replace tripthongs and diphthongs with single vowels
	word = tripthongRegex.ReplaceAllString(word, "a")
	word = diphthongRegex.ReplaceAllString(word, "a")

	// Count vowel groups
	syllables := vowelGroupRegex.FindAllString(word, -1)
	syllableCount := len(syllables)

	// Ensure at least one syllable for non-empty words
	if syllableCount == 0 && len(word) > 0 {
		syllableCount = 1
	}

	return syllableCount
}

// sanitizeInput validates and sanitizes input text
func (a *Analyzer) sanitizeInput(text string) string {
	if text == "" {
		return ""
	}

	// Remove null bytes and other control characters except newlines and tabs
	var result strings.Builder
	result.Grow(len(text))

	for _, r := range text {
		if r == '\n' || r == '\t' || !unicode.IsControl(r) {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// extractWords safely extracts words from text using pre-compiled regex
func (a *Analyzer) extractWords(text string) []string {
	if text == "" {
		return []string{}
	}

	matches := wordRegex.FindAllString(text, -1)
	if matches == nil {
		return []string{}
	}

	return matches
}

// extractSentences safely extracts sentences from text
func (a *Analyzer) extractSentences(text string) []string {
	if text == "" {
		return []string{}
	}

	// Split on sentence endings, keeping the delimiters
	sentences := sentenceEndRegex.Split(text, -1)
	if len(sentences) == 0 {
		return []string{}
	}

	// Filter out empty sentences and trim whitespace
	var result []string
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			result = append(result, sentence)
		}
	}

	return result
}

// extractParagraphs safely extracts paragraphs from text
func (a *Analyzer) extractParagraphs(text string) []string {
	if text == "" {
		return []string{}
	}

	paragraphs := strings.Split(text, "\n\n")
	if len(paragraphs) == 1 && paragraphs[0] == text {
		// No double newlines found, try single newlines
		paragraphs = strings.Split(text, "\n")
	}

	var result []string
	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph != "" {
			result = append(result, paragraph)
		}
	}

	if len(result) == 0 {
		return []string{text}
	}

	return result
}

// CalculateStats calculates comprehensive text statistics with improved security and performance
func (a *Analyzer) CalculateStats(text string) *TextStats {
	// Sanitize input
	text = a.sanitizeInput(text)
	if text == "" {
		return a.createEmptyStats()
	}

	// Extract text components
	words := a.extractWords(text)
	sentences := a.extractSentences(text)
	paragraphs := a.extractParagraphs(text)

	// Basic counts
	wordCount := len(words)
	sentenceCount := len(sentences)
	paragraphCount := len(paragraphs)

	if wordCount == 0 {
		return a.createEmptyStats()
	}

	// Calculate detailed statistics
	var totalWordLength, totalSyllables, complexWordCount int
	longestWord := ""
	wordFrequency := make(map[string]int)

	for _, word := range words {
		// Count letter characters only (excluding punctuation)
		cleanWord := strings.Trim(word, ".,!?:;'\"()[]{}")
		if cleanWord == "" {
			continue
		}

		totalWordLength += len(cleanWord)

		if len(cleanWord) > len(longestWord) {
			longestWord = cleanWord
		}

		syllables := a.countSyllables(cleanWord)
		totalSyllables += syllables

		if syllables >= 3 {
			complexWordCount++
		}

		wordFrequency[strings.ToLower(cleanWord)]++
	}

	// Calculate averages
	averageWordLength := float64(totalWordLength) / float64(wordCount)
	averageSentenceLength := float64(wordCount) / float64(sentenceCount)

	// Find most common word
	mostCommonWord, _ := a.findMostCommonWord(wordFrequency)

	uniqueWordCount := len(wordFrequency)

	// Calculate readability scores with safe division
	fleschKincaidGrade := a.calculateFleschKincaid(averageSentenceLength, float64(totalSyllables), float64(wordCount))
	gunningFogIndex := a.calculateGunningFog(averageSentenceLength, float64(complexWordCount), float64(wordCount))
	smogGrade := a.calculateSMOG(float64(complexWordCount), float64(sentenceCount))

	// Generate interpretations
	englishLevel := a.determineEnglishLevel(fleschKincaidGrade)
	smogInterpretation := a.interpretSMOG(smogGrade)
	fogInterpretation := a.interpretFog(gunningFogIndex)

	return &TextStats{
		WordCount:             wordCount,
		LetterCount:           totalWordLength,
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

// createEmptyStats returns an empty TextStats struct
func (a *Analyzer) createEmptyStats() *TextStats {
	return &TextStats{
		WordCount:             0,
		LetterCount:           0,
		SentenceCount:         0,
		ParagraphCount:        0,
		AverageWordLength:     0,
		AverageSentenceLength: 0,
		LongestWord:           "",
		MostCommonWord:        "",
		UniqueWordCount:       0,
		FleschKincaidGrade:    0,
		GunningFogIndex:       0,
		SMOGGrade:             0,
		EnglishLevel:          "",
		SMOGInterpretation:    "",
		FogInterpretation:     "",
	}
}

// findMostCommonWord finds the most common word and its frequency
func (a *Analyzer) findMostCommonWord(wordFrequency map[string]int) (string, int) {
	var mostCommonWord string
	maxFrequency := 0

	for word, frequency := range wordFrequency {
		if frequency > maxFrequency {
			mostCommonWord = word
			maxFrequency = frequency
		}
	}

	return mostCommonWord, maxFrequency
}

// calculateFleschKincaid calculates the Flesch-Kincaid Grade Level
func (a *Analyzer) calculateFleschKincaid(avgSentenceLength, avgSyllablesPerWord, wordCount float64) float64 {
	if wordCount == 0 {
		return 0
	}
	avgSyllablesPerWord = avgSyllablesPerWord / wordCount
	return 0.39*avgSentenceLength + 11.8*avgSyllablesPerWord - 15.59
}

// calculateGunningFog calculates the Gunning Fog Index
func (a *Analyzer) calculateGunningFog(avgSentenceLength, complexWordCount, wordCount float64) float64 {
	if wordCount == 0 {
		return 0
	}
	return 0.4 * (avgSentenceLength + 100*complexWordCount/wordCount)
}

// calculateSMOG calculates the SMOG Grade
func (a *Analyzer) calculateSMOG(complexWordCount, sentenceCount float64) float64 {
	if sentenceCount == 0 {
		return 0
	}
	return 1.043*math.Sqrt(complexWordCount*(30.0/sentenceCount)) + 3.1291
}

// determineEnglishLevel determines English level based on Flesch-Kincaid Grade
func (a *Analyzer) determineEnglishLevel(fleschKincaidGrade float64) string {
	switch {
	case fleschKincaidGrade <= 5:
		return "Basic"
	case fleschKincaidGrade <= 8:
		return "Intermediate"
	default:
		return "Advanced"
	}
}

// interpretSMOG provides interpretation of SMOG Grade
func (a *Analyzer) interpretSMOG(smogGrade float64) string {
	switch {
	case smogGrade <= 6:
		return "Basic English, easily understood by a wide audience, including children and those with basic reading skills."
	case smogGrade <= 9:
		return "Intermediate English, suitable for a general audience, including young adults and the average reader."
	case smogGrade <= 12:
		return "Upper Intermediate to Advanced English, suitable for high school students and adults with good reading skills."
	case smogGrade <= 16:
		return "Advanced English, suitable for college students and readers with strong comprehension skills."
	default:
		return "Very Advanced English, suitable for readers with higher education or specialized knowledge."
	}
}

// interpretFog provides interpretation of Gunning Fog Index
func (a *Analyzer) interpretFog(gunningFogIndex float64) string {
	switch {
	case gunningFogIndex <= 8:
		return "Basic English, easily understood by children and those with basic reading skills."
	case gunningFogIndex <= 12:
		return "Intermediate English, suitable for high school students."
	case gunningFogIndex <= 16:
		return "Upper Intermediate to Advanced English, suitable for college students."
	default:
		return "Very Advanced English, suitable for postgraduate students and professionals."
	}
}
