#!/bin/zsh

echo "================================================================================="
TARGET_PATH="/Users/eric/Dropbox (Personal)/cab_dbx/code/go/voyager/test1/test2"
# echo "DUMP:"
# go run . --dump --powerline $TARGET_PATH
# echo
echo "zsh, PowerLine: format=prompt"
go run . --format=prompt --powerline $TARGET_PATH
echo
echo "zsh, PowerLine: format=prompt  (RENDERED)"
print -P $(go run . --powerline $TARGET_PATH)
echo
echo "zsh, PowerLine: format=display"
go run . --format=display --powerline $TARGET_PATH
echo
echo
echo "zsh, Text: format=prompt"
go run . -format=prompt $TARGET_PATH
echo
echo "zsh, Text: format=prompt  (RENDERED)"
print -P $(go run . -format=prompt $TARGET_PATH)
echo
echo "zsh, Text: format=display)"
go run . -format=display $TARGET_PATH
echo
# echo "================================================================================="
echo "bash, Powerline: foramt=display"
go run . -format=display -shell=bash -powerline $TARGET_PATH
echo
echo
echo "================================================================================="