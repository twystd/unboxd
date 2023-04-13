# TODO

## IN PROGRESS

- [x] help
- [ ] help command
      - [x] list-folders
      - [x] list-files
      - [x] upload-file
      - [ ] delete-file
      - [ ] tag-file
      - [ ] untag-file
      - [ ] retag-file
      - [ ] list-templates
      - [ ] create-template
      - [ ] delete-template
      - [ ] get-template
      - [ ] version
      - [ ] help

- [x] version
- [x] Restructure commands with embedded struct

- [ ] Implement checkpointable pipeline that can be serialized and resumed
      - [x] Resume from last checkpoint
      - [x] Use logger for reporting errors, progress etc
      - [x] Store list-folders to TSV
      - [x] Fix path concatenation
      - [x] --tags
      - [x] --checkpoint-file
      - [x] --no-resume
      - [x] list-files
      - [x] --batch-size
      - [ ] Dedupe folders list
      - [ ] Dedupe files list
      - [ ] Store list-files to TSV
      - [ ] Include account ID + base file ID in checkpoint and verify on resume
      - [ ] Include command in checkpoint and verify on resume
      - [ ] Don't recurse into folders that can't match the glob
      - [ ] Move delay to end of loop
      - [ ] Checkpoint on SIGHUP
      - (?) Checkpoint on CTRL-C
      - [ ] SIGINFO

      - [ ] Configurable interval between requests
      - [ ] Backoff and retry on HTTP error
      - [ ] Store list-folders to DB
      - [ ] Store list-files to DB

- [ ] Restructure so that Box is just a wrapper around the API and the complexity devolves on
      e.g. the command implementation.
      - [x] list-folders
      - [ ] list-files
            - Glob.HasPrefix or somesuch
      - [ ] Make IDs strings
      - [ ] Return error if strconv.ParseUint fails for ID

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
- [x] Include CHANGELOG in CLI
- [ ] (?) Photo gallery
      https://github.com/anvaka/panzoom

## NOTES

1. https://github.com/youmark/pkcs8
2. https://github.com/smallstep/crypto/blob/v0.9.2/pemutil/pkcs8.go#L189
3. http://keepachangelog.com/en/1.0.0
