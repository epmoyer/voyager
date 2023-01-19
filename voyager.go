package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

const APP_NAME = "voyager"
const APP_VERSION = "1.5.2"

const ENABLE_DEBUG_INDICATOR = false
const ENABLE_BULLNOSE = false

const SYMBOL_PL_GIT_BRANCH = "\uf418"           // PowerLine: VCS Branch ()
const SYMBOL_PL_GIT_DETACHED = "\uf995"         // PowerLine: VCS Detached (秊)
const SYMBOL_PL_GIT_STAGED = "\uF0C7"           // PowerLine: Floppy Disk ()
const SYMBOL_PL_GIT_MODIFIED = "\uF448"         // PowerLine: Pencil ()
const SYMBOL_PL_GIT_BRANCH_AHEAD = "\uF0DE"     // PowerLine: Up-arrow
const SYMBOL_PL_GIT_BRANCH_BEHIND = "\uF0DD"    // PowerLine: Down-arrow
const SYMBOL_PL_GIT_BRANCH_UNTRACKED = "\uF128" // PowerLine: Question-mark ()
const SYMBOL_PL_SEPARATOR = "\ue0b0"            // PowerLine: Triangle-Right Separator
const SYMBOL_PL_BULLNOSE = "\ue0b6"             // PowerLine: Bullnose ()
const SYMBOL_PL_DOLLAR = "\uf155"               // Powerline: Dollar ()
const SYMBOL_PL_CHECK = "\uf00c"                // PowerLine: Check-mark ()
const SYMBOL_PL_X = "\uf00d"                    // PowerLine: X ()

const SYMBOL_TEXT_SEPARATOR = " ⟫ "

var SYMBOLS_POWERLINE = map[string]string{
	"branch":     SYMBOL_PL_GIT_BRANCH,
	"detached":   SYMBOL_PL_GIT_DETACHED,
	"staged":     SYMBOL_PL_GIT_STAGED + " ",
	"modified":   SYMBOL_PL_GIT_MODIFIED + " ",
	"untracked":  SYMBOL_PL_GIT_BRANCH_UNTRACKED + " ",
	"shell_bash": SYMBOL_PL_DOLLAR,
	"error":      SYMBOL_PL_X,
}
var SYMBOLS_TEXT = map[string]string{
	"branch":     "ʎ",
	"detached":   "➦",
	"staged":     "+",
	"modified":   "!",
	"untracked":  "?",
	"shell_bash": "$",
	"error":      "",
}

const COLOR_BG_DEFAULT = "#000000"
const COLOR_FG_DEFAULT = "#ffffff"
const COLOR_TEXT_FG_SEPARATOR = "#707070"

var STYLE_DEBUG = promptStyleT{
	ColorHexFGPowerline: "#000000",
	ColorHexBGPowerline: "#B7E2B7",
	ColorHexFGText:      "#B7E2B7",
}
var STYLE_ERROR = promptStyleT{
	ColorHexFGPowerline: "#000000",
	ColorHexBGPowerline: "#ff3030",
	ColorHexFGText:      "#ff3030",
}
var STYLE_SHELL = promptStyleT{
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
var STYLE_CONTEXT_ROOT = promptStyleT{
	ColorHexFGPowerline: "#000000",
	ColorHexBGPowerline: "#ff8080",
	ColorHexFGText:      "#ff3030",
}
var STYLE_GITROOT_PRE = promptStyleT{
	ColorHexFGPowerline: "#c0c0c0",
	ColorHexBGPowerline: "#4F6D6F",
	ColorHexFGText:      "#729E72",
}
var STYLE_GITROOT = promptStyleT{
	ColorHexFGPowerline: "#ffffff",
	ColorHexBGPowerline: "#4F6D6F",
	ColorHexFGText:      "#9EFF9E",
	Bold:                true,
}
var STYLE_GIT_INFO_CLEAN = promptStyleT{
	ColorHexFGPowerline: "#000000",
	ColorHexBGPowerline: "#A2C3C7",
	ColorHexFGText:      "#5EABF7",
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
	ColorHexFGText:      "#7A9CA1",
}

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		executable := filepath.Base(os.Args[0])
		fmt.Fprintf(w, "Usage of %s:\n", executable)
		fmt.Fprintf(w, "%s [options] <path>\n", executable)
		flag.PrintDefaults()
	}

	optVersion := flag.Bool("version", false,
		"Show version.")
	optDump := flag.Bool("dump", false,
		"Show all prompt components and all prompts in all formatting styles.")
	optPowerline := flag.Bool("powerline", false,
		"Render prompt using PowerLine font.")
	optShell := flag.String("shell", "zsh", "The shell to format the prompt for.")
	optUsername := flag.String("username", "", "Force the prompt username (for testing).")
	optPrintable := flag.Bool("printable", false,
		"Return a printable (ASCII Esc) string rather than a shell $PROMPT/$PS1 string.")
	flag.Parse()

	if *optVersion {
		showVersion()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) > 1 {
		// Too many args.
		os.Exit(1)
	}
	if len(args) == 0 {
		flag.Usage()
		os.Exit(0)
	}
	path := strings.TrimSpace(args[0])
	// fmt.Fprintf(os.Stderr, "args[0]:%#v\n", args[0])
	// fmt.Fprintf(os.Stderr, "path:%#v\n", path)

	promptInfo, _ := buildPromptInfo(path, *optUsername)

	prompt := promptT{}
	prompt.init(*optPowerline, *optShell)

	if *optDump {
		fmt.Printf("path:%#v\n", path)
		fmt.Printf("promptInfo:%#v\n", promptInfo)
		os.Exit(0)
	}

	prompt.renderPrompt(promptInfo)
	if *optPrintable {
		fmt.Print(prompt.TextPrintable)
	} else {
		fmt.Print(prompt.TextShell)
	}

	os.Exit(0)
}

func showVersion() {
	fmt.Printf("%s %s\n", APP_NAME, APP_VERSION)
}

func (prompt *promptT) renderPrompt(promptInfo promptInfoT) {
	var symbols map[string]string
	if prompt.isPowerline {
		symbols = SYMBOLS_POWERLINE
	} else {
		symbols = SYMBOLS_TEXT
	}

	// -----------------------
	// Debug
	// -----------------------
	if ENABLE_DEBUG_INDICATOR {
		prompt.addSegment(
			"",
			STYLE_DEBUG)
	}

	// -----------------------
	// Error
	// -----------------------
	if promptInfo.ReturnValue != 0 {
		prompt.addSegment(
			symbols["error"],
			STYLE_ERROR)
	}

	// -----------------------
	// Shell
	// -----------------------
	if prompt.colorizer.shell == "bash" {
		prompt.addSegment(
			symbols["shell_bash"],
			STYLE_SHELL)
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
	if promptInfo.ShowContext || promptInfo.IsRoot {
		styleContext := STYLE_CONTEXT
		if promptInfo.IsRoot {
			styleContext = STYLE_CONTEXT_ROOT
		}
		prompt.addSegment(
			fmt.Sprintf("%s@%s", promptInfo.Username, promptInfo.Hostname),
			styleContext)
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
	git := promptInfo.Git
	if git.IsRepo {
		style := STYLE_GIT_INFO_CLEAN
		if git.IsDirty {
			style = STYLE_GIT_INFO_DIRTY
		}
		if git.IsDetached {
			style = STYLE_GIT_INFO_DETACHED
		}
		text := git.render(prompt.isPowerline)
		prompt.addSegment(
			text,
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

	prompt.endSegments(promptInfo)
}

func buildPromptInfo(path string, optUsername string) (promptInfoT, error) {
	promptInfo := promptInfoT{}

	promptInfo.ShowContext = true

	pathGitRoot, pathGitSub := getPath(path)
	promptInfo.PathGitRootBeginning, promptInfo.PathGitRootFinal = chopPath(pathGitRoot)
	promptInfo.PathGitSub = pathGitSub

	// ---------------------
	// Previous command return value
	// ---------------------
	returnValue := os.Getenv("VGER_RETVAL")
	// fmt.Fprintf(os.Stderr, "VGER_RETVAL: %s\n", returnValue)
	if returnValue != "" {
		value, err := strconv.Atoi(returnValue)
		if err == nil {
			promptInfo.ReturnValue = value
			// fmt.Fprintf(os.Stderr, "value: %d\n", value)
		}
	}

	// ---------------------
	// User and Host
	// ---------------------
	user, err := user.Current()
	if err != nil {
		return promptInfo, err
	}
	if optUsername != "" {
		// Override username as specified on command line (for testing)
		promptInfo.Username = optUsername
	} else {
		promptInfo.Username = user.Username
	}
	if promptInfo.Username == "root" {
		promptInfo.IsRoot = true
	}
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
	if sshClient == "" && optUsername == "" {
		defaultUser := os.Getenv("DEFAULT_USER")
		// fmt.Printf("defaultUser:%#v", defaultUser)
		if defaultUser == promptInfo.Username {
			promptInfo.ShowContext = false
		}
	}

	// ---------------------
	// Git
	// ---------------------
	promptInfo.Git.update(path)

	return promptInfo, nil
}

func getPath(path string) (string, string) {

	usr, _ := user.Current()
	homeDir := usr.HomeDir
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", homeDir, 1)
	}

	pathGitRoot, pathGitSub := splitGitPath(path)

	if strings.HasPrefix(pathGitRoot, homeDir) {
		pathGitRoot = strings.Replace(pathGitRoot, homeDir, "~", 1)
	}
	pathGitRoot = shortenPath(pathGitRoot)

	return pathGitRoot, pathGitSub
}

func shortenPath(path string) string {
	truncationStartDepth := getPathTruncationStartDepth()
	pieces := strings.Split(path, "/")
	newPieces := []string{}
	var piece string
	for i := 0; i < len(pieces); i++ {
		piece = pieces[i]
		depth := len(pieces) - 1 - i
		if depth >= truncationStartDepth {
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

func finalComponent(path string) string {
	pieces := strings.Split(path, "/")
	return pieces[len(pieces)-1]
}

func getPathTruncationStartDepth() int {
	truncationStartDepthStr := os.Getenv("VGER_TRUNCATION_START_DEPTH")
	if truncationStartDepthStr == "" {
		return 1
	}
	truncationStartDepth, err := strconv.Atoi(truncationStartDepthStr)
	if err != nil {
		return 1
	}
	return truncationStartDepth
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
