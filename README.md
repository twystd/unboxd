![build](https://github.com/twystd/boxd/workflows/build/badge.svg)

# boxd

Somewhat eclectic Go CLI for managing files and templates in Box: 

- list-files
- delete-file
- list-templates
- create-template
- get-template
- delete-template

### Raison d'être

Mostly just because another project needed a Go interface to the Box API and it turned out to be convenient
to create it as a separate library and then the CLI turned out to be occasionally useful. So it by no means
even vaguely approximates the official Box API implementations - it just implements some functionality in a 
way that was useful for a particular requirement.

*IN DEVELOPMENT*

## Releases

| *Version* | *Description*               |
| --------- | ----------------------------|
|           |                             |
|           |                             |                                                                    

## Installation

### Building from source

Assuming you have `Go v1.17+` and `make` installed:

```
git clone https://github.com/twystd/boxd.git
cd boxd
make build
```

If you prefer not to use `make`:
```
git clone https://github.com/twystd/boxd.git
cd boxd
go build -o bin/ ./...
```

#### Dependencies

| *Module*                                   | *Version*  |
| -------------------------------------------| ---------- |
|                                            |            |

## Notes

1. https://github.com/golang/go/issues/8860


## References

## Attribution

