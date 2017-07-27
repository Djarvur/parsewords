# parsewords image:https://godoc.org/github.com/Djarvur/parsewords?status.svg["GoDoc",link="http://godoc.org/github.com/Djarvur/parsewords"] image:https://travis-ci.org/Djarvur/parsewords.svg["Build Status",link="https://travis-ci.org/Djarvur/parsewords"] image:https://coveralls.io/repos/Djarvur/parsewords/badge.svg?branch=master&service=github["Coverage Status",link="https://coveralls.io/github/Djarvur/parsewords?branch=master"]

Golang package based on [CPAN Text::ParseWords](http://search.cpan.org/~chorny/Text-ParseWords-3.30/lib/Text/ParseWords.pm) module.

Go regexps are little bit less powerfull than Perl,
so parser is little bit more complicated inside.

All the tests supplied with `Text::ParseWords` are implemented and passed.

## Benchmark

```
curl https://tools.ietf.org/rfc/rfc3501.txt > bench/bench.txt

$wc bench/bench.txt
    6051   28059  227639 test.txt

$perl bench/bench.pl < bench/bench.txt
100 iterations done in 9.986882s: 6604 words found in 227639 bytes of input

$go run bench/bench.go  < bench/bench.txt
100 iterations done in 6.572298304s: 6604 words found in 227639 bytes of input
```

Interesting: Go version is significantly faster, even Go regexps are known to be slower than Perl.