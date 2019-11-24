#!/usr/bin/env bash

set -o errexit
set -o nounset

cd "$(dirname "$(realpath "$0")")"/../..
go build ./cmd/url-parser

function clean() {
    rm -rf ".test-out"
    rm -rf ".test-err"
}

trap "clean" EXIT

urls=(
    "http://www.example.com"
    "postgres://user:abc{DEf1=ghi@example.com"
    "https://sub.domain.co.uk/path/file.html"
    "https://"
)

# component

./url-parser -c "host" > .test-out 2> .test-err <<< IFS="\n" "${urls[@]}"

expected=$'www.example.com\n\nsub.domain.co.uk'
result="$(cat .test-out)"
test "$expected" = "$result" || {
    printf >&2 $"Unexpected parser output:\n\nExpected:\n%s\n\nResult:\n%s\n" "$expected" "$result"
}

expectedErr=$'parse postgres://user:abc{DEf1=ghi@example.com: net/url: invalid userinfo'
resultErr="$(cat .test-err)"
test "$expectedErr" = "$resultErr" || {
    printf >&2 $"Unexpected parser output:\n\nExpected:\n%s\n\nResult:\n%s\n" "$expectedErr" "$resultErr"
}

# format

./url-parser -f "{scheme} - {tld} - {file}" > .test-out 2> .test-err <<< IFS="\n" "${urls[@]}"

expected=$'http - com - \n\nhttps - co.uk - file.html\nhttps -  - '
result="$(cat .test-out)"
test "$expected" = "$result" || {
    printf >&2 $"Unexpected parser output:\n\nExpected:\n%s\n\nResult:\n%s\n" "$expected" "$result"
}

expectedErr=$'parse postgres://user:abc{DEf1=ghi@example.com: net/url: invalid userinfo'
resultErr="$(cat .test-err)"
test "$expectedErr" = "$resultErr" || {
    printf >&2 $"Unexpected parser output:\n\nExpected:\n%s\n\nResult:\n%s\n" "$expectedErr" "$resultErr"
}
