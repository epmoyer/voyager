package main

// ANSII Terminal Escape codes
const ESC_CODE_RESET = "0"
const ESC_CODE_RESET_BOLD_DIM = "22"
const ESC_CODE_BOLD = "1"
const ESC_CODE_FG_RGB = "38;2"
const ESC_CODE_BG_RGB = "48;2"

func colorizeICS(text string, ICSColorFG string, ICSColorBG string, bold bool) string {
	result := "%{"
	if ICSColorFG != "" {
		result += "%F{" + ICSColorFG + "}"
	}
	if ICSColorBG != "" {
		result += "%K{" + ICSColorBG + "}"
	}
	if bold {
		result += "%B"
	}
	result += "%}"
	return result
}
