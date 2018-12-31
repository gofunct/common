package templates

import "text/template"

func TravisTemplate() *template.Template {
	return MustCreateTemplate("travis", `language: go

go: '1.11'

env:
  global:
  - DEP_RELEASE_TAG=v0.5.0
  - FILE_TO_DEPLOY="dist/*"

  # GITHUB_TOKEN
  # TODO: shold encrypt and set a github access token using "travis encrypt" command
  - secure: "..."
  - REVIEWDOG_GITHUB_API_TOKEN=$GITHUB_TOKEN

cache:
  directories:
  - $GOPATH/pkg/dep
  - $HOME/.cache/go-build

jobs:
  include:
  - name: lint
    install: make setup
    script: make lint
    if: type = 'pull_request'

  - &test
    install: make setup
    script: make test
    if: type != 'pull_request'

  - <<: *test
    go: master

  - <<: *test
    go: '1.10'

  - stage: deploy
    install: make setup
    script: make packages -j4
    deploy:
    - provider: releases
      skip_cleanup: true
      api_key: $GITHUB_TOKEN
      file_glob: true
      file: $FILE_TO_DEPLOY
      on:
        tags: true
    if: type != 'pull_request'
`)
}
