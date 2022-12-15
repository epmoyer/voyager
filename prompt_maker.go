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

const ENABLE_DEBUG_INDICATOR = false

const SYMBOL_PL_GIT_BRANCH = "\uf418"           // PowerLine: VCS Branch ()
const SYMBOL_PL_GIT_DETACHED = "\uf995"         // PowerLine: VCS Detached (秊)
const SYMBOL_PL_GIT_STAGED = "\uF0C7"           // PowerLine: Floppy Disk ()
const SYMBOL_PL_GIT_UNSTAGED = "\uF448"         // PowerLine: Pencil ()
const SYMBOL_PL_GIT_BRANCH_AHEAD = "\uF0DE"     // PowerLine: Up-arrow
const SYMBOL_PL_GIT_BRANCH_BEHIND = "\uF0DD"    // PowerLine: Down-arrow
const SYMBOL_PL_GIT_BRANCH_UNTRACKED = "\uF128" // PowerLine: Question-mark ()
const SYMBOL_PL_SEPARATOR = "\ue0b0"            // PowerLine: Triangle-Right Separator
const SYMBOL_PL_BULLNOSE = "\ue0b6"             // PowerLine: Bullnose ()
const SYMBOL_PL_CHECK = "\uf00c"                // PowerLine: Check-mark ()
const SYMBOL_PL_X = "\uf00d"                    // PowerLine: X ()

const COLOR_BG_DEFAULT = "#000000"
const COLOR_FG_DEFAULT = "#ffffff"
const COLOR_TEXT_FG_SEPARATOR = "#707070"

var STYLE_DEBUG = promptStyleT{
	ColorHexFGPowerline: "#000000",
	ColorHexBGPowerline: "#B8E3B8",
	ColorHexFGText:      "#B8E3B8",
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
var STYLE_GIT_INFO_DETACHED = promptStyleT{
	ColorHexFGPowerline: "#000000",
	ColorHexBGPowerline: "#FFAA55",
	ColorHexFGText:      "#FF8000",
}
var STYLE_GITSUB = promptStyleT{
	ColorHexFGPowerline: "#c0c0c0",
	ColorHexBGPowerline: "#515151",
	ColorHexFGText:      "#6D8B8F",
}

func main() {
	optDump := flag.Bool("dump", false,
		"Show all prompt components and all prompts in all formatting styles.")
	optPowerline := flag.Bool("powerline", false,
		"Render prompt using PowerLine font.")
	optShell := flag.String("shell", "zsh", "The shell to format the prompt for.")
	optPrintable := flag.Bool("printable", false,
		"Return a printable (ASCII Esc) string rather than a shell $PROMPT/$PS1 string.")
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

	prompt := promptT{}
	prompt.init(*optPowerline, *optShell)

	if *optDump {
		fmt.Printf("ENV:GP_FORMAT=%#v\n", os.Getenv("GP_FORMAT"))
		fmt.Printf("path:%#v\n", path)
		fmt.Printf("promptInfo:%#v\n", promptInfo)

		// isPowerline = false
		// prompt = renderPrompt(promptInfo, isPowerline)
		// fmt.Printf("TEXT PROMPT:\n%s\n%s\n", prompt.TextShell, prompt.TextPrintable)

		// isPowerline = true
		// prompt = renderPrompt(promptInfo, isPowerline)
		// fmt.Printf("POWERLINE PROMPT:\n%s\n%s\n", prompt.TextShell, prompt.TextPrintable)

		// prompt := promptT{}
		// prompt = prompt.addSegment(" conda ", STYLE_POWERLINE_CONDA, false)
		// prompt = prompt.addSegment(" context ", STYLE_POWERLINE_CONTEXT, true)
		// prompt = prompt.addSegment(" gitroot_pre/", STYLE_GITROOT_PRE, true)
		// prompt = prompt.addSegment("final ", STYLE_GITROOT, false)
		// prompt = prompt.addSegment(" "+SYMBOL_PL_GIT_BRANCH+" git_info "+SYMBOL_PL_GIT_UNSTAGED+" ", STYLE_GIT_INFO_DIRTY, true)
		// prompt = prompt.addSegment(" gitsub ", STYLE_GITSUB, true)
		// prompt = prompt.endSegments()
		// fmt.Printf("PROMPT POWERLINE SEGMENT TEST:\n%s\n", prompt.Prompt)

		// fmt.Println("-------------------------------------------------")
		os.Exit(0)
	}

	// isPowerline = (os.Getenv("GP_FORMAT") == "POWERLINE")
	prompt.renderPrompt(promptInfo)
	if *optPrintable {
		fmt.Print(prompt.TextPrintable)
	} else {
		fmt.Print(prompt.TextShell)
	}

	os.Exit(0)
}

func (prompt *promptT) renderPrompt(promptInfo promptInfoT) {
	// -----------------------
	// Debug
	// -----------------------
	if ENABLE_DEBUG_INDICATOR {
		prompt.addSegment(
			"",
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
	// Git Info
	// -----------------------
	// TODO: Detect clean/dirty
	// TODO: Do nothing if not in a git dir
	if promptInfo.GitBranch != "" {
		style := STYLE_GIT_INFO_CLEAN
		if promptInfo.GitStatus != "" {
			style = STYLE_GIT_INFO_DIRTY
		}
		if promptInfo.IsDetached {
			style = STYLE_GIT_INFO_DETACHED
		}
		var segmentText string
		if prompt.isPowerline {
			symbol := SYMBOL_PL_GIT_BRANCH
			if promptInfo.IsDetached {
				symbol = SYMBOL_PL_GIT_DETACHED
			}
			segmentText = fmt.Sprintf("%s %s", symbol, promptInfo.GitBranch)
			if promptInfo.GitStatus != "" {
				segmentText += " " + promptInfo.GitStatus
			}
		} else {
			segmentText = fmt.Sprint(promptInfo.GitBranch)
			// TODO: Probably don't use powerline fonts here. Find a way to do ASCII instead
			if promptInfo.GitStatus != "" {
				segmentText += " " + promptInfo.GitStatus
			}
		}
		prompt.addSegment(
			segmentText,
			style)
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
}

func buildPromptInfo(path string) (promptInfoT, error) {

	promptInfo := promptInfoT{}

	promptInfo.ShowContext = true

	pathGitRoot, pathGitSub := getPath(path)
	promptInfo.PathGitRootBeginning, promptInfo.PathGitRootFinal = chopPath(pathGitRoot)
	promptInfo.PathGitSub = pathGitSub

	// ---------------------
	// User and Host
	// ---------------------
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
	sshClient := os.Getenv("SSH_CLIENT")
	// fmt.Printf("sshClient:%#v", sshClient)
	if sshClient == "" {
		defaultUser := os.Getenv("DEFAULT_USER")
		// fmt.Printf("defaultUser:%#v", defaultUser)
		if defaultUser == promptInfo.Username {
			promptInfo.ShowContext = false
		}
	}

	// ---------------------
	// Git
	// ---------------------
	promptInfo.GitBranch, promptInfo.IsDetached = getGitBranch(path)
	if promptInfo.GitBranch != "" {
		promptInfo.GitStatus = getGitStatus(path)
	}

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

func getGitBranch(path string) (string, bool) {
	var cmd *exec.Cmd
	var e bytes.Buffer
	var out []byte
	var err error
	var reference string
	var isDetached bool

	cmd = exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Stderr = &e
	cmd.Dir = path
	out, err = cmd.Output()
	if err != nil {
		// This is not a git repo
		return "", isDetached
	}
	if strings.TrimSpace(string(out)) != "true" {
		// This is not a git repo
		return "", isDetached
	}

	cmd = exec.Command("git", "symbolic-ref", "HEAD")
	cmd.Stderr = &e
	cmd.Dir = path
	out, err = cmd.Output()
	if err == nil {
		reference = strings.TrimSpace(string(out))
	}
	if reference != "" {
		reference = finalComponent(reference)
	} else {
		// reference = "(other)"
		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
		var e bytes.Buffer
		cmd.Stderr = &e
		cmd.Dir = path
		out, err = cmd.Output()
		if err != nil {
			// This is not a git repo
			return "", isDetached
		}
		reference = "(" + strings.TrimSpace(string(out)) + ")"
		isDetached = true
	}
	return reference, isDetached
}

func getGitStatus(path string) string {
	var cmd *exec.Cmd
	var e bytes.Buffer
	var out []byte
	var err error
	var status string

	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Stderr = &e
	cmd.Dir = path
	out, err = cmd.Output()
	if err == nil {
		result := string(out)
		// TODO: These are sloppy checks.  Use proper regexes
		if strings.Contains(result, "??") {
			// UNTRACKED
			status += " "
		}
		if strings.Contains(result, "A ") {
			// STAGED
			status += " "
		}
		if strings.Contains(result, " M") {
			// MODIFIED
			status += " "
		}
	}
	return status
}

func finalComponent(path string) string {
	pieces := strings.Split(path, "/")
	return pieces[len(pieces)-1]
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
