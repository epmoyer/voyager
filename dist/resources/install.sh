#/bin/bash

RED=$'\033[31m'
YELLOW=$'\033[33m'
GREEN=$'\033[32m'
ENDCOLOR=$'\033[0m'

APP_NAME="voyager"

echo "Installing to /usr/local/bin.  You may be prompted for sudo permissions..."
sudo cp $APP_NAME /usr/local/bin
sudo chmod 755 /usr/local/bin/$APP_NAME
echo "${GREEN}   Copied.${ENDCOLOR}"

install_shell_snippet() {
    echo "         Adding to() function to $SHELL_INIT_SCRIPT..."
    cat $SNIPPET_SCRIPT >> $SHELL_INIT_SCRIPT
    echo "         ${GREEN}Added.${ENDCOLOR}"
}

query_install_shell_snippet() {
    read -p "      Add shell snippet to $SHELL_INIT_SCRIPT ? " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]
    then
        install_shell_snippet $SHELL_INIT_SCRIPT
    else
        echo "         ${YELLOW}(Skipped)${ENDCOLOR}"
    fi
}

check_shell_init_script () {
    echo "      looking for existing shell init snippet..."
    if grep -Eq "^# $APP_NAME:start" $SHELL_INIT_SCRIPT
    then
        echo "         ${GREEN}Found.${ENDCOLOR}"
    else
        echo "         Not found."
        query_install_shell_snippet $SHELL_INIT_SCRIPT $SNIPPET_SCRIPT
    fi
}

process_shell_script () {
    SHELL_INIT_SCRIPT=$1
    SNIPPET_SCRIPT=$2
    echo "   Looking for $SHELL_INIT_SCRIPT..."
    if test -f "$SHELL_INIT_SCRIPT"; then
        echo "      Found."
        check_shell_init_script $SHELL_INIT_SCRIPT $SNIPPET_SCRIPT
    else
        echo "      (Does not exist)"
    fi
}

echo "Adding to() function to shell script..."
process_shell_script ~/.bashrc shell_init_snippet_bash.sh
process_shell_script ~/.zshrc shell_init_snippet_zsh.sh

echo "${GREEN}Done.${ENDCOLOR}"