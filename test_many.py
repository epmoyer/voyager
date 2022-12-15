#!/usr/bin/env python3
import os

TEST_PATHS = [
    # "/Users/eric/temp/git_test_cases/normal",
    # r"/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker",
    r"/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/normal",
    r"/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/normal/subdir1",
    r"/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/normal/subdir1/subdir2",
    r"/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/detatched",
    r"/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/unstarted",
    r"/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/untracked",
    r"/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/edited",
    r"/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/staged",
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