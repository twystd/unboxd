# TODO

## IN PROGRESS

- [ ] help
- [ ] help command
- [x] version
- [ ] Restructure commands with embedded struct

- [ ] Implement checkpointable pipeline that can be serialized and resumed
      - [x] Resume from last checkpoint
      - [x] Use logger for reporting errors, progress etc
      - [x] Store list-folders to TSV
      - [x] Fix path concatenation
      - [x] --tags
      - [ ] --checkpoint-file
      - [ ] --no-resume
      - [ ] Checkpoint by base file ID

      - [ ] Dedupe folders list
      - [ ] Configurable interval between requests
      - [ ] Backoff and retry on HTTP error
      - [ ] Store list-folders to DB
      - [ ] Store list-files to TSV
      - [ ] Store list-files to DB

- [ ] Restructure so that Box is just a wrapper around the API and the complexity devolves on
      e.g. the command implementation.
      - [x] list-folders
      - [ ] list-files
            - Glob.HasPrefix or somesuch
      - [ ] Make IDs uint64
      - [ ] return error if strconv.ParseUint fails for ID


- [ ] glob
      - [ ] Rework to rather match on tokenised strings/DFA
      - [ ] `/alpha/**/today`

- [ ] list-folders
      - [ ] (?) should return 0 folder
      - (?) by folder ID

- [ ] list-files
      - (?) by file ID
      - [ ] just return list of File
      - [ ] List files in root dir

- [x] Upload file
      - [ ] Using folder name
      - [ ] (?) Byte streaming for uploading large files

- [x] Move file funcs to `files` package
      - [ ] (MAYBE) Reinstate FileID type so that maps are typed

- [x] Move folders funcs to `folders` package
- [x] File tags
- [x] Make public
- [x] Move template funcs to `templates` package
- [x] Replace FileID type with string
- [x] Authenticate with JWT credentials
- [x] Github workflow
- [x] `version` command
- [x] Move `Credentials` to `box` package

## TODO
- [ ] JWT auth
      - [ ] Marshal/unmarshal unit tests
      - [ ] Token refresh
      - [x] Authenticate
      - (?) Cache tokens to disk
            - With encryption (? GPG)
            - --no-cache option

- [ ] OAuth2
- [ ] App auth
- [ ] List folders by ID/name
- [ ] Templates for output
- [ ] Include CHANGELOG in CLI
      - https://bhupesh-v.github.io/why-how-add-changelog-in-your-next-cli/
      - http://keepachangelog.com/en/1.0.0

- [ ] (?) Photo gallery
      https://github.com/anvaka/panzoom

## NOTES

1. https://github.com/youmark/pkcs8
2. https://github.com/smallstep/crypto/blob/v0.9.2/pemutil/pkcs8.go#L189
