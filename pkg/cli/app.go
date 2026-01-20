package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"textstat/internal/errors"
	"textstat/pkg/extractor"
	"textstat/pkg/textstat"
)

// Config holds CLI configuration
type Config struct {
	FilePath     string
	FileType     string
	OutputFormat string
	Verbose      bool
}

// App represents the CLI application
type App struct {
	analyzer *textstat.Analyzer
	config   *Config
}

// NewApp creates a new CLI application
func NewApp() *App {
	return &App{
		analyzer: textstat.NewAnalyzer(),
	}
}

// ParseFlags parses command line flags and returns configuration
func (a *App) ParseFlags() (*Config, error) {
	config := &Config{}

	flag.StringVar(&config.FilePath, "file", "", "Path to the input file")
	flag.StringVar(&config.FileType, "type", "text", "Type of the input file (text, pdf, docx)")
	flag.StringVar(&config.OutputFormat, "format", "table", "Output format (table, json)")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")

	flag.Parse()

	// Validate file type
	switch config.FileType {
	case "text", "pdf", "docx": // Valid types
	default:
		return nil, errors.NewValidationError("invalid file type. Supported types: text, pdf, docx")
	}

	// Validate output format
	switch config.OutputFormat {
	case "table", "json": // Valid formats
	default:
		return nil, errors.NewValidationError("invalid output format. Supported formats: table, json")
	}

	return config, nil
}

// readFromStdin reads text from standard input
func (a *App) readFromStdin() (string, error) {
	var builder strings.Builder
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		builder.WriteString(scanner.Text())
		builder.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", errors.NewFileError("failed to read from standard input", err)
	}

	return builder.String(), nil
}

// readFromFile reads text from a file
func (a *App) readFromFile(filePath string, fileType string) (string, error) {
	ext, err := extractor.NewDocumentExtractor(filePath, extractor.FileType(fileType))
	if err != nil {
		return "", err
	}

	text, err := ext.ExtractText()
	if err != nil {
		return "", err
	}

	return text, nil
}

// getInputText gets input text from either file or stdin
func (a *App) getInputText(config *Config) (string, error) {
	if config.FilePath != "" {
		return a.readFromFile(config.FilePath, config.FileType)
	}
	return a.readFromStdin()
}

// printTableStats prints statistics in table format
func (a *App) printTableStats(stats *textstat.TextStats) error {
	const (
		minwidth = 0
		tabwidth = 0
		padding  = 3
		padchar  = ' '
		flags    = 0
	)

	w := tabwriter.NewWriter(os.Stdout, minwidth, tabwidth, padding, padchar, flags)

	// Print header with tabs
	fmt.Fprintln(w, "METRIC\tVALUE\tINTERPRETATION\t")

	// Print statistics with proper tab termination
	fmt.Fprintf(w, "Word count\t%d\t-\t\n", stats.WordCount)
	fmt.Fprintf(w, "Letter count\t%d\t-\t\n", stats.LetterCount)
	fmt.Fprintf(w, "Sentence count\t%d\t-\t\n", stats.SentenceCount)
	fmt.Fprintf(w, "Paragraph count\t%d\t-\t\n", stats.ParagraphCount)
	fmt.Fprintf(w, "Average word length\t%.2f\tcharacters\t\n", stats.AverageWordLength)
	fmt.Fprintf(w, "Average sentence length\t%.2f\twords\t\n", stats.AverageSentenceLength)
	fmt.Fprintf(w, "Longest word\t%s\t-\t\n", stats.LongestWord)
	fmt.Fprintf(w, "Most common word\t%s\t-\t\n", stats.MostCommonWord)
	fmt.Fprintf(w, "Unique word count\t%d\t-\t\n", stats.UniqueWordCount)
	fmt.Fprintf(w, "Flesch-Kincaid Grade Level\t%.2f\t%s\t\n", stats.FleschKincaidGrade, stats.EnglishLevel)
	fmt.Fprintf(w, "Gunning Fog Index\t%.2f\t%s\t\n", stats.GunningFogIndex, stats.FogInterpretation)
	fmt.Fprintf(w, "SMOG Grade\t%.2f\t%s\t\n", stats.SMOGGrade, stats.SMOGInterpretation)

	return w.Flush()
}

// printJSONStats prints statistics in JSON format
func (a *App) printJSONStats(stats *textstat.TextStats) {
	fmt.Printf("{\n")
	fmt.Printf("  \"word_count\": %d,\n", stats.WordCount)
	fmt.Printf("  \"letter_count\": %d,\n", stats.LetterCount)
	fmt.Printf("  \"sentence_count\": %d,\n", stats.SentenceCount)
	fmt.Printf("  \"paragraph_count\": %d,\n", stats.ParagraphCount)
	fmt.Printf("  \"average_word_length\": %.2f,\n", stats.AverageWordLength)
	fmt.Printf("  \"average_sentence_length\": %.2f,\n", stats.AverageSentenceLength)
	fmt.Printf("  \"longest_word\": \"%s\",\n", stats.LongestWord)
	fmt.Printf("  \"most_common_word\": \"%s\",\n", stats.MostCommonWord)
	fmt.Printf("  \"unique_word_count\": %d,\n", stats.UniqueWordCount)
	fmt.Printf("  \"flesch_kincaid_grade\": %.2f,\n", stats.FleschKincaidGrade)
	fmt.Printf("  \"gunning_fog_index\": %.2f,\n", stats.GunningFogIndex)
	fmt.Printf("  \"smog_grade\": %.2f,\n", stats.SMOGGrade)
	fmt.Printf("  \"english_level\": \"%s\",\n", stats.EnglishLevel)
	fmt.Printf("  \"smog_interpretation\": \"%s\",\n", stats.SMOGInterpretation)
	fmt.Printf("  \"fog_interpretation\": \"%s\"\n", stats.FogInterpretation)
	fmt.Printf("}\n")
}

// printStats prints statistics in the specified format
func (a *App) printStats(stats *textstat.TextStats, format string) error {
	switch format {
	case "json":
		a.printJSONStats(stats)
		return nil
	case "table":
		fallthrough
	default:
		return a.printTableStats(stats)
	}
}

// Run executes the CLI application
func (a *App) Run() error {
	config, err := a.ParseFlags()
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	a.config = config

	// Get input text
	text, err := a.getInputText(config)
	if err != nil {
		return fmt.Errorf("failed to get input text: %w", err)
	}

	// Calculate statistics
	stats := a.analyzer.CalculateStats(text)

	// Print results
	if err := a.printStats(stats, config.OutputFormat); err != nil {
		return fmt.Errorf("failed to print stats: %w", err)
	}

	return nil
}

// HandleError handles application errors and prints appropriate messages
func (a *App) HandleError(err error) {
	if err == nil {
		return
	}

	if textStatErr, ok := err.(*errors.TextStatError); ok {
		switch textStatErr.Type {
		case errors.ValidationError:
			fmt.Fprintf(os.Stderr, "Error: %s\n", textStatErr.Message)
		case errors.FileError:
			fmt.Fprintf(os.Stderr, "File Error: %s\n", textStatErr.Message)
		case errors.ParsingError:
			fmt.Fprintf(os.Stderr, "Parsing Error: %s\n", textStatErr.Message)
		case errors.ExtractionError:
			fmt.Fprintf(os.Stderr, "Extraction Error: %s\n", textStatErr.Message)
		case errors.UnsupportedFormatError:
			fmt.Fprintf(os.Stderr, "Unsupported Format: %s\n", textStatErr.Message)
		default:
			fmt.Fprintf(os.Stderr, "Error: %s\n", textStatErr.Message)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Unexpected error: %s\n", err.Error())
	}

	os.Exit(1)
}
