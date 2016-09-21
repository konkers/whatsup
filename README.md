# What's Up Doc(umentation) Generator

A simple C documentation generation based on libclang.

I wrote this because I was looking to auto generate some documentation for u
(https://github.com/konkers/u).  I wanted something that output HTML that
doesn't look like it was designed in the early 2000s, support C (instead of C++)
as a "first class citizen", and doesn't have a dizzying array of dependencies.
I didn't multiple language support.  I didn't need to output anything but HTML.

How complicated could it be.....

This is a **HUGE** work in progress.

## Getting go-clang to compile
```
CGO_LDFLAGS="-L`llvm-config --libdir`" go get -u github.com/go-clang/v3.9/...
```

## Building whatsup to run on OSX
```
go install -ldflags "-extldflags '-Xlinker -rpath -Xlinker `llvm-config --libdir`'" all
```

## TODO
- [ ] Macros
- [ ] Enums
- [ ] Links between index and details
- [ ] Type details
