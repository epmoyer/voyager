
# voyager:start ---------------------------------------------------------------

# Allow substitutions and expansions in the prompt
setopt prompt_subst

export VGER_APP="/Users/eric/Dropbox (Personal)/cab_dbx/code/go/voyager/voyager.production"
export VGER_OPT_POWERLINE="--powerline"
vger_build_prompt() {
  echo $("$VGER_APP" $VGER_OPT_POWERLINE "$(pwd)")
}
vger_text() {
    export VGER_OPT_POWERLINE=""
}
vger_pl() {
    export VGER_OPT_POWERLINE="--powerline"
}
vger_ver() {
    echo $("$VGER_APP" --version)
}

export PROMPT='$(vger_build_prompt)'
# voyager:end -----------------------------------------------------------------
