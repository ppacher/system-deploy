#!/bin/bash
set -e

[[ $(cat /etc/passwd | grep bar | grep zsh) != "" ]] || exit 1
[[ -d /home/bar/.oh-my-zsh ]] || exit 1