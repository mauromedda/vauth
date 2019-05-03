#!/bin/bash

VERSION=$(git describe --tags --exact-match)
REPO=$(basename $(pwd))
ARCHS="linux/amd64 darwin/amd64 windows/amd64"

set -e
if [[ -z "${VERSION}" ]]; then
  echo "No tag present, stopping build now."
  exit 0
fi

if [[ -z "${GITHUB_TOKEN}" ]]; then
  echo "Please set \$GITHUB_TOKEN environment variable"
  exit 1
fi

DEV_PLATFORM=${DEV_PLATFORM:-"./pkg/$(go env GOOS)_$(go env GOARCH)"}
for F in $(find ${DEV_PLATFORM} -mindepth 1 -maxdepth 1 -type f); do
    shasum -a 256 ${F}* > $(dirname ${F})/SHA256SUMS
    github-release upload --user mauromedda --repo ${REPO} --tag ${VERSION} --name ${F} --file ${F}
done

for F in $(find ${DEV_PLATFORM} -mindepth 1 -maxdepth 1 -type f); do
    github-release upload --user mauromedda --repo ${REPO} --tag ${VERSION} --name ${F} --file ${F}
done
