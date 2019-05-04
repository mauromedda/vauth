#!/bin/bash -x

VERSION=$(git describe --tags --exact-match)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
REPO=$(basename $(pwd))
ARCHS="linux/amd64 darwin/amd64 windows/amd64"
BUILD_DIR=.build
set -e
if [[ -z "${VERSION}" ]] ; then
  echo "No tag present or you are not in the master branch, stopping build now."
  exit 0
fi

if [[ -z "${GITHUB_TOKEN}" ]]; then
  echo "Please set \$GITHUB_TOKEN environment variable"
  exit 1
fi

echo "Create release"
github-release release \
    --user mauromedda \
    --repo ${REPO} \
    --tag ${VERSION} \
    --pre-release

cd ${BUILD_DIR} && shasum -a 256 * > SHA256SUMS

for F in *; do
    github-release upload --user mauromedda --repo ${REPO} --tag ${VERSION} --name ${F} --file ${F}
done

