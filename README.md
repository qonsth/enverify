# Enverify

## Description
Devops tool for verifying matches of the `.env` files

For example, it's useful for multiple `docker-compose` environments

## Installation
``` sh
go install github.com/qonsth/enverify@v0
```

## Usage
``` sh
enverify .env .env.example .env.dev .env.prod
```

## Advanced usage
Make `enverify-local.sh`, which can be useful for local verification and recreation of .env on different branch switching
``` sh
#!/usr/bin/env bash

recreate() {
  flag='.'
  until [[ $flag == 'y' || $flag == 'n' || $flag == '' ]]; do
      printf "\x1B[93mRecreate %s file by %s.example? y/n (default: n)\x1B[0m\n" "$1"
      read -n1 -s -r flag
  done

  if [[ $flag == 'y' ]]; then
    [[ -e "$1" ]] && cp "$1" "$1.bck"
    cp "$1.example" "$1"
  else
    exit 1
  fi
}

enverify .env .env.example .env.dev .env.prod || recreate .env
```
