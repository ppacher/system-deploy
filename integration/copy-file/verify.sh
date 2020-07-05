#!/bin/bash
set -e

[[ -e /tmp/${TEST}/dir/testfile ]] || exit 1
[[ "$(cat /tmp/${TEST}/dir/testfile)" == "some content" ]] || exit 1