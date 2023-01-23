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
			// prompt.PromptTextICS += colorizeICS(SYMBOL_PL_BULLNOSE, style.ICSColorBGPowerline, "", style.Bold)
			prompt.PromptTextICS += icsFormat(style.ICSColorBGPowerline, "", "") + SYMBOL_PL_BULLNOSE
		}
	} else {
		// -------------------
		//  Add Separator
		// -------------------
		if prompt.IsPowerLine {
			// separatorStyle := color.HEXStyle(COLOR_FG_DEFAULT, prompt.CurrentBGColorHex)
			// prompt.TextPrintable += separatorStyle.Sprint(" ")
			// separatorStyle = color.HEXStyle(prompt.CurrentBGColorHex, style.ColorHexBGPowerline)
			// prompt.TextPrintable += separatorStyle.Sprintf("%s", SYMBOL_PL_SEPARATOR)

			// // SHELL
			// prompt.TextShell += " "
			// prompt.TextShell += prompt.Colorizer.colorize(SYMBOL_PL_SEPARATOR, prompt.CurrentBGColorHex, style.ColorHexBGPowerline, style.Bold)

			// prompt.PromptTextICS += colorizeICS(" ", "", prompt.CurrentBGColorICS, false)
			prompt.PromptTextICS += " "
			// prompt.PromptTextICS += colorizeICS(SYMBOL_PL_SEPARATOR, prompt.CurrentBGColorICS, style.ICSColorBGPowerline, style.Bold)
			prompt.PromptTextICS += icsFormat(prompt.CurrentBGColorICS, style.ICSColorBGPowerline, "") + SYMBOL_PL_SEPARATOR
		} else {
			// separatorColor := color.HEX(COLOR_TEXT_FG_SEPARATOR)
			// // prompt.Prompt += separatorColor.Sprintf(" ⟫ ")
			// prompt.TextPrintable += separatorColor.Sprintf(SYMBOL_TEXT_SEPARATOR)

			// // SHELL
			// // TODO: Should probably declare a style for this
			// prompt.TextShell += prompt.Colorizer.colorize(SYMBOL_TEXT_SEPARATOR, COLOR_TEXT_FG_SEPARATOR, "", false)

			// prompt.PromptTextICS += colorizeICS(SYMBOL_TEXT_SEPARATOR, ICS_COLOR_TEXT_FG_SEPARATOR, "", false)
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
		prompt.CurrentBGColorICS = style.ICSColorBGPowerline
		// prompt.CurrentBGColorHex = style.ColorHexBGPowerline
		// appendStyle := color.HEXStyle(style.ColorHexFGPowerline, style.ColorHexBGPowerline)
		// if style.Bold {
		// 	appendStyle.SetOpts(color.Opts{color.OpBold})
		// }
		// prompt.TextPrintable += appendStyle.Sprintf("%s", text)

		// // SHELL
		// prompt.TextShell += prompt.Colorizer.colorize(text, style.ColorHexFGPowerline, style.ColorHexBGPowerline, style.Bold)

		// prompt.PromptTextICS += colorizeICS(text, style.ICSColorFGPowerline, style.ICSColorBGPowerline, style.Bold)
		// prompt.PromptTextICS += colorizeICS(text, style.ICSColorFGPowerline, style.ICSColorBGPowerline, style.Bold)
		prompt.PromptTextICS += icsFormat(style.ICSColorFGPowerline, style.ICSColorBGPowerline, icsBoldBoolToString(style.Bold)) + text
	} else {
		// appendColor := color.HEX(style.ColorHexFGText)
		// appendStyle := color.HEXStyle(style.ColorHexFGText)
		// if style.Bold {
		// 	appendStyle.SetOpts(color.Opts{color.OpBold})
		// }
		// prompt.TextPrintable += appendStyle.Sprintf("%s", text)

		// // SHELL
		// prompt.TextShell += prompt.Colorizer.colorize(text, style.ColorHexFGText, "", style.Bold)
		// prompt.PromptTextICS += colorizeICS(text, style.ICSColorFGPowerline, "", style.Bold)
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

		// // PRINTABLE
		// separatorStyle := color.HEXStyle(COLOR_BG_DEFAULT, prompt.CurrentBGColorHex)
		// prompt.TextPrintable += separatorStyle.Sprint(" ")
		// separatorStyle = color.HEXStyle(prompt.CurrentBGColorHex)
		// prompt.TextPrintable += separatorStyle.Sprintf("%s ", SYMBOL_PL_SEPARATOR)

		// // SHELL
		// prompt.TextShell += " "
		// prompt.TextShell += prompt.Colorizer.reset()
		// prompt.TextShell += prompt.Colorizer.colorize(SYMBOL_PL_SEPARATOR+" ", prompt.CurrentBGColorHex, "", false)
		// prompt.TextShell += prompt.Colorizer.reset()

		prompt.PromptTextICS += " "
		// prompt.PromptTextICS += colorizeICS(SYMBOL_PL_SEPARATOR+" ", prompt.CurrentBGColorICS, "", false)
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

		// // PRINTABLE
		// if promptInfo.IsRoot {
		// 	promptStyle := color.HEXStyle(STYLE_CONTEXT_ROOT.ColorHexFGText)
		// 	prompt.TextPrintable += promptStyle.Sprint(promptSymbol)
		// } else {
		// 	prompt.TextPrintable += promptSymbol
		// }

		// // SHELL
		// prompt.TextShell += prompt.Colorizer.reset()
		// if promptInfo.IsRoot {
		// 	prompt.TextShell += prompt.Colorizer.colorize(promptSymbol, STYLE_CONTEXT_ROOT.ColorHexFGText, "", false)
		// } else {
		// 	prompt.TextShell += promptSymbol
		// }
		// // TODO: The final reset should not be necessary (but is).  Is a trim() removing the final spaces somewhere?
		// prompt.TextShell += prompt.Colorizer.reset()

		if promptInfo.IsRoot {
			// prompt.PromptTextICS += colorizeICS(promptSymbol, STYLE_CONTEXT_ROOT.ICSColorFGText, "", false)
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
	// TODO: Implement ColorMode
	if optPrintable {
		display := icsRenderDisplay(prompt.PromptTextICS, ColorMode16m)
		debugDump(display)
		return display
	} else {
		return icsRenderPrompt(prompt.PromptTextICS, ColorMode16m, prompt.Shell)
	}
}

func debugDump(text string) {
	for _, character := range text {
		if character == '\033' {
			fmt.Printf("^")
		} else {
			fmt.Printf("%c", character)
		}
	}
	fmt.Println()
}
