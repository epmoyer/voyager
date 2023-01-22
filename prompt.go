package main

import (
	"strings"

	"github.com/gookit/color"
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

const (
	ColorMode16m = iota
	ColorMode256
	ColorMode16
	ColorModeNone
)

type promptT struct {
	TextPrintable     string
	TextShell         string
	CurrentBGColorHex string
	IsPowerLine       bool
	ColorMode         int
	Colorizer         colorizerT
}

type promptStyleT struct {
	ColorHexFGPowerline string
	ColorHexBGPowerline string
	ColorHexFGText      string

	Color256FGPowerline uint8
	Color256BGPowerline uint8
	Color256FGText      uint8

	Color16FGPowerLine color.Color
	Color16BGPowerLine color.Color
	Color16Text        color.Color

	ICSColorFGPowerline string
	ICSColorBGPowerline string
	ICSColorFGText      string

	Bold bool
}

func (prompt *promptT) init(isPowerline bool, shell string, optNoColor bool, optColor string) {
	prompt.Colorizer = colorizerT{}
	prompt.Colorizer.shell = shell
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
	if prompt.IsPowerLine && !(prompt.TextPrintable == "" && ENABLE_BULLNOSE) {
		// Powerline prompt gets a leading space
		text = " " + text
	}
	if prompt.TextPrintable == "" {
		// -------------------
		//  First segment: Start with bull-nose
		// -------------------
		if prompt.IsPowerLine && ENABLE_BULLNOSE {
			bullnoseStyle := color.HEXStyle(style.ColorHexBGPowerline)
			prompt.TextPrintable += bullnoseStyle.Sprint(SYMBOL_PL_BULLNOSE)

			// SHELL
			// prompt.TextShell += " "
			prompt.TextShell += prompt.Colorizer.colorize(SYMBOL_PL_BULLNOSE, style.ColorHexBGPowerline, "", style.Bold)
		}
	} else {
		// -------------------
		//  Add Separator
		// -------------------
		if prompt.IsPowerLine {
			separatorStyle := color.HEXStyle(COLOR_FG_DEFAULT, prompt.CurrentBGColorHex)
			prompt.TextPrintable += separatorStyle.Sprint(" ")
			separatorStyle = color.HEXStyle(prompt.CurrentBGColorHex, style.ColorHexBGPowerline)
			prompt.TextPrintable += separatorStyle.Sprintf("%s", SYMBOL_PL_SEPARATOR)

			// SHELL
			prompt.TextShell += " "
			prompt.TextShell += prompt.Colorizer.colorize(SYMBOL_PL_SEPARATOR, prompt.CurrentBGColorHex, style.ColorHexBGPowerline, style.Bold)
		} else {
			separatorColor := color.HEX(COLOR_TEXT_FG_SEPARATOR)
			// prompt.Prompt += separatorColor.Sprintf(" ⟫ ")
			prompt.TextPrintable += separatorColor.Sprintf(SYMBOL_TEXT_SEPARATOR)

			// SHELL
			// TODO: Should probably declare a style for this
			prompt.TextShell += prompt.Colorizer.colorize(SYMBOL_TEXT_SEPARATOR, COLOR_TEXT_FG_SEPARATOR, "", false)
		}
	}
	prompt.appendToSegment(text, style)
}

func (prompt *promptT) appendToSegment(text string, style promptStyleT) {
	if prompt.IsPowerLine {
		prompt.CurrentBGColorHex = style.ColorHexBGPowerline
		appendStyle := color.HEXStyle(style.ColorHexFGPowerline, style.ColorHexBGPowerline)
		if style.Bold {
			appendStyle.SetOpts(color.Opts{color.OpBold})
		}
		prompt.TextPrintable += appendStyle.Sprintf("%s", text)

		// SHELL
		prompt.TextShell += prompt.Colorizer.colorize(text, style.ColorHexFGPowerline, style.ColorHexBGPowerline, style.Bold)
	} else {
		// appendColor := color.HEX(style.ColorHexFGText)
		appendStyle := color.HEXStyle(style.ColorHexFGText)
		if style.Bold {
			appendStyle.SetOpts(color.Opts{color.OpBold})
		}
		prompt.TextPrintable += appendStyle.Sprintf("%s", text)

		// SHELL
		prompt.TextShell += prompt.Colorizer.colorize(text, style.ColorHexFGText, "", style.Bold)
	}
}

func (prompt *promptT) endSegments(promptInfo promptInfoT) {
	if prompt.IsPowerLine {
		// --------------------
		// Powerline
		// --------------------

		// PRINTABLE
		separatorStyle := color.HEXStyle(COLOR_BG_DEFAULT, prompt.CurrentBGColorHex)
		prompt.TextPrintable += separatorStyle.Sprint(" ")
		separatorStyle = color.HEXStyle(prompt.CurrentBGColorHex)
		prompt.TextPrintable += separatorStyle.Sprintf("%s ", SYMBOL_PL_SEPARATOR)

		// SHELL
		prompt.TextShell += " "
		prompt.TextShell += prompt.Colorizer.reset()
		prompt.TextShell += prompt.Colorizer.colorize(SYMBOL_PL_SEPARATOR+" ", prompt.CurrentBGColorHex, "", false)
		prompt.TextShell += prompt.Colorizer.reset()
	} else {
		// --------------------
		// Text
		// --------------------
		promptSymbol := "%"
		if prompt.Colorizer.shell == "bash" {
			promptSymbol = "$"
		}
		if promptInfo.IsRoot {
			promptSymbol = "#"
		}
		promptSymbol = " " + promptSymbol + " "

		// PRINTABLE
		if promptInfo.IsRoot {
			promptStyle := color.HEXStyle(STYLE_CONTEXT_ROOT.ColorHexFGText)
			prompt.TextPrintable += promptStyle.Sprint(promptSymbol)
		} else {
			prompt.TextPrintable += promptSymbol
		}

		// SHELL
		// Escape the % symbol for zsh
		if prompt.Colorizer.shell == "zsh" {
			promptSymbol = strings.Replace(promptSymbol, "%", "%%", -1)
		}
		prompt.TextShell += prompt.Colorizer.reset()
		if promptInfo.IsRoot {
			prompt.TextShell += prompt.Colorizer.colorize(promptSymbol, STYLE_CONTEXT_ROOT.ColorHexFGText, "", false)
		} else {
			prompt.TextShell += promptSymbol
		}
		// TODO: The final reset should not be necessary (but is).  Is a trim() removing the final spaces somewhere?
		prompt.TextShell += prompt.Colorizer.reset()
	}
}
