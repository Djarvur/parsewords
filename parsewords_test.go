package parsewords_test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/Djarvur/parsewords"
)

func getCaller(stackBack int) string {
	_, file, line, ok := runtime.Caller(stackBack + 1)
	if !ok {
		return "UNKNOWN"
	}

	if li := strings.LastIndex(file, "/"); li > 0 {
		file = file[li+1:]
	}

	return fmt.Sprintf("%s:%d", file, line)
}

func is(t *testing.T, got string, expected string) {
	if got != expected {
		t.Errorf("%s: Expected %q, got %q", getCaller(1), expected, got)
	}
}

func isInt(t *testing.T, got int, expected int) {
	if got != expected {
		t.Errorf("%s: Expected %v, got %v", getCaller(1), expected, got)
	}
}

func TestShellwords(t *testing.T) {
	words, err := parsewords.Shellwords(`foo "bar quiz" zoo`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, words[0], `foo`)
	is(t, words[1], `bar quiz`)
	is(t, words[2], `zoo`)
}

func TestQuotewords(t *testing.T) {
	// Test quotewords() with other parameters and null last field
	words, err := parsewords.Quotewords(`:+`, parsewords.KeepQuotes, `foo:::"bar:foo":zoo zoo:`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(words, ";"), `foo;"bar:foo";zoo zoo;`)
}

func TestQuotewordsDelimiters(t *testing.T) {
	// Test $keep eq 'delimiters' and last field zero
	words, err := parsewords.Quotewords(`\s+`, parsewords.KeepDelimiters, `4 3 2 1 0`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(words, ";"), `4; ;3; ;2; ;1; ;0`)
}

// Big ol' nasty test (thanks, Joerk!)
var BigOlNastyTest = `aaaa"bbbbb" cc\ cc \\\"dddd" eee\\\"ffff" "gg"`

func TestParseLineEscapedKeepQuotes(t *testing.T) {
	str := BigOlNastyTest

	// First with $keep == 1
	words, err := parsewords.ParseLine(`\s+`, parsewords.KeepQuotes, str)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(words, `|`), `aaaa"bbbbb"|cc\ cc|\\\"dddd" eee\\\"ffff"|"gg"`)
}

func TestParseLineEscapedKeepNothing(t *testing.T) {
	str := BigOlNastyTest

	// Now, $keep == 0
	words, err := parsewords.ParseLine(`\s+`, parsewords.KeepNothing, str)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(words, `|`), `aaaabbbbb|cc cc|\"dddd eee\"ffff|gg`)
}

func TestParseLineSinglequote(t *testing.T) {
	// Now test single quote behavior
	// Note: original test behavior is unclear for me so this one is modified to make sense.
	str := `aaaa"bbbbb" cc\ cc \\\"dddd\\' eee\\\"ffff\' gg`

	words, err := parsewords.ParseLine(`\s+`, parsewords.KeepNothing, str)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(words, `|`), `aaaabbbbb|cc cc|\"dddd\ eee\\\"ffff\|gg`)
}

func TestNestedQuotewords(t *testing.T) {
	// Make sure @nested_quotewords does the right thing
	words, err := parsewords.NestedQuotewords(`\s+`, parsewords.KeepNothing, `a b c`, `1 2 3`, `x y z`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	isInt(t, len(words), 3)
	isInt(t, len(words[0]), 3)
	isInt(t, len(words[1]), 3)
	isInt(t, len(words[2]), 3)
}

func TestShellwordsError(t *testing.T) {
	// Now test error return
	str := `foo bar baz"bach blech boop`

	words, err := parsewords.Shellwords(str)
	if err == nil {
		t.Errorf("Error expected but not received")
	}
	if len(words) > 0 {
		t.Errorf("Empty result expected but %d elements received", len(words))
	}

	words, err = parsewords.ParseLine(`s+`, parsewords.KeepNothing, str)
	if err == nil {
		t.Errorf("Error expected but not received")
	}
	if len(words) > 0 {
		t.Errorf("Empty result expected but %d elements received", len(words))
	}

	words, err = parsewords.Quotewords(`s+`, parsewords.KeepNothing, str)
	if err == nil {
		t.Errorf("Error expected but not received")
	}
	if len(words) > 0 {
		t.Errorf("Empty result expected but %d elements received", len(words))
	}

	lines, err := parsewords.NestedQuotewords(`s+`, parsewords.KeepNothing, str)
	if err == nil {
		t.Errorf("Error expected but not received")
	}
	if len(lines) > 0 {
		t.Errorf("Empty result expected but %d elements received", len(words))
	}
}

func TestEmptyFields(t *testing.T) {
	// Now test empty fields
	result, err := parsewords.ParseLine(":", parsewords.KeepNothing, `foo::0:"":::`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(result, "|"), `foo||0||||`)
}

func TestQuotedZero(t *testing.T) {
	// Test for 0 in quotes without $keep
	result, err := parsewords.ParseLine(":", parsewords.KeepNothing, `:"0":`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(result, "|"), `|0|`)
}

func TestQuotedOne(t *testing.T) {
	// Test for \001 in quoted string
	result, err := parsewords.ParseLine(":", parsewords.KeepNothing, `:"`+"\001"+`":`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(result, "|"), `|`+"\001"+`|`)
}

func TestPerlishSingleQuote(t *testing.T) {
	t.Skip("skipping unclear test.")
	// Now test perlish single quote behavior
	str := `aaaa"bbbbb" cc\ cc \\\\\"dddd\' eee\\\\\"\\\'ffff\' gg`
	result, err := parsewords.Quotewords(" ", parsewords.KeepNothing, str)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(result, "|"), `aaaabbbbb|cc cc|\"dddd eee\\\\"\'ffff|gg`)
}

func TestWhitespaceDelimiter(t *testing.T) {
	// test whitespace in the delimiters
	result, err := parsewords.Quotewords(" ", parsewords.KeepQuotes, `4 3 2 1 0`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(result, ";"), `4;3;2;1;0`)
}

func TestNewlineInsideQuotes(t *testing.T) {
	// [perl #30442] Text::ParseWords does not handle backslashed newline inside quoted text
	str := `"field1"	"field2` + "\n" + `still field2"	"field3"`

	result, err := parsewords.ParseLine("\t", parsewords.KeepQuotes, str)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(result, "|"), `"field1"|"field2`+"\n"+`still field2"|"field3"`)

	result, err = parsewords.ParseLine("\t", parsewords.KeepNothing, str)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(result, "|"), `field1|field2`+"\n"+`still field2|field3`)
}

func TestUnicode(t *testing.T) {
	// unicode
	str := "field1\u1234field2\\\u1234still field2\u1234field3"

	result, err := parsewords.ParseLine("\u1234", parsewords.KeepNothing, str)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(result, "|"), "field1|field2\u1234still field2|field3")
}

// Not relevant for Go
// # missing quote after matching regex used to hang after change #22997
// "1234" =~ /(1)(2)(3)(4)/;
// $string = qq{"missing quote};
// $result = join('|', shellwords($string));
// is($result, "");
//

func TestShellwordsStrip(t *testing.T) {
	// make sure shellwords strips out leading whitespace and trailng undefs
	// from parse_line, so it's behavior is more like /bin/sh
	result, err := parsewords.Shellwords(` aa \  \ bb `, ` \  `, `cc dd ee\ `)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(result, "|"), "aa| | bb| |cc|dd|ee ")

	result, err = parsewords.Shellwords(` aa \  \ bb `, ` \  `, `cc dd ee\  `)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(result, "|"), "aa| | bb| |cc|dd|ee ")

	result, err = parsewords.Shellwords(` aa \  \ bb `, ` \  `, `cc dd ee `)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	is(t, strings.Join(result, "|"), "aa| | bb| |cc|dd|ee")
}

// Not relevant for Go
// $SIG{ALRM} = sub {die "Timeout!"};
// alarm(3);
// @words = Text::ParseWords::old_shellwords("foo\\");
// is(@words, 1);
// alarm(0);
//
