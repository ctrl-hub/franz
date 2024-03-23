#!/usr/bin/env bash

set -o nounset
set -o pipefail

test -z "${DEBUG:-}" || {
    set -x
}

TAG="${1}"

# remove tag
git tag --delete "${TAG}"
git push origin ":refs/tags/${TAG}"

# republish tag
git tag "${1}"
git push origin "${1}"
