package main

import (
	"fmt"
	"image/color"
)

// ANSII Terminal Escape codes
const ESC_CODE_RESET = "0"
const ESC_CODE_BOLD = "1"
const ESC_CODE_FG_RGB = "38;2"
const ESC_CODE_BG_RGB = "48;2"

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
	var err error
	var colorFG color.RGBA
	var colorBG color.RGBA

	if colorHexFG != "" {
		colorFG, err = ParseHexColor(colorHexFG)
		if err != nil {
			return "(fg color err)"
		}
		formatFG = bash_esc(ESC_CODE_FG_RGB + fmt.Sprintf(";%d;%d;%d", colorFG.R, colorFG.G, colorFG.B))
	}
	if colorHexBG != "" {
		colorBG, err = ParseHexColor(colorHexBG)
		if err != nil {
			return "(fg color err)"
		}
		formatBG = bash_esc(ESC_CODE_BG_RGB + fmt.Sprintf(";%d;%d;%d", colorBG.R, colorBG.G, colorBG.B))
	}
	if bold {
		formatBold = bash_esc(ESC_CODE_BOLD)
	}
	result := formatFG + formatBG + formatBold + text + bash_esc(ESC_CODE_RESET)
	return result
}

func bash_esc(text string) string {
	// return "[\u001b[" + text + "m]"
	// return "\[\u001b[" + text + "m\]"
	return "\\[\\e[" + text + "m\\]"
}

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}
