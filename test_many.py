#!/usr/bin/env python3
import os

TEST_CASES = [
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/normal',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/normal',
        'username': 'eric',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/normal',
        'username': 'root',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/normal/subdir1',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/normal/subdir1/subdir2',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/detached',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/unstarted',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/untracked',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/edited',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/staged',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/prompt_maker/untracked/git_test_cases/renamed',
    },
]


def main():
    # print("hello")
    print()
    for test_case in TEST_CASES:
        path = test_case["path"]
        print(f'   {path}')
        options = ''
        username = test_case.get('username')
        if username:
            options = f'--username={username}'
        stream = os.popen(
            f'"/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/prompt_maker" --printable --powerline {options} {path}'
        )
        print(f'      {stream.read()}')
        stream = os.popen(
            f'"/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/prompt_maker" --printable {options} {path}'
        )
        print(f'      {stream.read()}')
        # stream = os.popen(
        #     f'"/Users/eric/Dropbox (Personal)/cab_dbx/code/go/prompt_maker/prompt_maker" --dump {options} {path}'
        # )
        # print(f'      {stream.read()}')
    print()


if __name__ == "__main__":
    main()
