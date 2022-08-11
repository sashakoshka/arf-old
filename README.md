<br>

# ![ARF](assets/logo.svg)

**This is being rewritten.**

[Click here to go to the new repository](https://git.tebibyte.media/sashakoshka/arf).

All the dangerous fun of C, and you can do it in style.

Arf is an extremely modular, minimalist programming language. Its syntax and
structure are designed to make writing a program feel like constructing a
machine out of physical parts.

It uses syntactical indentation and everything must be indented with 8 spaces
(the holy indentation size used in the Linux kernel).

What does arf stand for you ask? It doesn't, its just the sound that dogs make.
So far, there have been no successful backgronymming attempts. Ciao has also
been suggested as a name. 

This does not compile code yet, it just partially parses it. I don't even have
a language spec. Basically, don't expect a lot right now. Maybe expect a lot
some time later though.

Here is some example code:

```
:arf
module main
author "Sasha Koshka"
require "io"
---

func rr main
        > argc:Int
        > argv:{String}
        < status:Int
        ---
        io.println "Hello, world!"
```

Isn't that neat? Its like a document.

# Roadmap

These are things, in order, that have already been done or are planned:

- [x] Lexer
- [x] Parser
- [x] AST (I think this is done)
- [ ] Semantic analysis
- [ ] Transpiling to C
- [ ] Compile any file with options by passing command line arguments
- [ ] Rewrite compiler in Arf (this current one is a mess...)
- [ ] Convert Arf module to Arf header file (so we can have shared libraries)
- [ ] Convert C header file to Arf header file
- [ ] Modular Arf standard library (written in Arf)
  - [ ] I/O
  - [ ] String manipulation
  - [ ] Memory management
  - [ ] Networking
  - [ ] Multithreading

These aren't necessary for making Arf a full-fledged language, but they would
make it more useful:

- [ ] Produce LLVM IR instead of C
- [ ] Option to use C stdlib for memory allocation backend, for compatibility
      with libraries written in C
- [ ] Rewrite some coreutils
- [ ] Make sure it works with essential libraries
  - [ ] XCB
  - [ ] Wayland
  - [ ] OpenGL
  - [ ] Cairo
  - [ ] GTK
  - [ ] OpenSSL (or similar)
