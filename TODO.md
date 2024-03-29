# TODO

## IN PROGRESS

- [ ] Implement checkpointable pipeline that can be serialized and resumed
      - [ ] Don't recurse into folders that can't match the glob
      - [ ] Checkpoint on SIGHUP
      - (?) Checkpoint on CTRL-C
      - [ ] SIGINFO
      - [ ] Backoff and retry on HTTP error
      - [ ] Store list-folders to sqlite3 DB
      - [ ] Store list-files to sqlite3 DB
      - [ ] Use cached file/folder lists for queries

- [ ] update-template

- [ ] Restructure so that Box is just a wrapper around the API and the complexity devolves on
      e.g. the command implementation.
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
      - (?) Cache tokens to disk
            - With encryption (? GPG)
            - --no-cache option

- [ ] OAuth2
- [ ] App auth
- [ ] List folders by ID/name
- [ ] Templates for output
- [ ] (?) Photo gallery
      https://github.com/anvaka/panzoom

## NOTES

1. https://github.com/youmark/pkcs8
2. https://github.com/smallstep/crypto/blob/v0.9.2/pemutil/pkcs8.go#L189
3. http://keepachangelog.com/en/1.0.0
