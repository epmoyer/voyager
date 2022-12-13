package main

import "github.com/gookit/color"

// type powerlinePromptT struct {
// 	Prompt            string
// 	CurrentBGColorHex string
// }

type textPromptT struct {
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

func (prompt textPromptT) addSegment(text string, style promptStyleT, withSeparator bool) textPromptT {
	if prompt.Prompt != "" && withSeparator {
		if prompt.isPowerline {
			separatorStyle := color.HEXStyle(prompt.CurrentBGColorHex, style.ColorHexBGPowerline)
			prompt.Prompt += separatorStyle.Sprintf("\ue0b0")
		} else {
			separatorColor := color.HEX(COLOR_TEXT_FG_SEPARATOR)
			prompt.Prompt += separatorColor.Sprintf("⟫")
		}
	}

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
	return prompt
}

func (prompt textPromptT) endSegments() textPromptT {
	if prompt.isPowerline {
		separatorStyle := color.HEX(prompt.CurrentBGColorHex)
		prompt.Prompt += separatorStyle.Sprintf(SYMBOL_SEPARATOR)
	} else {
		prompt.Prompt += "$ "
	}
	return prompt
}

// func (prompt textPromptT) addSegment(text string, colorHexFG string, withSeparator bool) textPromptT {
// 	if prompt.Prompt != "" && withSeparator {
// 		separatorColor := color.HEX(COLOR_TEXT_FG_SEPARATOR)
// 		prompt.Prompt += separatorColor.Sprintf("⟫")
// 	}
// 	prompt.CurrentBGColorHex = colorHexFG
// 	appendColor := color.HEX(colorHexFG)
// 	prompt.Prompt += appendColor.Sprintf("%s", text)
// 	return prompt
// }

// func (prompt textPromptT) endSegments() textPromptT {
// 	prompt.Prompt += "$ "
// 	return prompt
// }
