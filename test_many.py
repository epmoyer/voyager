#!/usr/bin/env python3
import os
import subprocess

# from rich import print as rprint
from rich.console import Console
from rich.theme import Theme

# --------------------
# Rich output console
# --------------------
# fmt: off
THEME = Theme({
    "case": "#d0d0d0",
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
        'name': 'Error',
        'path': r'/usr/local/bin',
        'environment_vars': {'VGER_RETVAL': "1"}
    },
    {
        'name': 'With Context',
        'path': r'/usr/local/bin',
        'username': 'eric',
    },
    {
        'name': 'With Context (as root)',
        'path': r'/usr/local/bin',
        'username': 'root',
    },
    {
        'name': 'Git repo',
        'path': r'./untracked/git_test_cases/normal',
    },
    {
        'name': 'Git repo, in subdirectory',
        'path': r'./untracked/git_test_cases/normal/subdir1',
    },
    {
        'name': 'Git repo, in second subdirectory',
        'path': r'./untracked/git_test_cases/normal/subdir1/subdir2',
    },
    {
        'name': 'Git repo, w/ detached head',
        'path': r'./untracked/git_test_cases/detached',
    },
    {
        'name': 'Git repo, unstarted (new repo, no check-ins)',
        'path': r'./untracked/git_test_cases/unstarted',
    },
    {
        'name': 'Git repo, untracked change',
        'path': r'./untracked/git_test_cases/untracked',
    },
    {
        'name': 'Git repo, edited file',
        'path': r'./untracked/git_test_cases/edited',
    },
    {
        'name': 'Git repo, staged file',
        'path': r'./untracked/git_test_cases/staged',
    },
    {
        'name': 'Git repo, renamed file',
        'path': r'./untracked/git_test_cases/renamed',
    },
]


def main():
    # print("hello")
    print()
    for test_case in TEST_CASES:
        name = test_case.get('name', '(unnamed test case)')
        rprint(f'[case]{name}:[/case]')
        path = test_case["path"]
        rprint(f'   [path]{path}[/path]')

        command_line_args = [
            '/Users/eric/Dropbox (Personal)/cab_dbx/code/go/voyager/voyager', 
            '--printable', 
        ]
        username = test_case.get('username')
        if username:
            command_line_args.append(f'--username={username}')
        
        environment = os.environ.copy()
        environment_vars = test_case.get('environment_vars')
        if environment_vars:
            for key, value in environment_vars.items():
                environment[key] = value
        
        output = subprocess.check_output(
            command_line_args + ['--powerline', path],
            env=environment,
        )
        print(f'   {output.decode("utf-8")}')

        output = subprocess.check_output(
            command_line_args + [path],
            env=environment,
        )
        print(f'   {output.decode("utf-8")}')
    print()


if __name__ == "__main__":
    main()
