package main

// ANSII Terminal Escape codes
const ESC_CODE_RESET = "0"
const ESC_CODE_BOLD = "1"
const ESC_CODE_FG_RGB = "38"
const ESC_CODE_BG_RGB = "48"

type colorizerT struct {
	shell string
}

func (colorizer colorizerT) colorize(text string, colorHexFG string, colorHexBG string, bold bool) string {
	if colorizer.shell == "zsh" {
		return colorizer._colorize_zsh(text, colorHexFG, colorHexBG, bold)
	} else if colorizer.shell == "bash" {
		return colorizer._colorize_bash(text, colorHexFG, colorHexBG, bold)
	}
	return "(Unknown shell)"
}

func (colorizer colorizerT) reset() string {
	if colorizer.shell == "zsh" {
		return "%{%f%b%k%}"
	} else if colorizer.shell == "bash" {
		return bash_esc(ESC_CODE_RESET)
	}
	return "(Unknown shell)"
}

func (colorizer colorizerT) _colorize_zsh(text string, colorHexFG string, colorHexBG string, bold bool) string {
	var formatFG string
	var formatBG string
	var formatBold string
	if colorHexFG != "" {
		formatFG = "%F{" + colorHexFG + "}"
	}
	if colorHexBG != "" {
		formatBG = "%K{" + colorHexBG + "}"
	}
	if bold {
		formatBold = "%B"
	}
	result := "%{" + formatFG + formatBG + formatBold + "%}" + text + "%b"
	return result
}

func (colorizer colorizerT) _colorize_bash(text string, colorHexFG string, colorHexBG string, bold bool) string {
	var formatFG string
	var formatBG string
	var formatBold string
	if colorHexFG != "" {
		formatFG = bash_esc(ESC_CODE_FG_RGB + ";255;128;0")
	}
	if colorHexBG != "" {
		formatFG = bash_esc(ESC_CODE_BG_RGB + ";200;200;255")
	}
	if bold {
		formatBold = bash_esc(ESC_CODE_BOLD)
	}
	result := formatFG + formatBG + formatBold + text + bash_esc(ESC_CODE_RESET)
	return result
}

func bash_esc(text string) string {
	// return "[\u001b[" + text + "m]"
	return "\u001b[" + text + "m"
}
