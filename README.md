![build](https://github.com/twystd/unboxd/workflows/build/badge.svg)

# unboxd

A somewhat eclectic Go CLI for managing files and templates in [Box](box.com): 

- list-folders

- list-files
- upload-file
- delete-file
- tag-file
- untag-file

- list-templates  
- create-template
- get-template
- delete-template

Currently supports authentication and authorisation using either Box _client_ or _JWT_ credentials.

### Raison d'être

Mostly just because another project needed a Go interface to the Box API and it turned out to be convenient
to create it as a separate library and then the CLI turned out to be occasionally useful. So it by no means
even vaguely approximates the official Box API implementations - it just implements some functionality in a 
way that was useful for a particular requirement.

## Releases

*EARLY DEVELOPMENT*

| *Version* | *Description*               |
| --------- | ----------------------------|
|           |                             |
|           |                             |                                                                    

## Installation

### Building from source

Assuming you have `Go v1.17+` and `make` installed:

```
git clone https://github.com/twystd/unboxd.git
cd unboxd
make build
```

If you prefer not to use `make`:
```
git clone https://github.com/twystd/unboxd.git
cd unboxd
go build -o bin/ ./...
```

#### Dependencies

| *Module*                                             | *Version*  |
| -----------------------------------------------------| ---------- |
| [youmark:PKCS8](https://github.com/youmark/pkcs8)    | (latest)   |
| [cristalhq:JWT](https://github.com/cristalhq/jwt/v4) | v4         |


## Notes

1. https://github.com/golang/go/issues/8860


## References

