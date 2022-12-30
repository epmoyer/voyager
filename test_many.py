#!/usr/bin/env python3
import os
# from rich import print as rprint
from rich.console import Console
from rich.theme import Theme

# --------------------
# Rich output console
# --------------------
# fmt: off
THEME = Theme({
    "case": "#c0c0c0",
    "path": "#808080",
})
# fmt: on
CONSOLE = Console(highlight=False, color_system='256', theme=THEME)
rprint = CONSOLE.print

TEST_CASES = [
    {
        'name': 'Normal',
        'path': r'/usr/local/bin',
    },
    {
        'name': 'With Context',
        'path': r'/usr/local/bin',
        # 'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/voyager/untracked/git_test_cases/normal',
        'username': 'eric',
    },
    {
        'name': 'With Context (as root)',
        'path': r'/usr/local/bin',
        # 'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/voyager/untracked/git_test_cases/normal',
        'username': 'root',
    },
    {
        'name': 'Normal, git repo',
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/voyager/untracked/git_test_cases/normal',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/voyager/untracked/git_test_cases/normal/subdir1',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/voyager/untracked/git_test_cases/normal/subdir1/subdir2',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/voyager/untracked/git_test_cases/detached',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/voyager/untracked/git_test_cases/unstarted',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/voyager/untracked/git_test_cases/untracked',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/voyager/untracked/git_test_cases/edited',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/voyager/untracked/git_test_cases/staged',
    },
    {
        'path': r'/Users/eric/Dropbox\ \(Personal\)/cab_dbx/code/go/voyager/untracked/git_test_cases/renamed',
    },
]


def main():
    # print("hello")
    print()
    for test_case in TEST_CASES:
        name = test_case.get('name', '(unnamed test case)')
        rprint(f'[case]{name}[/case]')
        path = test_case["path"]
        rprint(f'   [path]{path}[/path]')
        options = ''
        username = test_case.get('username')
        if username:
            options = f'--username={username}'
        stream = os.popen(
            f'"/Users/eric/Dropbox (Personal)/cab_dbx/code/go/voyager/voyager" --printable --powerline {options} {path}'
        )
        print(f'   {stream.read()}')
        stream = os.popen(
            f'"/Users/eric/Dropbox (Personal)/cab_dbx/code/go/voyager/voyager" --printable {options} {path}'
        )
        print(f'   {stream.read()}')
        # stream = os.popen(
        #     f'"/Users/eric/Dropbox (Personal)/cab_dbx/code/go/voyager/voyager" --dump {options} {path}'
        # )
        # print(f'   {stream.read()}')
    print()


if __name__ == "__main__":
    main()
