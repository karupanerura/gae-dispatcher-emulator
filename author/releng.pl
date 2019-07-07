use strict;
use warnings;
use feature qw/say/;

use version;
use Cwd qw/realpath/;
use File::Spec::Functions qw/catdir updir/;
use File::Basename qw/dirname/;

my $root_dir = realpath(catdir(dirname($0), updir()));
chdir $root_dir;

# check deps
{
    if (!$ENV{GITHUB_TOKEN}) {
        say STDERR "GITHUB_TOKEN env is required";
        exit 1;
    }
    for my $cmd (qw/go git gox ghr/) {
        my $has_err = system "type $cmd > /dev/null 2>&1";
        if ($has_err) {
            say STDERR "$cmd is required";
            exit 1;
        }
    }
}

print 'Next Version: ';
chomp(my $version = <STDIN>);

# check version
{
    if ($version !~ /^[0-9]+\.[0-9]+\.[0-9]+$/) {
        die "Invalid version format: $version";
    }

    my $current_version = do {
        open my $fh, '<', 'VERSION' or die $!;
        <$fh>;
    };
    chomp $current_version;

    if (version->parse($version) <= version->parse($current_version)) {
        die "New version $version is not newer than $current_version";
    }
}

# renew version
{
    open my $fh, '>', 'VERSION' or die $!;
    print $fh $version;
    close $fh or die $!;
}

# edit release note
while (1) {
    system $ENV{EDITOR} || 'vim', 'CHANGELOG.md';

    my $code = system 'git', 'diff', '--check', '--exit-code', 'CHANGELOG.md';
    last if $code != 0;
}

# generate & test
system 'go', 'generate';
system 'go', 'test';

# build
system 'gox', '-output', 'build/{{.Dir}}_{{.OS}}_{{.Arch}}', './...';

while (1) {
    print 'OK? (y/n): ';
    chomp(my $yn = <STDIN>);
    last   if $yn =~ /^y(?:es)?$/i;
    exit 1 if $yn =~ /^n(?:o)?$/i;
}

# release
system 'git', 'add', qw/VERSION version_gen.go CHANGELOG.md/;
system 'git', 'commit', '-m', "release $version";
system 'git', 'tag', "v$version";
system 'git', 'push';
system 'git', 'push', '--tags';
system 'ghr', "v$version", './build';
