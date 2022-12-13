package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// const SYMBOL_GIT_BRANCH = "\ue0a0"           // PowerLine: VCS Branch
// const SYMBOL_GIT_BRANCH = "\ue725"           // PowerLine: VCS Branch
const SYMBOL_GIT_BRANCH = "\uf418"           // PowerLine: VCS Branch
const SYMBOL_GIT_STAGED = "\uF0C7"           // PowerLine: Floppy Disk
const SYMBOL_GIT_UNSTAGED = "\uF448"         // PowerLine: Pencil
const SYMBOL_GIT_BRANCH_AHEAD = "\uF0DE"     // PowerLine: Up-arrow
const SYMBOL_GIT_BRANCH_BEHIND = "\uF0DD"    // PowerLine: Down-arrow
const SYMBOL_GIT_BRANCH_UNTRACKED = "\uF128" // PowerLine: Question-mark
const SYMBOL_SEPARATOR = "\ue0b0"            // PowerLine: Triangle-Right Separator

const ENABLE_BOLD = false
const COLOR_FG_BOLD = "#ffffff"
const COLOR_BG_DEFAULT = "#000000"

const COLOR_TEXT_FG_PATH_CONDA = "#4040ff"
const COLOR_TEXT_FG_PATH_CONTEXT = "#C040BE"
const COLOR_TEXT_FG_PATH_GITROOT = "#009000"
const COLOR_TEXT_FG_PATH_GITROOT_FINAL = "#30FF30"
const COLOR_TEXT_FG_PATH_GITSUB = "#6D8B8F"

const COLOR_TEXT_FG_SEPARATOR = "#707070"
const COLOR_TEXT_FG_GIT_INFO_CLEAN = "#A2C3C7"
const COLOR_TEXT_FG_GIT_INFO_DIRTY = "#E2D47D"

var STYLE_POWERLINE_CONDA = promptStyleT{ColorHexFG: "#202020", ColorHexBG: "#5EABF7"}
var STYLE_POWERLINE_CONTEXT = promptStyleT{ColorHexFG: "#000000", ColorHexBG: "#B294BF"}
var STYLE_GITROOT_PRE = promptStyleT{ColorHexFG: "#c0c0c0", ColorHexBG: "#4F6D6F"}
var STYLE_GITROOT = promptStyleT{ColorHexFG: "#ffffff", ColorHexBG: "#4F6D6F", Bold: true}
var STYLE_GIT_INFO_CLEAN = promptStyleT{ColorHexFG: "#000000", ColorHexBG: "#A2C3C7"}
var STYLE_GIT_INFO_DIRTY = promptStyleT{ColorHexFG: "#000000", ColorHexBG: "#E2D47D"}
var STYLE_GITSUB = promptStyleT{ColorHexFG: "#c0c0c0", ColorHexBG: "#515151"}

type promptInfoT struct {
	CondaEnvironment     string
	Username             string
	UserHomeDir          string
	ShowContext          bool
	Hostname             string
	PathGitRootBeginning string
	PathGitRootFinal     string
	PathGitSub           string
	GitBranch            string
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

		// prompt := promptT{}

		// prompt = prompt.addSegment(" conda ", STYLE_POWERLINE_CONDA, false)
		// prompt = prompt.addSegment(" context ", STYLE_POWERLINE_CONTEXT, true)
		// prompt = prompt.addSegment(" gitroot_pre/", STYLE_GITROOT_PRE, true)
		// prompt = prompt.addSegment("final ", STYLE_GITROOT, false)
		// prompt = prompt.addSegment(" "+SYMBOL_GIT_BRANCH+" git_info "+SYMBOL_GIT_UNSTAGED+" ", STYLE_GIT_INFO_DIRTY, true)
		// prompt = prompt.addSegment(" gitsub ", STYLE_GITSUB, true)

		// prompt = prompt.endSegments()
		// fmt.Printf("PROMPT POWERLINE SEGMENT TEST:\n%s\n", prompt.Prompt)
		fmt.Println("-------------------------------------------------")
	}

	os.Exit(0)
}

func renderPrompt(usePowerline bool, promptInfo promptInfoT) string {
	prompt := textPromptT{}
	// -----------------------
	// Conda Environment
	// -----------------------
	if promptInfo.CondaEnvironment != "" {
		prompt = prompt.addSegment(
			fmt.Sprintf(" %s ", promptInfo.CondaEnvironment),
			COLOR_TEXT_FG_PATH_CONDA,
			true)
	}

	// -----------------------
	// Context
	// -----------------------
	if promptInfo.ShowContext {
		prompt = prompt.addSegment(
			fmt.Sprintf(" %s@%s ", promptInfo.Username, promptInfo.Hostname),
			COLOR_TEXT_FG_PATH_CONTEXT, true)
	}

	// -----------------------
	// Git root directory
	// -----------------------
	prompt = prompt.addSegment(
		fmt.Sprintf(" %s", promptInfo.PathGitRootBeginning),
		COLOR_TEXT_FG_PATH_GITROOT, true)
	prompt = prompt.addSegment(
		fmt.Sprintf("%s ", promptInfo.PathGitRootFinal),
		COLOR_TEXT_FG_PATH_GITROOT_FINAL, false)

	// -----------------------
	// Git Status
	// -----------------------
	// TODO: Detect clean/dirty
	// TODO: Do nothing if not in a git dir
	if promptInfo.GitBranch != "" {
		prompt = prompt.addSegment(
			fmt.Sprintf(" %s ", promptInfo.GitBranch),
			COLOR_TEXT_FG_GIT_INFO_CLEAN,
			true)
	}

	// -----------------------
	// Sub-directory within Git Repo
	// -----------------------
	if promptInfo.PathGitSub != "" {
		prompt = prompt.addSegment(
			fmt.Sprintf(" %s ", promptInfo.PathGitSub),
			COLOR_TEXT_FG_PATH_GITSUB,
			true)
	}

	// contextColor := color.HEX(COLOR_TEXT_FG_PATH_CONTEXT)
	// context := contextColor.Sprintf("%s", promptInfo.Username+"@"+promptInfo.Hostname)

	// basePathColor := color.HEX(COLOR_TEXT_FG_PATH_GITROOT)
	// basePathBeginning := basePathColor.Sprintf("%s", promptInfo.PathGitRootBeginning)

	// basePathFinalColor := color.HEX(COLOR_TEXT_FG_PATH_GITROOT_FINAL)
	// basePathFinal := basePathFinalColor.Sprintf("%s", promptInfo.PathGitRootFinal)

	// gitColor := color.HEX(COLOR_TEXT_FG_GIT_INFO_CLEAN)
	// gitInfo := gitColor.Sprintf("%s", promptInfo.GitBranch)

	// subPathColor := color.HEX(COLOR_TEXT_FG_PATH_GITSUB)
	// subPath := subPathColor.Sprintf("%s", promptInfo.PathGitSub)

	// separatorColor := color.HEX(COLOR_TEXT_FG_SEPARATOR)
	// separator := separatorColor.Sprintf(" ⟫ ")

	// separator := separatorColor.Sprintf(" > ")

	// prompt := context + separator + basePathBeginning + basePathFinal + separator + gitInfo + separator + subPath + " $"

	prompt = prompt.endSegments()

	return prompt.Prompt
}

func renderPromptPowerline(promptInfo promptInfoT) string {
	// -----------------------
	// Conda Environment
	// -----------------------
	prompt := powerlinePromptT{}
	if promptInfo.CondaEnvironment != "" {
		prompt = prompt.addSegment(
			fmt.Sprintf(" %s ", promptInfo.CondaEnvironment),
			STYLE_POWERLINE_CONDA,
			true)
	}

	// -----------------------
	// Context
	// -----------------------
	if promptInfo.ShowContext {
		prompt = prompt.addSegment(
			fmt.Sprintf(" %s@%s ", promptInfo.Username, promptInfo.Hostname),
			STYLE_POWERLINE_CONTEXT, true)
	}

	// -----------------------
	// Git root directory
	// -----------------------
	prompt = prompt.addSegment(
		fmt.Sprintf(" %s", promptInfo.PathGitRootBeginning),
		STYLE_GITROOT_PRE, true)
	prompt = prompt.addSegment(
		fmt.Sprintf("%s ", promptInfo.PathGitRootFinal),
		STYLE_GITROOT, false)

	// -----------------------
	// Git Status
	// -----------------------
	// TODO: Detect clean/dirty
	// TODO: Do nothing if not in a git dir
	if promptInfo.GitBranch != "" {
		prompt = prompt.addSegment(
			fmt.Sprintf(" %s %s ", SYMBOL_GIT_BRANCH, promptInfo.GitBranch),
			STYLE_GIT_INFO_CLEAN,
			true)
	}

	// -----------------------
	// Sub-directory within Git Repo
	// -----------------------
	if promptInfo.PathGitSub != "" {
		prompt = prompt.addSegment(
			fmt.Sprintf(" %s ", promptInfo.PathGitSub),
			STYLE_GITSUB,
			true)
	}

	prompt = prompt.endSegments()

	return prompt.Prompt
}

func buildPromptInfo(path string) (promptInfoT, error) {

	promptInfo := promptInfoT{}

	// TODO: Show conditionally
	promptInfo.ShowContext = true

	pathGitRoot, pathGitSub := getPath(path)
	promptInfo.PathGitRootBeginning, promptInfo.PathGitRootFinal = chopPath(pathGitRoot)
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

func chopPath(path string) (string, string) {
	pieces := strings.Split(path, "/")
	newPieces := []string{}
	var piece string
	var finalComponent string
	for i := 0; i < len(pieces); i++ {
		piece = pieces[i]
		if i < (len(pieces) - 1) {
			newPieces = append(newPieces, piece)
		} else {
			finalComponent = piece
		}
	}
	if len(newPieces) > 0 {
		// This will cause a trailing slash in the base path if a base path exists
		newPieces = append(newPieces, "")
	}
	return strings.Join(newPieces, "/"), finalComponent
}
