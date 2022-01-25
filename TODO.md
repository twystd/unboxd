# TODO

## IN PROGRESS

- [x] Move file funcs to `files` package
- [ ] Move folders funcs to `folders` package
- [ ] Upload file
- [ ] Make public

- [ ] Tags
      - [x] list
      - [x] add
      - [x] delete
      - [ ] update


- [ ] (MAYBE) reinstate FileID and FolderID types so that maps are typed

- [x] Replace FileID type with string
- [x] Authenticate with JWT credentials
- [x] Github workflow
- [x] `version` command
- [x] Move `Credentials` to `box` package

## TODO
- [ ] JWT auth
      - [ ] `Load` unit tests
      - [ ] Token refresh
      - [x] Authenticate
            - https://github.com/youmark/pkcs8
            - https://github.com/smallstep/crypto/blob/v0.9.2/pemutil/pkcs8.go#L189

- [ ] List folders by ID/name
- [ ] Add file tag
- [ ] Delete file tag
- [ ] Templates for output

- [ ] Include CHANGELOG in CLI
      - https://bhupesh-v.github.io/why-how-add-changelog-in-your-next-cli/
      - http://keepachangelog.com/en/1.0.0
