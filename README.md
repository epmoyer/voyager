

## TODO
- Need to run `setopt prompt_subst` per https://stackoverflow.com/questions/15212152/zsh-prompt-function-not-running
    - Consider setting in the install script?
    - Nope, it belongs in the .zshrc (so add it to the voyager block)
- Figure out how to make it work in the standard apple "terminal" app, because presumably a bunch of people use that.
    - Seems that terminal just doesn't support RGB (24 bit colors):
        - https://github.com/mawww/kakoune/issues/1057
    - Probably have to add a separate palette just for xterm-256color support, and an option to enable that palette.