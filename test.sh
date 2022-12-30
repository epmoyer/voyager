#!/bin/zsh

echo "================================================================================="
TARGET_PATH="/Users/eric/Dropbox (Personal)/cab_dbx/code/go/voyager/test1/test2"
# echo "DUMP:"
# go run . --dump --powerline $TARGET_PATH
# echo
echo "AS PROMPT STRING (ZSH, POWERLINE):"
go run . --powerline $TARGET_PATH
echo
echo
echo "RENDERED PROMPT STRING (ZSH, POWERLINE):"
print -P $(go run . --powerline $TARGET_PATH)
echo
echo "AS PRINTABLE (ZSH, POWERLINE):"
go run . --printable --powerline $TARGET_PATH
echo
echo
echo "AS PROMPT STRING (ZSH, TEXT):"
go run . $TARGET_PATH
echo
echo
echo "RENDERED PROMPT STRING (ZSH, POWERLINE):"
print -P $(go run . $TARGET_PATH)
echo
echo "AS PRINTABLE (TEXT):"
go run . --printable $TARGET_PATH
echo
# echo "================================================================================="
echo "AS PROMPT STRING (BASH, POWERLINE):"
go run . --shell=bash --powerline $TARGET_PATH
echo
echo
echo "================================================================================="