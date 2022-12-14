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
}

type promptT struct {
	Prompt            string
	CurrentBGColorHex string
	isPowerline       bool
}

type promptStyleT struct {
	ColorHexFGPowerline string
	ColorHexBGPowerline string
	ColorHexFGText      string
	Bold                bool
}

func (prompt *promptT) addSegment(text string, style promptStyleT) {
	if prompt.isPowerline {
		// Powerline prompt gets a leading space
		text = " " + text
	}
	if prompt.Prompt != "" {
		// -------------------
		//  Add Separator
		// -------------------
		if prompt.isPowerline {
			separatorStyle := color.HEXStyle(style.ColorHexBGPowerline, prompt.CurrentBGColorHex)
			prompt.Prompt += separatorStyle.Sprint(" ")
			separatorStyle = color.HEXStyle(prompt.CurrentBGColorHex, style.ColorHexBGPowerline)
			prompt.Prompt += separatorStyle.Sprintf("%s", SYMBOL_PL_SEPARATOR)
		} else {
			separatorColor := color.HEX(COLOR_TEXT_FG_SEPARATOR)
			// prompt.Prompt += separatorColor.Sprintf(" ⟫ ")
			prompt.Prompt += separatorColor.Sprintf("⟫")
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
		prompt.Prompt += appendStyle.Sprintf("%s", text)
	} else {
		appendColor := color.HEX(style.ColorHexFGText)
		prompt.Prompt += appendColor.Sprintf("%s", text)
	}
}

func (prompt *promptT) endSegments() {
	if prompt.isPowerline {
		separatorStyle := color.HEXStyle(COLOR_BG_DEFAULT, prompt.CurrentBGColorHex)
		prompt.Prompt += separatorStyle.Sprint(" ")
		separatorStyle = color.HEXStyle(prompt.CurrentBGColorHex)
		prompt.Prompt += separatorStyle.Sprintf("%s ", SYMBOL_PL_SEPARATOR)
	} else {
		prompt.Prompt += " $ "
	}
}
