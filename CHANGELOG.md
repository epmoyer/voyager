# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## **(Unreleased)**

## 1.10.0 - 2024-04-02
### Changed
- If a virtual environment is running then show the basename (final path component) of $VIRTUAL_ENV instead of the default `venv` text.
    - This fix causes us to now display the actual name of a named `virtualenvwrapper` venv instead of the generic `venv`.

## 1.9.0 - 2023-06-30
### Changed
- Add support for Python virtual environments (venv).
    - Voyager prompt now displays "venv" as the current virtual environment if it detects that `$VIRTUAL_ENV` is set.
    - Voyager explicitly sets `VIRTUAL_ENV_DISABLE_PROMPT=1` during shell initialization, which tells Python venv to not manipulate the shell prompt (because Voyager is handling it). Without this, activating a venv would cause venv to inject "(venv) " before the prompt string, which is now redundant (because Voyager is now including "venv" inside it's own prompt string).

## 1.8.0 - 2023-03-15
### Fixed
- Properly render branch names that contain a slash.
    - Was previously rendering only the final "component" of the branch name.
### Added
- Test case for branch names containing a slash.

## 1.7.0 - 2023-01-30
### Changed
- Update revision for initial release.

## 1.6.1 - 2023-01-30
### Changed
- Use new `-ssh` flag instead of reading `SSH_CLIENT` inside voyager.
### Fixed
- Fix usage text for `-format`
- Virtual environment (now via `-virtualenv`)

## 1.6.0 - 2023-01-28
### Added
- Support color modes
### Changed
- Adopt ICS internally
- Snippet install location
- Improve demo coverage

## 1.5.2 - 2023-01-19
### Changed
- Adopt pink for `STYLE_ERROR`

## 1.5.2 - 2023-01-19
### Added
- Indicate errors with "X" prefix in dedicated error section. (via `VGER_RETVAL` environment variable).

## 1.5.1 - 2023-01-18
### Added
- Show error (X) if previous command returned nonzero value.
### Changed
- Refactor build system

## 1.5.0 - 2022-12-30
### Added
- Set path truncation depth with `VGER_TRUNCATION_START_DEPTH`
### Changed
- Rename project to 'voyager"

## 1.4.0 - 2022-12-29
### Changed
- Use `%` as final prompt symbol in text mode if shell is zsh.
- If user is root:
    - Text Mode:
        - Use red `#` as final prompt symbol
    - Powerline Mode:
        - Show context (root@host), and colorize red
### Fixed
- Robustly parse two letter status codes from `git status --porcelain`.
    - Any non-blank non-? status code in column 1 now considered as "Staged" (Index).
    - Any non-blank non-? status code in column 2 now considered as  "Modified" (Working Tree).
## 1.3.0 - 2022-12-28
### Fixed
- Robustly parse two letter status codes from `git status --porcelain`.
    - Any non-blank non-? status code in column 1 now considered as "Staged" (Index).
    - Any non-blank non-? status code in column 2 now considered as  "Modified" (Working Tree).

## 1.2.0 - 2022-12-28
### Changed
- Include git rename operations as `.IsStaged`.

## 1.1.0 - 2022-12-17
### Changed
- Remove bullnose from front of prompt (was rendering with a small gap at some scalings in iTerm2 Build 3.4.18, even when using their built-in Powerline symbols )

## 1.0.0 - 2022-12-16
Initial Release
