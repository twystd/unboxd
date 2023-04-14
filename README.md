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

Mostly because another project needed a Go interface to the Box [Content API](https://developer.box.com/reference/)
and it turned out to be convenient to create it as a separate library and then the CLI turned out to be occasionally
useful. So it by no means even vaguely approximates the official Box API implementations - it just implements some
functionality in a way that was useful for a particular requirement.

## Releases

*EARLY DEVELOPMENT*

| *Version* | *Description*               |
| --------- | ----------------------------|
|           |                             |
|           |                             |

## Installation

### Building from source

Assuming you have `Go v1.20+` and `make` installed:

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


## unboxd

Usage: ```unboxd [options] <command> <arguments>```

General_ commands:

- `help`
- `version`

Folder commands:
- [`list-folders`](#list-folders)

File commands:
- [`list-files`](#list-files)
- [`upload-file`](#upload-file)
- [`delete-file`](#delete-file)
- [`tag-file`](#tag-file)
- [`untag-file`](#untag-file)
- [`retag-file`](#retag-file)

Template commands:
- [`list-templates`](#list-templates)
- [`get-template`](#get-template)
- [`create-template`](#create-template)
- [`delete-template`](#delete-template)

### General

#### `help`

Displays the usage information and a list of available commands. Command specific help displays the detailed usage for that command.

```
unboxd help

  Examples:

  unboxd help
  uhboxd help list-folders
```

#### `version`

Displays the current application version.

```
unboxd version

  Example:

  unboxd version
```

### Folder commands

The folder commands wrap the Box _Folder_ API:
```
unboxd list-folders
```

#### `list-folders`

Retrieves a list of folders matching the (optionally) supplied path

```
unboxd [options] list-folders [path]

  Options:
  --credentials <file> Sets the file containing the Box API credentials
  --debug              Displays verbose debugging information

  Example:

  unboxd --debug --credentials .credentials  list-folders /

  123456789 /unboxd
  987654321 /unboxd/photos
  876543219 /unboxd/docs
  765432198 /unboxd/docs/public
  …

```


### File commands

The file commands wrap the Box _File_ API:
```
unboxd list-files
unboxd upload-file
unboxd delete-file
unboxd tag-file
unboxd untag-file
unboxd retag-file
```


### Template commands

The file commands wrap the Box _Template_ API:
```
unboxd list-templates
unboxd get-template
unboxd create-template
unboxd delete-template
```


## Notes

1. https://github.com/golang/go/issues/8860


## References
