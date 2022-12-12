package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/gookit/color"
)

const ENABLE_BOLD = false
const COLOR_FG_BOLD = "#ffffff"

type promptInfoT struct {
	Username    string
	UserHomeDir string
	Hostname    string
	PathGitRoot string
	PathGitSub  string
}

func main() {
	optDump := flag.Bool("dump", false,
		"Show all prompt components and all prompts in all formatting styles.")
	// requestGitSub := flag.Bool("gitsub", false,
	// 	"Get subdirectory relative to the git root.")
	// requestBoth := flag.Bool("both", false,
	// 	"Get current git root dir, and subdirectory relative to the git root, separated by `|`")
	flag.Parse()

	args := flag.Args()
	if len(args) > 1 {
		// Too many args.
		os.Exit(1)
	}
	path := strings.TrimSpace(args[0])
	// getPath(path)

	promptInfo, _ := buildPromptInfo(path)

	if *optDump {
		fmt.Println("Dump:")
		fmt.Println(path)
		fmt.Printf("%#v\n", promptInfo)
		promptText := renderPrompt(false, promptInfo)
		fmt.Println("-------------------------------------------------")
		fmt.Printf("PROMPT TEXT:\n%s\n", promptText)
		// promptPowerline = renderPrompt(true)
		promptPowerline := renderPromptPowerline(promptInfo)
		fmt.Printf("PROMPT POWERLINE:\n%s\n", promptPowerline)
		fmt.Println("-------------------------------------------------")
	}

	os.Exit(0)
}

func renderPrompt(usePowerline bool, promptInfo promptInfoT) string {
	contextColor := color.HEX("#C040BE")
	context := contextColor.Sprintf("%s", promptInfo.Username+"@"+promptInfo.Hostname)

	basePathColor := color.HEX("#808080")
	basePath := basePathColor.Sprintf("%s", promptInfo.PathGitRoot)

	subPathColor := color.HEX("#8080e0")
	subPath := subPathColor.Sprintf("%s", promptInfo.PathGitSub)

	separatorColor := color.HEX("#404040")
	separator := separatorColor.Sprintf(" ⟫ ")

	prompt := context + separator + basePath + separator + subPath + " $"
	return prompt
}

func renderPromptPowerline(promptInfo promptInfoT) string {
	separator := "\ue0b0"

	contextColor := color.HEXStyle("#000000", "#C040BE")
	context := contextColor.Sprintf(" %s ", promptInfo.Username+"@"+promptInfo.Hostname)

	// basePathColor := color.HEX("#808080")
	basePathColor := color.HEXStyle("#000000", "#A2C3C7")
	basePath := basePathColor.Sprintf(" %s", promptInfo.PathGitRoot)

	subPathColor := color.HEX("#8080e0")
	subPath := subPathColor.Sprintf("%s", promptInfo.PathGitSub)

	// separatorColor := color.HEX("#404040")
	// separator := separatorColor.Sprintf(" ⟫ ")

	prompt := context + makeSeparator("#C040BE", "#A2C3C7") + basePath + separator + subPath + " $"
	return prompt
}

func makeSeparator(oldColor string, newColor string) string {
	separatorStyle := color.HEXStyle(oldColor, newColor)
	return separatorStyle.Sprintf("\ue0b0")
}

func buildPromptInfo(path string) (promptInfoT, error) {

	promptInfo := promptInfoT{}

	pathGitRoot, pathGitSub := getPath(path)
	promptInfo.PathGitRoot = pathGitRoot
	promptInfo.PathGitSub = pathGitSub

	user, err := user.Current()
	if err != nil {
		return promptInfo, err
	}
	promptInfo.Username = user.Username
	promptInfo.UserHomeDir = user.HomeDir

	hostname, err := os.Hostname()
	if err != nil {
		return promptInfo, err
	}
	if strings.HasSuffix(hostname, ".local") {
		hostname = strings.Replace(hostname, ".local", "", 1)
	}
	promptInfo.Hostname = hostname

	return promptInfo, nil
}

func getPath(path string) (string, string) {

	usr, _ := user.Current()
	dir := usr.HomeDir
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", dir, 1)
	}

	pathGitRoot, pathGitSub := splitGitPath(path)

	if strings.HasPrefix(pathGitRoot, dir) {
		pathGitRoot = strings.Replace(pathGitRoot, dir, "~", 1)
	}
	pathGitRoot = shortenPath(pathGitRoot)

	return pathGitRoot, pathGitSub
}

func shortenPath(path string) string {
	pieces := strings.Split(path, "/")
	newPieces := []string{}
	var piece string
	for i := 0; i < len(pieces); i++ {
		piece = pieces[i]
		if i < (len(pieces) - 1) {
			newPieces = append(newPieces, shorten(piece))
			continue
		}
		if ENABLE_BOLD {
			piece = "%B%F{" + COLOR_FG_BOLD + "}" + piece + "%b%f"
		}
		newPieces = append(newPieces, piece)
	}
	return strings.Join(newPieces, "/")
}

func shorten(pathComponent string) string {
	if len(pathComponent) < 2 {
		return pathComponent
	}
	return pathComponent[0:1]
}

func splitGitPath(path string) (string, string) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var e bytes.Buffer
	cmd.Stderr = &e
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("🔴 Not a git repo")
		// This is not a git repo
		return path, ""
	}
	pathGitRoot := strings.TrimSpace(string(out))
	pathGitSub := strings.Replace(path, pathGitRoot, "", 1)
	if strings.HasPrefix(pathGitSub, "/") {
		pathGitSub = strings.Replace(pathGitSub, "/", "", 1)
	}

	return pathGitRoot, pathGitSub
}
