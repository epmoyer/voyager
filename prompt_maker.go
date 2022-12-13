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
const COLOR_BG_DEFAULT = "#000000"

const COLOR_TEXT_FG_PATH_CONTEXT = "#C040BE"
const COLOR_TEXT_FG_PATH_GITROOT = "#00d000"
const COLOR_TEXT_FG_PATH_GITSUB = "#6D8B8F"
const COLOR_TEXT_FG_SEPARATOR = "#d0d0d0"
const COLOR_TEXT_FG_GIT_INFO_CLEAN = "#A2C3C7"
const COLOR_TEXT_FG_GIT_INFO_DIRTY = "#E2D47D"

const COLOR_POWERLINE_FG_CONTEXT = "#000000"
const COLOR_POWERLINE_BG_CONTEXT = "#B294BF"
const COLOR_POWERLINE_FG_PATH_GITROOT_PRE = "#c0c0c0"
const COLOR_POWERLINE_FG_PATH_GITROOT = "#ffffff"
const COLOR_POWERLINE_BG_PATH_GITROOT = "#4F6D6F"
const COLOR_POWERLINE_FG_PATH_GITSUB = "#c0c0c0"
const COLOR_POWERLINE_BG_PATH_GITSUB = "#515151"
const COLOR_POWERLINE_FG_GIT_INFO = "#000000"
const COLOR_POWERLINE_BG_GIT_INFO_CLEAN = "#A2C3C7"
const COLOR_POWERLINE_BG_GIT_INFO_DIRTY = "#E2D47D"

type promptInfoT struct {
	Username    string
	UserHomeDir string
	Hostname    string
	PathGitRoot string
	PathGitSub  string
	GitBranch   string
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

		prompt := promptT{}
		prompt = prompt.addSegment(" test_context ", COLOR_POWERLINE_FG_CONTEXT, COLOR_POWERLINE_BG_CONTEXT, false)
		prompt = prompt.addSegment(" test_gitroot_dark/", COLOR_POWERLINE_FG_PATH_GITROOT_PRE, COLOR_POWERLINE_BG_PATH_GITROOT, true)
		prompt = prompt.addSegment("test_gitroot_light ", COLOR_POWERLINE_FG_PATH_GITROOT, COLOR_POWERLINE_BG_PATH_GITROOT, false)
		prompt = prompt.endSegments()
		fmt.Printf("PROMPT POWERLINE SEGMENT TEST:\n%s\n", prompt.Prompt)
		fmt.Println("-------------------------------------------------")
	}

	os.Exit(0)
}

func renderPrompt(usePowerline bool, promptInfo promptInfoT) string {
	contextColor := color.HEX(COLOR_TEXT_FG_PATH_CONTEXT)
	context := contextColor.Sprintf("%s", promptInfo.Username+"@"+promptInfo.Hostname)

	basePathColor := color.HEX(COLOR_TEXT_FG_PATH_GITROOT)
	basePath := basePathColor.Sprintf("%s", promptInfo.PathGitRoot)

	gitColor := color.HEX(COLOR_TEXT_FG_GIT_INFO_CLEAN)
	gitInfo := gitColor.Sprintf("%s", promptInfo.GitBranch)

	subPathColor := color.HEX(COLOR_TEXT_FG_PATH_GITSUB)
	subPath := subPathColor.Sprintf("%s", promptInfo.PathGitSub)

	separatorColor := color.HEX(COLOR_TEXT_FG_SEPARATOR)
	separator := separatorColor.Sprintf(" ⟫ ")

	prompt := context + separator + basePath + separator + gitInfo + separator + subPath + " $"
	return prompt
}

func renderPromptPowerline(promptInfo promptInfoT) string {
	// separator := "\ue0b0"

	contextColor := color.HEXStyle("#000000", COLOR_POWERLINE_BG_CONTEXT)
	context := contextColor.Sprintf(" %s ", promptInfo.Username+"@"+promptInfo.Hostname)

	basePathColor := color.HEXStyle(COLOR_POWERLINE_FG_PATH_GITROOT, COLOR_POWERLINE_BG_PATH_GITROOT)
	basePath := basePathColor.Sprintf(" %s ", promptInfo.PathGitRoot)

	subPathColor := color.HEXStyle(COLOR_POWERLINE_FG_PATH_GITSUB, COLOR_POWERLINE_BG_PATH_GITSUB)
	subPath := subPathColor.Sprintf(" %s ", promptInfo.PathGitSub)

	gitBackgroundColor := COLOR_POWERLINE_BG_GIT_INFO_CLEAN
	gitColor := color.HEXStyle(COLOR_POWERLINE_FG_GIT_INFO, gitBackgroundColor)
	gitInfo := gitColor.Sprintf(" %s ", promptInfo.GitBranch)

	prompt := (context +
		makeSeparator(COLOR_POWERLINE_BG_CONTEXT, COLOR_POWERLINE_BG_PATH_GITROOT) +
		basePath +
		makeSeparator(COLOR_POWERLINE_BG_PATH_GITROOT, gitBackgroundColor) +
		gitInfo +
		makeSeparator(gitBackgroundColor, COLOR_POWERLINE_BG_PATH_GITSUB) +
		subPath +
		makeSeparator(COLOR_POWERLINE_BG_PATH_GITSUB, COLOR_BG_DEFAULT) +
		" ")
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

	promptInfo.GitBranch = getGitBranch(path)

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

func getGitBranch(path string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	var e bytes.Buffer
	cmd.Stderr = &e
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		// This is not a git repo
		return ""
	}
	// TODO: If blank call "git rev-parse --short HEAD" for hash
	return strings.TrimSpace(string(out))
}
