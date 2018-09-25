# vGo (Go Modules Support) Makefile

## Dependencies

The following dependencies must be installed in order to use the Makefile. Any
dependencies not listed here are automatically installed by the Makefile.

- GNU Make
- Bash and standard unix utilities such as sed, grep, cat, etc
- Go v1.9+

## Repository Conventions

This Makefile assumes that the repository follows the conventions below:

- Go source code is in the `src` folder
- Binaries to be built are in `src/cmd/<name>`
- [Go Modules](https://github.com/golang/go/wiki/Modules) is used for dependency management
