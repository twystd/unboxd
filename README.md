![build](https://github.com/twystd/boxd/workflows/build/badge.svg)

# boxd
Go CLI for managing files and templates in Box: 

- list-templates
- create-template
- delete-template
- get-template
- list-files
- delete-file

### Raison d'être

Mostly just because I needed a Go interface to the Box API for another project and it turned out to be
convenient to create it as a separate library and then the CLI turned out to be useful. So it by no means
supersedes the official Box API implementations - it just implements some of them in a way that was 
useful for a particular requirement.

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

## References

## Attribution

