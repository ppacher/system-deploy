#!/bin/bash

set -e

hasUser() {
    local output=$(cat /etc/passwd | grep $1)
    if [[ "$output" == "" ]]; then
        return 1
    fi
    return 0
}

addUser() {
    local user=$1
    shift

    if hasUser ${user} ; then
        echo "User ${user} already created";
        return 0
    fi

    useradd -m -U -s /bin/zsh -G $@ ${user} || exit 1
    passwd -d ${user} || exit 1
}

groupadd -f admin
groupadd -f docker

addUser $1 admin,sudo,docker,adm 

exit 0