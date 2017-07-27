# parsewords

Golang package based on CPAN Text::ParseWords module

Go regexps are little bit less powerfull than Perl,
so parser is little bit more complicated inside.

Al the tests supplied with `Text::ParseWords` areimplemented and passed.

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