package main

import "github.com/gookit/color"

type promptT struct {
	Prompt            string
	CurrentBGColorHex string
}

func (prompt promptT) addSegment(text string, FGColorHex string, BGColorHex string, withSeparator bool) promptT {
	if prompt.Prompt != "" && withSeparator {
		separatorStyle := color.HEXStyle(prompt.CurrentBGColorHex, BGColorHex)
		prompt.Prompt += separatorStyle.Sprintf("\ue0b0")
	}
	prompt.CurrentBGColorHex = BGColorHex
	appendStyle := color.HEXStyle(FGColorHex, BGColorHex)
	prompt.Prompt += appendStyle.Sprintf("%s", text)
	return prompt
}

func (prompt promptT) endSegments() promptT {
	if prompt.Prompt != "" {
		separatorStyle := color.HEX(prompt.CurrentBGColorHex)
		prompt.Prompt += separatorStyle.Sprintf(SYMBOL_SEPARATOR)
	}
	return prompt
}
