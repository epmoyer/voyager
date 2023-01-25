#!/usr/bin/env python3
import os
import subprocess
from pathlib import Path
import textwrap

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
    "renderer": "#00ffff",
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
        rprint(f'{indent(1)}[path]{path}[/path]')

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
        print(f'{indent(1)}{output.decode("utf-8")}')

        # ------------------------
        # Render Text Prompt
        # ------------------------
        command_line_args.remove('--powerline')
        output = subprocess.check_output(command_line_args, env=environment,)
        print(f'{indent(1)}{output.decode("utf-8")}')
    print()


@cli.command()
@click.option('-c', '--colors', 'enable_color_modes', is_flag=True, help='Run test in all color modes')
@click.option('-n', '--nowrap', 'disable_text_wrap', is_flag=True, help='Do not text wrap.')
def formats(enable_color_modes, disable_text_wrap):
    for shell in ['zsh', 'bash']:
        rprint(f'[shell]{shell}[/shell]')
        for presentation in ("PowerLine", "Text"):
            rprint(f'{indent(1)}[presentation]{presentation}[/presentation]')
            show_formats(shell, presentation, disable_text_wrap)

def show_formats(shell, presentation, disable_text_wrap):
    TARGET_PATH = str(Path("./test1/test2").absolute())
    for _format in ('ics', 'prompt_debug', 'prompt', 'display_debug', 'display'):
        rprint(f'{indent(2)}[format]{_format}[/format]')
        command_line_args = [
            './voyager',
            f'-format={_format}',
            f'-shell={shell}',
        ]
        if presentation == 'PowerLine':
            command_line_args.append('-powerline')
        command_line_args.append(TARGET_PATH)

        # Execute
        output = subprocess.check_output(command_line_args)
        output = output.decode("utf-8")

        # Render prompt
        if _format == 'prompt':
            output = render_prompt(output, shell)
        
        if _format in ('ics', 'prompt_debug', 'display_debug') and not disable_text_wrap:
            wrapped = textwrap.wrap(output, width=80)
            for line in wrapped:
                print(f'{indent(3)}{line}')
                # print(wrapped)
        else:
            print(f'{indent(3)}{output}')


def render_prompt(prompt_text, shell):
    renderer = f'(no prompt renderer implemented for shell: "{shell}")'
    rendered_output = "(no renderer)"

    if shell == 'bash':
        renderer = 'Simulated bash renderer'
        rendered_output = prompt_text.replace(r'\[', '').replace(r'\]', '')
    elif shell == 'zsh':
        renderer = 'zsh (using "print -P")'
        command_line_args = [
            'zsh',
            '-c',
            # 'print',
            # '-P',
            f'print -P "{prompt_text}"'
        ]
        output = subprocess.check_output(command_line_args)
        rendered_output = output.decode("utf-8")

    rprint(f'{indent(3)}[renderer]Renderer: {renderer}[/renderer]')
    return rendered_output

def indent(level):
    return " " * 4 * level
if __name__ == "__main__":
    cli()
