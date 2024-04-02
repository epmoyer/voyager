package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

const APP_NAME = "voyager"
const APP_VERSION = "1.10.0"

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

const ICS_COLOR_TEXT_FG_SEPARATOR = "white:241:#707070"
const ICS_RESET_ALL = "%{%f%k%b%}"

const DEBUG_ENABLE = false

var colorMode int = ColorMode16

var STYLE_DEBUG = promptStyleT{
	ICSColorBGPowerline: "brightgreen:151:#B7E2B7",
	ICSColorFGPowerline: "black:0:#000000",
	ICSColorFGText:      "brightgreen:151:#B7E2B7",
}
var STYLE_ERROR = promptStyleT{
	ICSColorBGPowerline: "brightred:212:#ff92c5",
	ICSColorFGPowerline: "black:16:#000000",
	ICSColorFGText:      "red:212:#ff92c5",
}
var STYLE_SHELL = promptStyleT{
	ICSColorBGPowerline: "yellow:151:#B8E3B8",
	ICSColorFGPowerline: "black:16:#000000",
	ICSColorFGText:      "yellow:151:#B8E3B8",
}
var STYLE_VIRTUAL_ENVIRONMENT = promptStyleT{
	ICSColorBGPowerline: "blue:75:#5EABF7",
	ICSColorFGPowerline: "black:16:#202020",
	ICSColorFGText:      "blue:63:#4040ff",
}
var STYLE_CONTEXT = promptStyleT{
	ICSColorBGPowerline: "brightmagenta:139:#B294BF",
	ICSColorFGPowerline: "black:16:#000000",
	ICSColorFGText:      "brightmagenta:133:#C040BE",
}
var STYLE_CONTEXT_ROOT = promptStyleT{
	ICSColorBGPowerline: "brightred:210:#ff8080",
	ICSColorFGPowerline: "black:16:#000000",
	ICSColorFGText:      "red:197:#ff3030",
}
var STYLE_GITROOT_PRE = promptStyleT{
	ICSColorBGPowerline: "green:66:#4F6D6F",
	ICSColorFGPowerline: "black:251:#c0c0c0",
	ICSColorFGText:      "green:70:#729E72",
}
var STYLE_GITROOT = promptStyleT{
	ICSColorBGPowerline: "green:66:#4F6D6F",
	ICSColorFGPowerline: "brightwhite:231:#ffffff",
	ICSColorFGText:      "white:157:#9EFF9E",
	Bold:                true,
}
var STYLE_GIT_INFO_CLEAN = promptStyleT{
	ICSColorBGPowerline: "cyan:152:#A2C3C7",
	ICSColorFGPowerline: "black:16:#000000",
	ICSColorFGText:      "cyan:75:#5EABF7",
}
var STYLE_GIT_INFO_DIRTY = promptStyleT{
	ICSColorBGPowerline: "brightyellow:186:#E2D47D",
	ICSColorFGPowerline: "black:16:#000000",
	ICSColorFGText:      "brightyellow:186:#E2D47D",
}
var STYLE_GIT_INFO_DETACHED = promptStyleT{
	ICSColorBGPowerline: "magenta:215:#FFAA55",
	ICSColorFGPowerline: "black:16:#000000",
	ICSColorFGText:      "magenta:208:#FF8000",
}
var STYLE_GITSUB = promptStyleT{
	ICSColorBGPowerline: "brightblack:59:#515151",
	ICSColorFGPowerline: "white:145:#c0c0c0",
	ICSColorFGText:      "brightblack:109:#7A9CA1",
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
		"DEBUG: Show all internal prompt data components.")
	optPowerline := flag.Bool("powerline", false,
		"Render prompt using PowerLine font.")
	optShell := flag.String("shell", "zsh", "The shell to format the prompt for.")
	optSSH := flag.Bool("ssh", false, "User is currently SSH'd into this machine.")
	optUsername := flag.String("username", "", "DEBUG: Force the prompt username.")
	optDefaultUser := flag.String("defaultuser", "", "The default username (don't show user/host for this user).")
	optError := flag.Bool("showerror", false, "Show the error indicator.")
	optColor := flag.String("color", "16m",
		"Set color mode. Can be set to any of: [\"16\", \"256\", \"16m\", \"none\"].")
	optFormat := flag.String("format", "prompt",
		"Output format. Can be any of: [\"prompt\", \"prompt_debug\", \"display\", \"display_debug\", \"ics\"].")
	optTruncationStartDepth := flag.Int("truncation", 1,
		"How many path components (right to left) to show in full. The rest will be truncated to a single character.")
	optVirtualEnv := flag.String("virtualenv", "",
		"Virtual environment name to display.")
	flag.Parse()

	setColorMode(*optColor)

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

	promptInfo, _ := buildPromptInfo(
		path,
		*optUsername,
		*optError,
		*optDefaultUser,
		*optTruncationStartDepth,
		*optVirtualEnv,
		*optSSH)

	prompt := promptT{}
	prompt.init(*optPowerline, *optShell)

	if *optDump {
		fmt.Printf("path:%#v\n", path)
		fmt.Printf("promptInfo:%#v\n", promptInfo)
		os.Exit(0)
	}

	prompt.build(promptInfo)
	fmt.Print(prompt.render(*optFormat))

	os.Exit(0)
}

func showVersion() {
	fmt.Printf("%s %s\n", APP_NAME, APP_VERSION)
}

func (prompt *promptT) build(promptInfo promptInfoT) {
	var symbols map[string]string
	if prompt.IsPowerLine {
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
	if promptInfo.ShowErrorIndicator {
		prompt.addSegment(
			symbols["error"],
			STYLE_ERROR)
	}

	// -----------------------
	// Shell
	// -----------------------
	if prompt.Shell == "bash" {
		prompt.addSegment(
			symbols["shell_bash"],
			STYLE_SHELL)
	}

	// -----------------------
	// Virtual Environment
	// -----------------------
	if promptInfo.VirtualEnvironment != "" {
		prompt.addSegment(
			fmt.Sprint(promptInfo.VirtualEnvironment),
			STYLE_VIRTUAL_ENVIRONMENT)
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
		text := git.render(prompt.IsPowerLine)
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

func buildPromptInfo(
	path string,
	optUsername string,
	optError bool,
	optDefaultUser string,
	optTruncationStartDepth int,
	optVirtualEnv string,
	optSSH bool,
) (promptInfoT, error) {

	promptInfo := promptInfoT{}

	promptInfo.ShowContext = true
	promptInfo.VirtualEnvironment = optVirtualEnv

	pathGitRoot, pathGitSub := getPath(path, optTruncationStartDepth)
	promptInfo.PathGitRootBeginning, promptInfo.PathGitRootFinal = chopPath(pathGitRoot)
	promptInfo.PathGitSub = pathGitSub

	// ---------------------
	// Show error indicator
	// ---------------------
	if optError {
		promptInfo.ShowErrorIndicator = true
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
	username := promptInfo.Username
	if optUsername != "" {
		// A username was injected on the command line (for testing)
		username = optUsername
	}

	// Show/Hide context (i.e. user & hostname)
	if optSSH {
		// Always show context when connected over SSH.
		promptInfo.ShowContext = true
	} else if username != optDefaultUser {
		// Show context if the current user is NOT the default user
		promptInfo.ShowContext = true
	} else {
		promptInfo.ShowContext = false
	}

	// ---------------------
	// Git
	// ---------------------
	promptInfo.Git.update(path)

	return promptInfo, nil
}

func getPath(path string, optTruncationStartDepth int) (string, string) {

	usr, _ := user.Current()
	homeDir := usr.HomeDir
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", homeDir, 1)
	}

	pathGitRoot, pathGitSub := splitGitPath(path)

	if strings.HasPrefix(pathGitRoot, homeDir) {
		pathGitRoot = strings.Replace(pathGitRoot, homeDir, "~", 1)
	}
	pathGitRoot = shortenPath(pathGitRoot, optTruncationStartDepth)

	return pathGitRoot, pathGitSub
}

func shortenPath(path string, optTruncationStartDepth int) string {
	pieces := strings.Split(path, "/")
	newPieces := []string{}
	var piece string
	for i := 0; i < len(pieces); i++ {
		piece = pieces[i]
		depth := len(pieces) - 1 - i
		if depth >= optTruncationStartDepth {
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

func setColorMode(optColor string) {
	switch optColor {
	case "none":
		colorMode = ColorModeNone
	case "16":
		colorMode = ColorMode16
	case "256":
		colorMode = ColorMode256
	// "16m"
	default:
		colorMode = ColorMode16m
	}
}
