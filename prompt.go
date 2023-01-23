package main

import (
	"fmt"
	"strings"
)

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
	IsRoot               bool
	ReturnValue          int
}

type promptT struct {
	PromptTextICS     string
	CurrentBGColorICS string

	IsPowerLine bool
	ColorMode   int
	Shell       string
}

type promptStyleT struct {
	ICSColorFGPowerline string
	ICSColorBGPowerline string
	ICSColorFGText      string

	Bold bool
}

func (prompt *promptT) init(isPowerline bool, shell string, optNoColor bool, optColor string) {
	prompt.Shell = shell
	prompt.IsPowerLine = isPowerline

	// --------------------
	// Set color mode
	// --------------------
	if optNoColor {
		prompt.ColorMode = ColorModeNone
		return
	}
	switch optColor {
	case "16":
		prompt.ColorMode = ColorMode16
	case "256":
		prompt.ColorMode = ColorMode256
	case "16m":
		prompt.ColorMode = ColorMode16m
	}
}

func (prompt *promptT) addSegment(text string, style promptStyleT) {
	if prompt.IsPowerLine && !(prompt.PromptTextICS == "" && ENABLE_BULLNOSE) {
		// Powerline prompt gets a leading space
		text = " " + text
	}
	if prompt.PromptTextICS == "" {
		// -------------------
		//  First segment: Start with bull-nose
		// -------------------
		if prompt.IsPowerLine && ENABLE_BULLNOSE {
			prompt.PromptTextICS += icsFormat(style.ICSColorBGPowerline, "", "") + SYMBOL_PL_BULLNOSE
		}
	} else {
		// -------------------
		//  Add Separator
		// -------------------
		if prompt.IsPowerLine {
			prompt.PromptTextICS += " "
			prompt.PromptTextICS += icsFormat(prompt.CurrentBGColorICS, style.ICSColorBGPowerline, "") + SYMBOL_PL_SEPARATOR
		} else {
			prompt.PromptTextICS += icsFormat(ICS_COLOR_TEXT_FG_SEPARATOR, "", icsBoldBoolToString(style.Bold)) + SYMBOL_TEXT_SEPARATOR
			if style.Bold {
				// Clear Bold
				prompt.PromptTextICS += icsFormat("", "", "clear")
			}
		}
	}
	prompt.appendToSegment(text, style)
}

func (prompt *promptT) appendToSegment(text string, style promptStyleT) {
	if prompt.IsPowerLine {
		// --------------------
		// Powerline
		// --------------------
		prompt.CurrentBGColorICS = style.ICSColorBGPowerline
		prompt.PromptTextICS += icsFormat(style.ICSColorFGPowerline, style.ICSColorBGPowerline, icsBoldBoolToString(style.Bold)) + text
	} else {
		// --------------------
		// Text
		// --------------------
		prompt.PromptTextICS += icsFormat(style.ICSColorFGText, "", icsBoldBoolToString(style.Bold)) + text
	}
	if style.Bold {
		// Clear Bold
		prompt.PromptTextICS += icsFormat("", "", "clear")
	}
}

func (prompt *promptT) endSegments(promptInfo promptInfoT) {
	if prompt.IsPowerLine {
		// --------------------
		// Powerline
		// --------------------
		prompt.PromptTextICS += " "
		prompt.PromptTextICS += icsFormat(prompt.CurrentBGColorICS, "clear", "clear") + SYMBOL_PL_SEPARATOR + " "
		prompt.PromptTextICS += ICS_RESET_ALL
	} else {
		// --------------------
		// Text
		// --------------------
		promptSymbol := "%"
		if prompt.Shell == "bash" {
			promptSymbol = "$"
		}
		if promptInfo.IsRoot {
			promptSymbol = "#"
		}
		promptSymbol = " " + promptSymbol + " "
		// Escape the % symbol for zsh
		if prompt.Shell == "zsh" {
			promptSymbol = strings.Replace(promptSymbol, "%", "%%", -1)
		}

		if promptInfo.IsRoot {
			prompt.PromptTextICS += icsFormat(STYLE_CONTEXT_ROOT.ICSColorFGText, "", icsBoldBoolToString(STYLE_CONTEXT_ROOT.Bold)) + promptSymbol
		} else {
			prompt.PromptTextICS += ICS_RESET_ALL
			prompt.PromptTextICS += promptSymbol
		}
		// TODO: The final reset should not be necessary (but is).  Is a trim() removing the final spaces somewhere?
		prompt.PromptTextICS += ICS_RESET_ALL
	}
}

func (prompt *promptT) render(optPrintable bool) string {
	if optPrintable {
		display := icsRenderDisplay(prompt.PromptTextICS, colorMode)
		debugDump(display)
		return display
	} else {
		return icsRenderPrompt(prompt.PromptTextICS, colorMode, prompt.Shell)
	}
}

func debugDump(text string) {
	if !DEBUG_ENABLE {
		return
	}
	for _, character := range text {
		if character == '\033' {
			fmt.Printf("^")
		} else {
			fmt.Printf("%c", character)
		}
	}
	fmt.Println()
}
