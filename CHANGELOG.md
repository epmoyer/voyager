# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.5.0 - 2022-12-30
(WIP)
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
