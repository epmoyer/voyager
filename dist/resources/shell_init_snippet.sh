
# voyager:start ---------------------------------------------------------------

# Allow substitutions and expansions in the prompt
setopt prompt_subst

export VGER_OPT_POWERLINE="--powerline"
vger_build_prompt() {
  echo $(voyager $VGER_OPT_POWERLINE "$(pwd)")
}
export PROMPT='$(vger_build_prompt)'

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
