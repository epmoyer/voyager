package main

import "github.com/gookit/color"

type powerlinePromptT struct {
	Prompt            string
	CurrentBGColorHex string
}

type promptStyleT struct {
	ColorHexFG string
	ColorHexBG string
	Bold       bool
}

func (prompt powerlinePromptT) addSegment(text string, style promptStyleT, withSeparator bool) powerlinePromptT {
	if prompt.Prompt != "" && withSeparator {
		separatorStyle := color.HEXStyle(prompt.CurrentBGColorHex, style.ColorHexBG)
		prompt.Prompt += separatorStyle.Sprintf("\ue0b0")
	}
	prompt.CurrentBGColorHex = style.ColorHexBG
	appendStyle := color.HEXStyle(style.ColorHexFG, style.ColorHexBG)
	if style.Bold {
		appendStyle.SetOpts(color.Opts{color.OpBold})
	}
	prompt.Prompt += appendStyle.Sprintf("%s", text)
	return prompt
}

func (prompt powerlinePromptT) endSegments() powerlinePromptT {
	if prompt.Prompt != "" {
		separatorStyle := color.HEX(prompt.CurrentBGColorHex)
		prompt.Prompt += separatorStyle.Sprintf(SYMBOL_SEPARATOR)
	}
	return prompt
}
