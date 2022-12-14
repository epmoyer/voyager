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

const ENABLE_DEBUG_INDICATOR = true

// const SYMBOL_PL_GIT_BRANCH = "\ue0a0"           // PowerLine: VCS Branch
// const SYMBOL_PL_GIT_BRANCH = "\ue725"           // PowerLine: VCS Branch
const SYMBOL_PL_GIT_BRANCH = "\uf418"           // PowerLine: VCS Branch
const SYMBOL_PL_GIT_STAGED = "\uF0C7"           // PowerLine: Floppy Disk
const SYMBOL_PL_GIT_UNSTAGED = "\uF448"         // PowerLine: Pencil
const SYMBOL_PL_GIT_BRANCH_AHEAD = "\uF0DE"     // PowerLine: Up-arrow
const SYMBOL_PL_GIT_BRANCH_BEHIND = "\uF0DD"    // PowerLine: Down-arrow
const SYMBOL_PL_GIT_BRANCH_UNTRACKED = "\uF128" // PowerLine: Question-mark
const SYMBOL_PL_SEPARATOR = "\ue0b0"            // PowerLine: Triangle-Right Separator

const COLOR_BG_DEFAULT = "#000000"
const COLOR_TEXT_FG_SEPARATOR = "#707070"

var STYLE_DEBUG = promptStyleT{
	ColorHexFGPowerline: "#000000",
	ColorHexBGPowerline: "#FFA500",
	ColorHexFGText:      "#FFA500",
}
var STYLE_CONDA = promptStyleT{
	ColorHexFGPowerline: "#202020",
	ColorHexBGPowerline: "#5EABF7",
	ColorHexFGText:      "#4040ff",
}
var STYLE_CONTEXT = promptStyleT{
	ColorHexFGPowerline: "#000000",
	ColorHexBGPowerline: "#B294BF",
	ColorHexFGText:      "#C040BE",
}
var STYLE_GITROOT_PRE = promptStyleT{
	ColorHexFGPowerline: "#c0c0c0",
	ColorHexBGPowerline: "#4F6D6F",
	ColorHexFGText:      "#009000",
}
var STYLE_GITROOT = promptStyleT{
	ColorHexFGPowerline: "#ffffff",
	ColorHexBGPowerline: "#4F6D6F",
	ColorHexFGText:      "#30FF30",
	Bold:                true,
}
var STYLE_GIT_INFO_CLEAN = promptStyleT{
	ColorHexFGPowerline: "#000000",
	ColorHexBGPowerline: "#A2C3C7",
	ColorHexFGText:      "#A2C3C7",
}
var STYLE_GIT_INFO_DIRTY = promptStyleT{
	ColorHexFGPowerline: "#000000",
	ColorHexBGPowerline: "#E2D47D",
	ColorHexFGText:      "#E2D47D",
}
var STYLE_GITSUB = promptStyleT{
	ColorHexFGPowerline: "#c0c0c0",
	ColorHexBGPowerline: "#515151",
	ColorHexFGText:      "#6D8B8F",
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

	var isPowerline bool
	var prompt string

	if *optDump {
		fmt.Println("Dump:")
		fmt.Printf("ENV:GP_FORMAT=%#v\n", os.Getenv("GP_FORMAT"))
		fmt.Println(path)
		fmt.Printf("%#v\n", promptInfo)
		fmt.Println("-------------------------------------------------")

		isPowerline = false
		prompt = renderPrompt(promptInfo, isPowerline)
		fmt.Printf("TEXT PROMPT:\n%s\n", prompt)

		isPowerline = true
		prompt = renderPrompt(promptInfo, isPowerline)
		fmt.Printf("POWERLINE PROMPT:\n%s\n", prompt)

		// prompt := promptT{}
		// prompt = prompt.addSegment(" conda ", STYLE_POWERLINE_CONDA, false)
		// prompt = prompt.addSegment(" context ", STYLE_POWERLINE_CONTEXT, true)
		// prompt = prompt.addSegment(" gitroot_pre/", STYLE_GITROOT_PRE, true)
		// prompt = prompt.addSegment("final ", STYLE_GITROOT, false)
		// prompt = prompt.addSegment(" "+SYMBOL_PL_GIT_BRANCH+" git_info "+SYMBOL_PL_GIT_UNSTAGED+" ", STYLE_GIT_INFO_DIRTY, true)
		// prompt = prompt.addSegment(" gitsub ", STYLE_GITSUB, true)
		// prompt = prompt.endSegments()
		// fmt.Printf("PROMPT POWERLINE SEGMENT TEST:\n%s\n", prompt.Prompt)

		fmt.Println("-------------------------------------------------")
	}

	isPowerline = (os.Getenv("GP_FORMAT") == "POWERLINE")
	prompt = renderPrompt(promptInfo, isPowerline)
	fmt.Print(prompt)

	os.Exit(0)
}

func renderPrompt(promptInfo promptInfoT, isPowerline bool) string {
	prompt := promptT{
		isPowerline: isPowerline,
	}

	// -----------------------
	// Debug
	// -----------------------
	if ENABLE_DEBUG_INDICATOR {
		prompt.addSegment(
			"Debug",
			STYLE_DEBUG)
	}

	// -----------------------
	// Conda Environment
	// -----------------------
	if promptInfo.CondaEnvironment != "" {
		prompt.addSegment(
			fmt.Sprint(promptInfo.CondaEnvironment),
			STYLE_CONDA)
	}

	// -----------------------
	// Context
	// -----------------------
	if promptInfo.ShowContext {
		prompt.addSegment(
			fmt.Sprintf("%s@%s", promptInfo.Username, promptInfo.Hostname),
			STYLE_CONTEXT)
	}

	// -----------------------
	// Git root directory
	// -----------------------
	prompt.addSegment(
		fmt.Sprint(promptInfo.PathGitRootBeginning),
		STYLE_GITROOT_PRE)
	prompt.appendToSegment(
		fmt.Sprint(promptInfo.PathGitRootFinal),
		STYLE_GITROOT)

	// -----------------------
	// Git Status
	// -----------------------
	// TODO: Detect clean/dirty
	// TODO: Do nothing if not in a git dir
	if promptInfo.GitBranch != "" {
		var segmentText string
		if isPowerline {
			segmentText = fmt.Sprintf("%s %s", SYMBOL_PL_GIT_BRANCH, promptInfo.GitBranch)
		} else {
			segmentText = fmt.Sprint(promptInfo.GitBranch)
		}
		prompt.addSegment(
			segmentText,
			STYLE_GIT_INFO_CLEAN)
	}

	// -----------------------
	// Sub-directory within Git Repo
	// -----------------------
	if promptInfo.PathGitSub != "" {
		prompt.addSegment(
			fmt.Sprint(promptInfo.PathGitSub),
			STYLE_GITSUB)
	}

	prompt.endSegments()

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
			piece = shorten(piece)
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
