#!/usr/bin/env python3
import os

TEST_PATHS = [
    "/Users/eric/temp/prompt_examples/normal",
    "/Users/eric/temp/prompt_examples/normal/subdir1",
    "/Users/eric/temp/prompt_examples/normal/subdir1/subdir2",
    "/Users/eric/temp/prompt_examples/detatched",
    "/Users/eric/temp/prompt_examples/unstarted",
    "/Users/eric/temp/prompt_examples/untracked",
    "/Users/eric/temp/prompt_examples/edited",
    "/Users/eric/temp/prompt_examples/staged",
]

def main():
    # print("hello")
    print()
    for path in TEST_PATHS:
        print(f'   {path}')
        stream = os.popen(f'"/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/prompt_maker" --printable --powerline {path}')
        print(f'      {stream.read()}')
        stream = os.popen(f'"/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/prompt_maker" --printable {path}')
        print(f'      {stream.read()}')
    print()

if __name__ == "__main__":
    main()