# ------------------------------------------
# Voyager Prompt
# ------------------------------------------
export VGER_OPT_POWERLINE="--powerline"
vger_build_prompt(){
    PS1=`voyager $VGER_OPT_POWERLINE --shell=bash "$(pwd)"`
}
PROMPT_COMMAND=vger_build_prompt
vger_text() {
    export VGER_OPT_POWERLINE=""
}
vger_pl() {
    export VGER_OPT_POWERLINE="--powerline"
}
vger_ver() {
    echo $(voyager --version)
}