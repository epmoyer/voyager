package main

import "github.com/gookit/color"

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
	GitStatus            string
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

func (prompt *promptT) init(isPowerline bool) {
	prompt.colorizer = colorizerT{}
	prompt.isPowerline = isPowerline
}

func (prompt *promptT) addSegment(text string, style promptStyleT) {
	if prompt.isPowerline && prompt.TextPrintable != "" {
		// Powerline prompt gets a leading space
		text = " " + text
	}
	if prompt.TextPrintable == "" {
		// -------------------
		//  First segment: Start with bull-nose
		// -------------------
		if prompt.isPowerline {
			bullnoseStyle := color.HEXStyle(style.ColorHexBGPowerline)
			prompt.TextPrintable += bullnoseStyle.Sprint(SYMBOL_PL_BULLNOSE)

			// SHELL
			prompt.TextShell += " "
			prompt.TextShell += prompt.colorizer.colorize(SYMBOL_PL_BULLNOSE, style.ColorHexBGPowerline, "")
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
			prompt.TextShell += prompt.colorizer.colorize(SYMBOL_PL_SEPARATOR, prompt.CurrentBGColorHex, style.ColorHexBGPowerline)
		} else {
			separatorColor := color.HEX(COLOR_TEXT_FG_SEPARATOR)
			// prompt.Prompt += separatorColor.Sprintf(" ⟫ ")
			prompt.TextPrintable += separatorColor.Sprintf("⟫")
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
		prompt.TextShell += prompt.colorizer.colorize(text, style.ColorHexFGPowerline, style.ColorHexBGPowerline)
	} else {
		appendColor := color.HEX(style.ColorHexFGText)
		prompt.TextPrintable += appendColor.Sprintf("%s", text)
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
		prompt.TextShell += prompt.colorizer.colorize(SYMBOL_PL_SEPARATOR+" ", prompt.CurrentBGColorHex, "")
		prompt.TextShell += prompt.colorizer.reset()
	} else {
		prompt.TextPrintable += " $ "
	}
}
