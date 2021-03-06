// Package parsewords - parse text into an array of tokens or array of arrays
package parsewords

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type keepType int

// Possible keep argument values
//
// The keep argument is a semi-boolean flag.  If not KeepNothing, then the tokens are
// split on the specified delimiter, but all other characters (including
// quotes and backslashes) are kept in the tokens.  If keep is KeepNothing then the
// Quotewords() functions remove all quotes and backslashes that are
// not themselves backslash-escaped or inside of single quotes (i.e.,
// Quotewords() tries to interpret these characters just like the Bourne
// shell).
//
// As an additional feature, keep may be KeepDelimiters value which
// causes the functions to preserve the delimiters in each string as
// tokens in the token lists, in addition to preserving quote and
// backslash characters.
const (
	KeepNothing    keepType = 0
	KeepQuotes     keepType = 1
	KeepDelimiters keepType = 2
)

// Errors might be returned.
var (
	ErrQuoteUnclosed = errors.New("quote unclosed")
)

var (
	unslashRegexp    = regexp.MustCompile(`(?s:\\(.))`)               // nolint: gochecknoglobals
	unslashReplace   = `$1`                                           // nolint: gochecknoglobals
	trimLeftRegexp   = regexp.MustCompile(`^\s+`)                     // nolint: gochecknoglobals
	trimLeftReplace  = ``                                             // nolint: gochecknoglobals
	trimRightRegexp  = regexp.MustCompile(`((?:\\\s)|(?:[^\\]))\s+$`) // nolint: gochecknoglobals
	trimRightReplace = `$1`                                           // nolint: gochecknoglobals
)

// Shellwords is written as a special case of &quotewords(), and it
// does token parsing with whitespace as a delimiter-- similar to most
// Unix shells.
func Shellwords(lines ...string) ([]string, error) {
	wholeLine := trimRightRegexp.ReplaceAllString(
		trimLeftRegexp.ReplaceAllString(
			strings.Join(lines, ""),
			trimLeftReplace,
		),
		trimRightReplace,
	)

	if len(wholeLine) == 0 {
		return nil, nil
	}

	allwords, err := ParseLine(`\s+`, KeepNothing, wholeLine)
	if err != nil {
		return nil, err
	}

	return allwords, nil
}

// ParseLine does tokenizing on a single string.
func ParseLine(delimiter string, keep keepType, line string) ([]string, error) {
	delimiting, err := regexp.Compile(delimiter)
	if err != nil {
		return nil, err
	}

	return ParseLinePrecompiled(delimiting, keep, line)
}

// ParseLinePrecompiled does tokenizing on a single string using provided regular expression as a delimiter.
func ParseLinePrecompiled(delimiting *regexp.Regexp, keep keepType, line string) ([]string, error) {
	words, err := smartSplit(delimiting, line)
	if len(words) == 0 || err != nil {
		return nil, err
	}

	pieces := make([]string, 0, len(words))

	for _, word := range words {
		piece := ""

		for _, token := range word.tokens {
			str := line[token.on:token.off]
			if keep == KeepNothing {
				str = unquote(str)
			}

			piece += str
		}

		pieces = append(pieces, piece)

		if keep == KeepDelimiters && word.delimiter.on > 0 {
			pieces = append(pieces, line[word.delimiter.on:word.delimiter.off])
		}
	}

	return pieces, nil
}

func unquote(str string) string {
	switch {
	case len(str) == 0:
		return ""
	case str[0] == '\'':
		if len(str) > 2 { // nolint: gomnd
			return str[1 : len(str)-1]
		}

		return ""
	case str[0] == '"':
		if len(str) > 2 { // nolint: gomnd
			return unslashRegexp.ReplaceAllString(str[1:len(str)-1], unslashReplace)
		}

		return ""
	}

	return unslashRegexp.ReplaceAllString(str, unslashReplace)
}

// Quotewords and NestedQuotewords functions accept a delimiter
// (which can be a regular expression)
// and a list of lines and then breaks those lines up into a list of
// words ignoring delimiters that appear inside quotes.  Quotewords()
// returns all of the tokens in a single long list, while NestedQuotewords()
// returns a list of token lists corresponding to the elements of lines[].
func Quotewords(delimiter string, keep keepType, lines ...string) ([]string, error) {
	delimiting, err := regexp.Compile(delimiter)
	if err != nil {
		return nil, err
	}

	return QuotewordsPrecompiled(delimiting, keep, lines...)
}

// QuotewordsPrecompiled and NestedQuotewordsPrecompiled functions accept a delimiter
// (which is a regular expression)
// and a list of lines and then breaks those lines up into a list of
// words ignoring delimiters that appear inside quotes.  Quotewords()
// returns all of the tokens in a single long list, while NestedQuotewords()
// returns a list of token lists corresponding to the elements of lines[].
func QuotewordsPrecompiled(delimiting *regexp.Regexp, keep keepType, lines ...string) ([]string, error) {
	allwords := make([]string, 0, len(lines))

	for _, line := range lines {
		words, err := ParseLinePrecompiled(delimiting, keep, line)
		if err != nil {
			return nil, err
		}

		if !(len(words) > 0 || len(line) == 0) {
			return nil, nil
		}

		allwords = append(allwords, words...)
	}

	return allwords, nil
}

// NestedQuotewords and Quotewords functions accept a delimiter
// (which can be a regular expression)
// and a list of lines and then breaks those lines up into a list of
// words ignoring delimiters that appear inside quotes.  Quotewords()
// returns all of the tokens in a single long list, while NestedQuotewords()
// returns a list of token lists corresponding to the elements of lines[].
func NestedQuotewords(delimiter string, keep keepType, lines ...string) ([][]string, error) {
	delimiting, err := regexp.Compile(delimiter)
	if err != nil {
		return nil, err
	}

	return NestedQuotewordsPrecompiled(delimiting, keep, lines...)
}

// NestedQuotewordsPrecompiled and QuotewordsPrecompiled functions accept a delimiter
// (which is a regular expression)
// and a list of lines and then breaks those lines up into a list of
// words ignoring delimiters that appear inside quotes.  Quotewords()
// returns all of the tokens in a single long list, while NestedQuotewords()
// returns a list of token lists corresponding to the elements of lines[].
func NestedQuotewordsPrecompiled(delimiting *regexp.Regexp, keep keepType, lines ...string) ([][]string, error) {
	allwords := make([][]string, 0, len(lines))

	for _, line := range lines {
		words, err := ParseLinePrecompiled(delimiting, keep, line)
		if err != nil {
			return nil, err
		}

		if len(words) > 0 && words[len(words)-1] == "" {
			words = words[:len(words)-1]
		}

		if !(len(words) > 0 || len(line) == 0) {
			return nil, nil
		}

		allwords = append(allwords, words)
	}

	return allwords, nil
}

type substring struct {
	on  int
	off int
}

type delimitedWord struct {
	tokens    []substring
	delimiter substring
}

func delimitedWordNew(token substring, delimiter substring) delimitedWord {
	return delimitedWord{
		tokens:    []substring{token},
		delimiter: delimiter,
	}
}

func smartSplit(delimiting *regexp.Regexp, line string) ([]delimitedWord, error) {
	quoted, slashed, err := enumerateQuotes(line)
	if err != nil {
		return nil, err
	}

	delimiters := enumerateDelimiters(line, quoted, slashed, delimiting)

	if len(delimiters) == 0 {
		return []delimitedWord{delimitedWordNew(substring{0, len(line)}, badSubstring)}, nil
	}

	return enumerateWords(line, delimiters, quoted), nil
}

func enumerateQuotes(line string) ([]substring, []int, error) { // nolint: gocyclo
	quoted := make([]substring, 0, 10)
	slashed := make([]int, 0, 10)

	sqOn := -1
	dqOn := -1
	slashOn := false

	for ri, sym := range line {
		if slashOn && sqOn < 0 {
			slashed = append(slashed, ri)
		}

		switch {
		case sym == '\\':
			slashOn = !slashOn
		case sym == '"' && dqOn < 0 && sqOn < 0 && !slashOn:
			dqOn = ri
		case sym == '"' && dqOn >= 0 && !slashOn:
			quoted = append(quoted, substring{dqOn, ri + 1})
			dqOn = -1
		case sym == '\'' && dqOn < 0 && sqOn < 0 && !slashOn:
			sqOn = ri
			slashOn = false
		case sym == '\'' && sqOn >= 0:
			quoted = append(quoted, substring{sqOn, ri + 1})
			sqOn = -1
			slashOn = false
		default:
			slashOn = false
		}
	}

	if sqOn >= 0 {
		return nil, nil, fmt.Errorf("single quote unclosed: %d: %w", sqOn, ErrQuoteUnclosed)
	}

	if dqOn >= 0 {
		return nil, nil, fmt.Errorf("double quote unclosed: %d: %w", sqOn, ErrQuoteUnclosed)
	}

	return quoted, slashed, nil
}

func enumerateWords(
	line string,
	delimiters []substring,
	quoted []substring,
) []delimitedWord {
	words := make([]delimitedWord, 0, len(delimiters))
	curPos := 0
	curQuoted := 0

	for _, delimiter := range delimiters {
		word := delimitedWord{make([]substring, 0, 1), badSubstring}

		for curPos < delimiter.on {
			if curQuoted < len(quoted) {
				if curPos == quoted[curQuoted].on {
					word.tokens = append(word.tokens, quoted[curQuoted])
					curPos = quoted[curQuoted].off
					curQuoted++

					continue
				}

				if delimiter.on > quoted[curQuoted].on {
					word.tokens = append(word.tokens, substring{curPos, quoted[curQuoted].on})
					curPos = quoted[curQuoted].on

					continue
				}
			}

			word.tokens = append(word.tokens, substring{curPos, delimiter.on})
			curPos = delimiter.off
		}

		curPos = delimiter.off
		word.delimiter = delimiter
		words = append(words, word)
	}

	switch {
	case curPos < len(line):
		words = append(words, delimitedWordNew(substring{curPos, len(line)}, badSubstring))
	case len(delimiters) > 0 && delimiters[len(delimiters)-1].off == len(line):
		words = append(words, delimitedWordNew(substring{len(line), len(line)}, badSubstring))
	}

	return words
}

func enumerateDelimiters(
	line string,
	quoted []substring,
	slashed []int,
	delimiting *regexp.Regexp,
) []substring {
	matches := delimiting.FindAllStringIndex(line, -1)
	delimiters := make([]substring, 0, len(matches))

	for _, match := range matches {
		delimeter := checkDelimiter(line, quoted, slashed, delimiting, substring{match[0], match[1]})
		if delimeter.on >= 0 {
			delimiters = append(delimiters, delimeter)
		}
	}

	return delimiters
}

var badSubstring = substring{on: -1, off: -1} // nolint: gochecknoglobals

func checkDelimiter(
	line string,
	quoted []substring,
	slashed []int,
	delimiting *regexp.Regexp,
	delimiter substring,
) substring {
	for _, quote := range quoted {
		if delimiter.on >= quote.on && delimiter.on < quote.off {
			return badSubstring
		}
	}

	for _, slash := range slashed {
		if delimiter.on == slash {
			delimiter.on++
			if delimiter.on < delimiter.off && delimiting.MatchString(line[delimiter.on:delimiter.off]) {
				return delimiter
			}

			return badSubstring
		}
	}

	return delimiter
}
