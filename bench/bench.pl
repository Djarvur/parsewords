use Text::ParseWords;
use Time::HiRes qw(time);

my $iterations = 100;                      

my $text = do { local $/; <STDIN> };

my @words = ();

my $started = time();
for (my $ri = 0; $ri < $iterations; $ri++) {
  @words = shellwords($text);
}
my $spent = time() - $started;

printf(
  "%d iterations done in %fs: %d words found in %d bytes of input\n",
  $iterations,
  $spent,
  scalar(@words), 
  length($text),
)