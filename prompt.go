package main

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/gookit/color"
)

type gitInfoT struct {
	IsRepo      bool
	Branch      string
	IsDetached  bool
	IsDirty     bool
	IsStaged    bool
	IsModified  bool
	IsUntracked bool
}

type promptInfoT struct {
	CondaEnvironment     string
	Username             string
	UserHomeDir          string
	ShowContext          bool
	Hostname             string
	PathGitRootBeginning string
	PathGitRootFinal     string
	PathGitSub           string
	Git                  gitInfoT
}

type promptT struct {
	TextPrintable     string
	TextShell         string
	CurrentBGColorHex string
	isPowerline       bool
	colorizer         colorizerT
}

type promptStyleT struct {
	ColorHexFGPowerline string
	ColorHexBGPowerline string
	ColorHexFGText      string
	Bold                bool
}

func (prompt *promptT) init(isPowerline bool, shell string) {
	prompt.colorizer = colorizerT{}
	prompt.colorizer.shell = shell
	prompt.isPowerline = isPowerline
}

func (prompt *promptT) addSegment(text string, style promptStyleT) {
	if prompt.isPowerline && !(prompt.TextPrintable == "" && ENABLE_BULLNOSE) {
		// Powerline prompt gets a leading space
		text = " " + text
	}
	if prompt.TextPrintable == "" {
		// -------------------
		//  First segment: Start with bull-nose
		// -------------------
		if prompt.isPowerline && ENABLE_BULLNOSE {
			bullnoseStyle := color.HEXStyle(style.ColorHexBGPowerline)
			prompt.TextPrintable += bullnoseStyle.Sprint(SYMBOL_PL_BULLNOSE)

			// SHELL
			// prompt.TextShell += " "
			prompt.TextShell += prompt.colorizer.colorize(SYMBOL_PL_BULLNOSE, style.ColorHexBGPowerline, "", style.Bold)
		}
	} else {
		// -------------------
		//  Add Separator
		// -------------------
		if prompt.isPowerline {
			separatorStyle := color.HEXStyle(COLOR_FG_DEFAULT, prompt.CurrentBGColorHex)
			prompt.TextPrintable += separatorStyle.Sprint(" ")
			separatorStyle = color.HEXStyle(prompt.CurrentBGColorHex, style.ColorHexBGPowerline)
			prompt.TextPrintable += separatorStyle.Sprintf("%s", SYMBOL_PL_SEPARATOR)

			// SHELL
			prompt.TextShell += " "
			prompt.TextShell += prompt.colorizer.colorize(SYMBOL_PL_SEPARATOR, prompt.CurrentBGColorHex, style.ColorHexBGPowerline, style.Bold)
		} else {
			separatorColor := color.HEX(COLOR_TEXT_FG_SEPARATOR)
			// prompt.Prompt += separatorColor.Sprintf(" ⟫ ")
			prompt.TextPrintable += separatorColor.Sprintf(SYMBOL_TEXT_SEPARATOR)

			// SHELL
			// TODO: Should probably declare a style for this
			prompt.TextShell += prompt.colorizer.colorize(SYMBOL_TEXT_SEPARATOR, COLOR_TEXT_FG_SEPARATOR, "", false)
		}
	}
	prompt.appendToSegment(text, style)
}

func (prompt *promptT) appendToSegment(text string, style promptStyleT) {
	if prompt.isPowerline {
		prompt.CurrentBGColorHex = style.ColorHexBGPowerline
		appendStyle := color.HEXStyle(style.ColorHexFGPowerline, style.ColorHexBGPowerline)
		if style.Bold {
			appendStyle.SetOpts(color.Opts{color.OpBold})
		}
		prompt.TextPrintable += appendStyle.Sprintf("%s", text)

		// SHELL
		prompt.TextShell += prompt.colorizer.colorize(text, style.ColorHexFGPowerline, style.ColorHexBGPowerline, style.Bold)
	} else {
		// appendColor := color.HEX(style.ColorHexFGText)
		appendStyle := color.HEXStyle(style.ColorHexFGText)
		if style.Bold {
			appendStyle.SetOpts(color.Opts{color.OpBold})
		}
		prompt.TextPrintable += appendStyle.Sprintf("%s", text)

		// SHELL
		prompt.TextShell += prompt.colorizer.colorize(text, style.ColorHexFGText, "", style.Bold)
	}
}

func (prompt *promptT) endSegments() {
	if prompt.isPowerline {
		separatorStyle := color.HEXStyle(COLOR_BG_DEFAULT, prompt.CurrentBGColorHex)
		prompt.TextPrintable += separatorStyle.Sprint(" ")
		separatorStyle = color.HEXStyle(prompt.CurrentBGColorHex)
		prompt.TextPrintable += separatorStyle.Sprintf("%s ", SYMBOL_PL_SEPARATOR)

		// SHELL
		prompt.TextShell += " "
		prompt.TextShell += prompt.colorizer.reset()
		prompt.TextShell += prompt.colorizer.colorize(SYMBOL_PL_SEPARATOR+" ", prompt.CurrentBGColorHex, "", false)
		prompt.TextShell += prompt.colorizer.reset()
	} else {
		prompt.TextPrintable += " $ "
		// SHELL
		// TODO: The final reset should not be necessary (but is).  Is a trim() removing the final spaces somewhere?
		prompt.TextShell += prompt.colorizer.reset() + " $ " + prompt.colorizer.reset()
	}
}

func (gitInfo *gitInfoT) update(path string) {
	var cmd *exec.Cmd
	var e bytes.Buffer
	var out []byte
	var err error
	var branchName string

	cmd = exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Stderr = &e
	cmd.Dir = path
	out, err = cmd.Output()
	if err != nil {
		// This is not a git repo
		return
	}
	if strings.TrimSpace(string(out)) != "true" {
		// This is not a git repo
		return
	}

	cmd = exec.Command("git", "symbolic-ref", "HEAD")
	cmd.Stderr = &e
	cmd.Dir = path
	out, err = cmd.Output()
	if err == nil {
		branchName = strings.TrimSpace(string(out))
	}
	if branchName != "" {
		branchName = finalComponent(branchName)
	} else {
		// reference = "(other)"
		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
		var e bytes.Buffer
		cmd.Stderr = &e
		cmd.Dir = path
		out, err = cmd.Output()
		if err != nil {
			// This is not a git repo
			return
		}
		branchName = "(" + strings.TrimSpace(string(out)) + ")"
		gitInfo.IsDetached = true
	}
	if branchName == "" {
		// This is not a git repo
		return
	}
	gitInfo.IsRepo = true
	gitInfo.Branch = branchName

	// ---------------------------
	// Git Status
	// ---------------------------
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Stderr = &e
	cmd.Dir = path
	out, err = cmd.Output()
	if err == nil {
		result := string(out)
		statusIndex, statusWorking := extractPorcelainStatus(result)
		switch statusIndex {
		case " ":
			// Nothing Staged
			break
		case "?":
			gitInfo.IsUntracked = true
		default:
			gitInfo.IsStaged = true
		}
		switch statusWorking {
		case " ":
			// Nothing Modified
			break
		case "?":
			gitInfo.IsUntracked = true
		default:
			gitInfo.IsModified = true
		}
	}
	gitInfo.IsDirty = (gitInfo.IsUntracked ||
		gitInfo.IsStaged ||
		gitInfo.IsModified)
}

func extractPorcelainStatus(line string) (string, string) {
	if len(line) < 2 {
		return " ", " "
	}
	statusIndex := line[0:1]
	statusWorking := line[1:2]
	return statusIndex, statusWorking
}

func (git gitInfoT) render(isPowerline bool) string {
	var symbols map[string]string
	if isPowerline {
		symbols = SYMBOLS_POWERLINE
	} else {
		symbols = SYMBOLS_TEXT
	}

	text := ""
	indicator := "branch"
	if git.IsDetached {
		indicator = "detached"
	}
	text = symbols[indicator] + " " + git.Branch

	status := ""
	if git.IsStaged {
		status += symbols["staged"]
	}
	if git.IsModified {
		status += symbols["modified"]
	}
	if git.IsUntracked {
		status += symbols["untracked"]
	}
	if status != "" {
		text += " " + status
	}
	return text
}
