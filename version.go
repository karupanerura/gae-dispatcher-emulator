package gaedispemu

//go:generate bash -c "cat VERSION | xargs -n1 printf 'package gaedispemu\n\n// Version is a current version of the package.\nconst Version = \"%s\"\n' > version_gen.go"
