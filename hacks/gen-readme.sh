#!/bin/bash

go build -o /tmp/deploy ./cmd/deploy 

BUILTIN=./pkg/actions/builtin/

function gendoc() {
    dir=$2

    if [[ "$dir" == "" ]]; then dir=$1; fi
    
    /tmp/deploy describe $1 --markdown > ${BUILTIN}/${dir}/$1.md ;
}

gendoc installpackages platform
gendoc platform
gendoc systemd
gendoc copy
gendoc exec
gendoc onchange
gendoc editfile