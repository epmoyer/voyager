
# voyager:start ---------------------------------------------------------------
# voyager:snippet_version:1

# Allow substitutions and expansions in the prompt
setopt prompt_subst

export VGER_OPT_POWERLINE="--powerline"
vger_build_prompt() {
   export VGER_RETVAL=$?
   echo $(voyager $VGER_OPT_POWERLINE "$(pwd)")
}
export PROMPT='$(vger_build_prompt)'

# By default Voyager will truncate all but the last directory of your base (pre-git-repo) path.
# To truncate at a different start depth, uncomment the following.  For example, setting it to
# 3 will show the final 3 path components.  To show all path components, set it to a large
# number (e.g. 1000).
# export VGER_TRUNCATION_START_DEPTH=3

# --------------------
# Helper aliases
# --------------------
# These aliases provide commands you can run in the shell to 
# change the look of the voyager prompt on-the-fly

# Set prompt to "text" mode
alias vger_text="export VGER_OPT_POWERLINE=''"
# Set prompt to "powerline" mode
alias vger_pl="export VGER_OPT_POWERLINE=--powerline"
# Truncate all but the final path element.
alias vger_short="export VGER_TRUNCATION_START_DEPTH=1"
# Don't truncate any path elements.
alias vger_long="export VGER_TRUNCATION_START_DEPTH=1000"
# voyager:end -----------------------------------------------------------------
