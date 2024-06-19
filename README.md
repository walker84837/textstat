# textstat: text file analyser

TextStat is a command-line tool for analyzing text files to provide various
readability statistics and metrics.

This project provides a comprehensive set of features to analyze text, including
word count, sentence count, paragraph count, average word length, average
sentence length, the longest word, the most common word, unique word count, and
readability scores like Flesch-Kincaid Grade Level, Gunning Fog Index, and SMOG
Grade.

Development is active, and new features are being added. Feel free to contribute
or suggest improvements!

## Table of Contents

  - [Usage](#usage)
  - [Support](#support)
  - [Roadmap](#roadmap)
  - [License](#license)

## Usage

Ensure you have Go installed. Clone the repository and run the following
command:

``` console
$ make
```

To analyze a text file, use the following command:

``` console
$ textstat -file path/to/yourfile.txt
```

Alternatively, you can pipe text input directly:

``` console
$ echo "Your text here" | textstat
```

TextStat accepts plain text files. Ensure that the text files are properly
formatted with clear sentences and paragraphs to get accurate statistics. Avoid
using files with complex formatting like HTML, Word documents, or PDFs without
converting them to plain text first.

## Roadmap

  - [ ] Add support for extracting and analyzing text from Word documents (.docx).
  - [ ] Integrate PDF text extraction for analysis.
  - [ ] Implement more readability metrics, such as Coleman-Liau Index and ARI.
  - [ ] Enhance the word frequency analysis to exclude common stopwords.
  - [ ] Add options for outputting results in different formats (e.g., JSON, CSV).

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
Before contributing, please follow these guidelines:

  - Follow the [code of conduct](CODE_OF_CONDUCT.md).
  - Keep a consistent coding style. To ensure your coding style remains the
    same, format your code with:
    ``` console
    $ go fmt path/to/source_code
    ```
  - Use the stable version of Go.
  - Prefer using the standard library over reinventing the wheel.

For support, please [open an
issue](https://github.com/walker84837/textstat/issues).

## License

This project is licensed under the [BSD-3-Clause License](LICENSE.md).
