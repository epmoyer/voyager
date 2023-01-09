
# voyager:start ---------------------------------------------------------------

# Allow substitutions and expansions in the prompt
setopt prompt_subst

# By default Voyager will truncate all but the last directory of your base (pre-git-repo) path.
# To truncate at a different start depth, uncomment the following.  For example, setting it to
# 3 will show the final 3 path components.  To show all path components, set it to a large 
# number (e.g. 1000).
# export VGER_TRUNCATION_START_DEPTH=3

export VGER_OPT_POWERLINE="--powerline"
vger_build_prompt() {
  echo $(voyager $VGER_OPT_POWERLINE "$(pwd)")
}
vger_text() {
    export VGER_OPT_POWERLINE=""
}
vger_pl() {
    export VGER_OPT_POWERLINE="--powerline"
}
vger_ver() {
    echo $(voyager --version)
}

export PROMPT='$(vger_build_prompt)'
# voyager:end -----------------------------------------------------------------
