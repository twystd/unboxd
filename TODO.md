# TODO

## IN PROGRESS

- [ ] Tags
      - [x] list
      - [x] add
      - [x] delete
      - [x] update

- [x] Upload file
      - [x] With folder ID
      - [ ] With folder name
      - [ ] Display file ID,name
      - [ ] (?) Byte streaming for uploading large files

- [ ] Walk folder tree
      - [ ] (MAYBE) Provide file and folder ID lookup functions and externalise the list/upload logic

- [x] Move file funcs to `files` package
      - [x] Cleanify function names
      - [ ] (MAYBE) Reinstate FileID type so that maps are typed
            - (OR) at minimum make file ID `uint64`
            - (OR) just return list of File
            - (OR) generalised BoxID

- [x] Move folders funcs to `folders` package
      - [x] Cleanify function names
      - [ ] (MAYBE) Reinstate FolderID type so that maps are typed
            - (OR) at minimum make folder ID `uint64`
            - (OR) just return list of Folder

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
    
- [ ] Oauth2
- [ ] App auth    
- [ ] List folders by ID/name
- [ ] Templates for output
- [ ] Include CHANGELOG in CLI
      - https://bhupesh-v.github.io/why-how-add-changelog-in-your-next-cli/
      - http://keepachangelog.com/en/1.0.0


## NOTES

1. https://github.com/youmark/pkcs8
2. https://github.com/smallstep/crypto/blob/v0.9.2/pemutil/pkcs8.go#L189
