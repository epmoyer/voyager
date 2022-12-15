#!/bin/zsh

# echo "================================================================================="
# TARGET_PATH="/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/test1"
# echo "DUMP:"
# go run . --dump $TARGET_PATH
# echo
# echo "AS ZSH PROMPT STRING:"
# go run . $TARGET_PATH
# echo
# echo
# echo "RENDERED PROMPT STRING:"
# print -P $(go run . $TARGET_PATH)
# echo
# echo
# echo "AS PRINTABLE:"
# go run . --printable $TARGET_PATH
# echo
echo "================================================================================="
TARGET_PATH="/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/test1/test2"
echo "DUMP:"
go run . --dump --powerline $TARGET_PATH
echo
echo "AS ZSH PROMPT STRING (POWERLINE):"
go run . --powerline $TARGET_PATH
echo
echo
echo "RENDERED PROMPT STRING (POWERLINE):"
print -P $(go run . --powerline $TARGET_PATH)
echo
echo "AS PRINTABLE (POWERLINE):"
go run . --printable --powerline $TARGET_PATH
echo
echo
echo "AS ZSH PROMPT STRING (TEXT):"
go run . $TARGET_PATH
echo
echo
echo "RENDERED PROMPT STRING (POWERLINE):"
print -P $(go run . $TARGET_PATH)
echo
echo "AS PRINTABLE (TEXT):"
go run . --printable $TARGET_PATH
echo
echo "================================================================================="
# go run . --dump "/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/"
# echo XXX
# print -P $(go run . "/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker")
# echo XXX
# print -P $(go run . "/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/test1/test2")