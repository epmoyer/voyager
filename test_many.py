#!/usr/bin/env python3
import os
import subprocess
from pathlib import Path

# Library
import click
from rich.console import Console
from rich.theme import Theme

# --------------------
# Rich output console
# --------------------
# fmt: off
THEME = Theme({
    "case": "#d0d0d0",
    "path": "#808080",
    "shell": "#00ff00",
    # "shell": "#ffff00",
    "presentation": "#ff8000",
    "format": "#ff00ff",
    # "format": "#00ffff",
})
# fmt: on
CONSOLE = Console(highlight=False, color_system='256', theme=THEME)
rprint = CONSOLE.print  # rich print

# fmt:off
TEST_CASES = [
    {
        'name': 'Normal',
        'path': r'/usr/local/bin',
    },
    {
        'name': 'Normal, No path truncation',
        'path': r'/usr/local/bin',
        'environment_vars': {
            'VGER_TRUNCATION_START_DEPTH': "1000"
        }
    },
    {
        'name': 'Error',
        'path': r'/usr/local/bin',
        'environment_vars': {
            'VGER_RETVAL': "1"
        }
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
        'path': r'./git_test_cases/normal',
    },
    {
        'name': 'Git repo, in subdirectory',
        'path': r'./git_test_cases/normal/subdir1',
    },
    {
        'name': 'Git repo, in second subdirectory',
        'path': r'./git_test_cases/normal/subdir1/subdir2',
    },
    {
        'name': 'Git repo, w/ detached head',
        'path': r'./git_test_cases/detached',
    },
    {
        'name': 'Git repo, unstarted (new repo, no check-ins)',
        'path': r'./git_test_cases/unstarted',
    },
    {
        'name': 'Git repo, untracked change',
        'path': r'./git_test_cases/untracked',
    },
    {
        'name': 'Git repo, edited file',
        'path': r'./git_test_cases/edited',
    },
    {
        'name': 'Git repo, staged file',
        'path': r'./git_test_cases/staged',
    },
    {
        'name': 'Git repo, renamed file',
        'path': r'./git_test_cases/renamed',
    },
]
# fmt: on

@click.group()
def cli():
    pass

@cli.command()
@click.option('-c', '--colors', 'enable_color_modes', is_flag=True, help='Run test in all color modes')
def cases(enable_color_modes):
    if not enable_color_modes:
        run_tests()
        return
    for arg in ['--color=16m', '--color=256', '--color=16', '--no-color']:
        print('-' * 50 + ' ' + arg)
        run_tests([arg])

def run_tests(extra_args=None):
    print()
    for test_case in TEST_CASES:
        name = test_case.get('name', '(unnamed test case)')
        rprint(f'[case]{name}:[/case]')
        path = test_case["path"]
        path = str(Path(path).absolute())
        rprint(f'   [path]{path}[/path]')

        command_line_args = [
            './voyager',
            '--format=display',
            '--powerline',
        ]
        if extra_args:
            command_line_args += extra_args
        username = test_case.get('username')
        if username:
            command_line_args.append(f'--username={username}')
        command_line_args.append(path)

        environment = os.environ.copy()
        environment_vars = test_case.get('environment_vars')
        if environment_vars:
            for key, value in environment_vars.items():
                environment[key] = value

        # ------------------------
        # Render Powerline Prompt
        # ------------------------
        output = subprocess.check_output(command_line_args, env=environment,)
        print(f'   {output.decode("utf-8")}')

        # ------------------------
        # Render Text Prompt
        # ------------------------
        command_line_args.remove('--powerline')
        output = subprocess.check_output(command_line_args, env=environment,)
        print(f'   {output.decode("utf-8")}')
    print()


@cli.command()
@click.option('-c', '--colors', 'enable_color_modes', is_flag=True, help='Run test in all color modes')
def formats(enable_color_modes):
    for shell in ['zsh', 'bash']:
        rprint(f'[shell]{shell}[/shell]')
        for presentation in ("PowerLine", "Text"):
            rprint(f'   [presentation]{presentation}[/presentation]')
            show_formats(shell, presentation)

def show_formats(shell, presentation):
    TARGET_PATH = str(Path("./test1/test2").absolute())
    for _format in ('ics', 'prompt', 'display_debug', 'display'):
        rprint(f'      [format]{_format}[/format]')
        command_line_args = [
            './voyager',
            f'-format={_format}',
            f'-shell={shell}',
        ]
        if presentation == 'PowerLine':
            command_line_args.append('-powerline')
        command_line_args.append(TARGET_PATH)

        # ------------------------
        # Render Powerline Prompt
        # ------------------------
        output = subprocess.check_output(command_line_args)
        print(f'         {output.decode("utf-8")}')


if __name__ == "__main__":
    cli()
