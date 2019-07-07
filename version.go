package gaedispemu

//go:generate bash -c "cat VERSION | xargs -n1 printf 'package gaedispemu\n\nconst Version = \"%s\"\n' > version_gen.go"
