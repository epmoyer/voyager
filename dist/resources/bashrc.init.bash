export VGER_OPT_POWERLINE="-powerline"
vger_build_prompt() {
    VGER_RETVAL=$?
 
    VGER_OPT_ERROR=""
    if [ $VGER_RETVAL -ne 0 ]; then
        VGER_OPT_ERROR="-showerror"
    fi

    VGER_OPT_DEFAULTUSER=""
    if [ ! -z "$VGER_DEFAULT_USER" ]; then
        VGER_OPT_DEFAULTUSER="-defaultuser=$VGER_DEFAULT_USER"
    fi

    VGER_OPT_TRUNCATION=""
     if [ ! -z "$VGER_TRUNCATION_START_DEPTH" ]; then
        VGER_OPT_TRUNCATION="-truncation=$VGER_TRUNCATION_START_DEPTH"
    fi

    VGER_OPT_COLOR=""
    if [ ! -z "$VGER_COLOR" ]; then
        VGER_OPT_COLOR="-color=$VGER_COLOR"
    fi

    VGER_OPT_VIRTUAL_ENVIRONMENT=""
    if [ ! -z "$CONDA_DEFAULT_ENV" ] && [ $CONDA_DEFAULT_ENV != "base" ]; then
        # Anaconda virtual environment is running
        VGER_OPT_VIRTUAL_ENVIRONMENT="-virtualenv=$CONDA_DEFAULT_ENV"
    fi
    if [ ! -z "$VIRTUAL_ENV" ]; then
        # Python virturl environment (i.e. venv) is running
        VGER_OPT_VIRTUAL_ENVIRONMENT="-virtualenv=venv"
    fi

    VGER_OPT_SSH=""
    if [ ! -z "$SSH_CLIENT" ]; then
        VGER_OPT_SSH="-ssh"
    fi

    PS1=`voyager $VGER_OPT_POWERLINE $VGER_OPT_ERROR $VGER_OPT_DEFAULTUSER $VGER_OPT_TRUNCATION $VGER_OPT_COLOR $VGER_OPT_VIRTUAL_ENVIRONMENT $VGER_OPT_SSH -shell=bash "$(pwd)"`
}
export PROMPT_COMMAND=vger_build_prompt

# Voyager will detect when a Python virtual environment (venv) is running and display it in
# the prompt, so we set the following to tell venv not to inject it's own "(venv) " string
# into the prompt environment variable
export VIRTUAL_ENV_DISABLE_PROMPT=1

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
alias vger_pl="export VGER_OPT_POWERLINE=-powerline"
# Truncate all but the final path element.
alias vger_short="export VGER_TRUNCATION_START_DEPTH=1"
# Don't truncate any path elements.
alias vger_long="export VGER_TRUNCATION_START_DEPTH=1000"
# Set color modes
alias vger_16m="export VGER_COLOR=16m"
alias vger_256="export VGER_COLOR=256"
alias vger_16="export VGER_COLOR=16"
alias vger_none="export VGER_COLOR=none"
