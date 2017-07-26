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

func qq(words ...string) string {
	return strings.Join(words, " ")
}

func TestShellwords(t *testing.T) {
	words, err := parsewords.Shellwords(qq(`foo`, `"bar quiz"`, `zoo`))
	if err != nil {
		t.Fatalf("Unexpected eror: %v", err)
	}

	is(t, words[0], `foo`)
	is(t, words[1], `bar quiz`)
	is(t, words[2], `zoo`)
}

func TestQuotewords(t *testing.T) {
	// Test quotewords() with other parameters and null last field
	words, err := parsewords.Quotewords(`:+`, parsewords.KeepQuotes, `foo:::"bar:foo":zoo zoo:`)
	if err != nil {
		t.Fatalf("Unexpected eror: %v", err)
	}

	is(t, strings.Join(words, ";"), `foo;"bar:foo";zoo zoo`)
}

func TestQuotewordsDelimiters(t *testing.T) {
	// Test $keep eq 'delimiters' and last field zero
	words, err := parsewords.Quotewords(`\s+`, parsewords.KeepDelimiters, `4 3 2 1 0`)
	if err != nil {
		t.Fatalf("Unexpected eror: %v", err)
	}

	is(t, strings.Join(words, ";"), qq(`4;`, `;3;`, `;2;`, `;1;`, `;0`))
}

// Big ol' nasty test (thanks, Joerk!)
var BigOlNastyTest = `aaaa"bbbbb" cc\ cc \\\"dddd" eee\\\"ffff" "gg"`

func TestParseLineEscapedKeepQuotes(t *testing.T) {
	str := BigOlNastyTest

	// First with $keep == 1
	words, err := parsewords.ParseLine(`\s+`, parsewords.KeepQuotes, str)
	if err != nil {
		t.Fatalf("Unexpected eror: %v", err)
	}

	is(t, strings.Join(words, `|`), `aaaa"bbbbb"|cc\ cc|\\\"dddd" eee\\\"ffff"|"gg"`)
}

func TestParseLineEscapedKeepNothing(t *testing.T) {
	str := BigOlNastyTest

	// Now, $keep == 0
	words, err := parsewords.ParseLine(`\s+`, parsewords.KeepNothing, str)
	if err != nil {
		t.Fatalf("Unexpected eror: %v", err)
	}

	is(t, strings.Join(words, `|`), `aaaabbbbb|cc cc|\"dddd eee\"ffff|gg`)
}

func TestParseLineSinglequote(t *testing.T) {
	// Now test single quote behavior
	// Note: original test behavior is unclear for me so this one is modified to make sense.
	str := `aaaa"bbbbb" cc\ cc \\\"dddd\\' eee\\\"ffff\' gg`

	words, err := parsewords.ParseLine(`\s+`, parsewords.KeepNothing, str)
	if err != nil {
		t.Fatalf("Unexpected eror: %v", err)
	}

	is(t, strings.Join(words, `|`), `aaaabbbbb|cc cc|\"dddd\ eee\\\"ffff\|gg`)
}

func TestNestedQuotewords(t *testing.T) {
	// Make sure @nested_quotewords does the right thing
	words, err := parsewords.NestedQuotewords(`\s+`, parsewords.KeepNothing, `a b c`, `1 2 3`, `x y z`)
	if err != nil {
		t.Fatalf("Unexpected eror: %v", err)
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
}

// @words = parse_line('s+', 0, $string);
// is(@words, 0);
//
// @words = quotewords('s+', 0, $string);
// is(@words, 0);
//
// {
//   # Gonna get some more undefined things back
//   no warnings 'uninitialized' ;
//
//   @words = nested_quotewords('s+', 0, $string);
//   is(@words, 0);
//
//   # Now test empty fields
//   $result = join('|', parse_line(':', 0, 'foo::0:"":::'));
//   is($result, 'foo||0||||');
//
//   # Test for 0 in quotes without $keep
//   $result = join('|', parse_line(':', 0, ':"0":'));
//   is($result, '|0|');
//
//   # Test for \001 in quoted string
//   $result = join('|', parse_line(':', 0, ':"' . "\001" . '":'));
//   is($result, "|\1|");
//
// }
//
// # Now test perlish single quote behavior
// $Text::ParseWords::PERL_SINGLE_QUOTE = 1;
// $string = 'aaaa"bbbbb" cc\ cc \\\\\"dddd\' eee\\\\\"\\\'ffff\' gg';
// $result = join('|', parse_line('\s+', 0, $string));
// is($result, 'aaaabbbbb|cc cc|\"dddd eee\\\\"\'ffff|gg');
//
// # test whitespace in the delimiters
// @words = quotewords(' ', 1, '4 3 2 1 0');
// is(join(";", @words), qq(4;3;2;1;0));
//
// # [perl #30442] Text::ParseWords does not handle backslashed newline inside quoted text
// $string = qq{"field1"	"field2\\\nstill field2"	"field3"};
//
// $result = join('|', parse_line("\t", 1, $string));
// is($result, qq{"field1"|"field2\\\nstill field2"|"field3"});
//
// $result = join('|', parse_line("\t", 0, $string));
// is($result, "field1|field2\nstill field2|field3");
//
// SKIP: { # unicode
//   skip "No unicode",1 if $]<5.008;
//   $string = qq{"field1"\x{1234}"field2\\\x{1234}still field2"\x{1234}"field3"};
//   $result = join('|', parse_line("\x{1234}", 0, $string));
//   is($result, "field1|field2\x{1234}still field2|field3",'Unicode');
// }
//
// # missing quote after matching regex used to hang after change #22997
// "1234" =~ /(1)(2)(3)(4)/;
// $string = qq{"missing quote};
// $result = join('|', shellwords($string));
// is($result, "");
//
// # make sure shellwords strips out leading whitespace and trailng undefs
// # from parse_line, so it's behavior is more like /bin/sh
// $result = join('|', shellwords(" aa \\  \\ bb ", " \\  ", "cc dd ee\\ "));
// is($result, "aa| | bb| |cc|dd|ee ");
//
// $SIG{ALRM} = sub {die "Timeout!"};
// alarm(3);
// @words = Text::ParseWords::old_shellwords("foo\\");
// is(@words, 1);
// alarm(0);
//
