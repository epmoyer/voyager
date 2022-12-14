#!/bin/zsh
go run . --dump "/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/test1"
go run . --dump "/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/"
echo XXX
print -P $(go run . "/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker")
echo XXX
print -P $(go run . "/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/test1/test2")