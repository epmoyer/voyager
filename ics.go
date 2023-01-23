package main

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
)

const (
	ColorMode16m = iota
	ColorMode256
	ColorMode16
	ColorModeNone
)

var ESCAPE_RESET_FOREGROUND = "\033[39m"
var ESCAPE_RESET_BACKGROUND = "\033[49m"
var ESCAPE_BOLD = "\033[1m"
var ESCAPE_RESET_BOLD = "\033[0m"

var COLOR_CODES_16 = map[string]uint8{
	"black":         30,
	"red":           31,
	"green":         32,
	"yellow":        33,
	"blue":          34,
	"magenta":       35,
	"cyan":          36,
	"white":         37,
	"brightblack":   90,
	"brightred":     91,
	"brightgreen":   92,
	"brightyellow":  93,
	"brightblue":    94,
	"brightmagenta": 95,
	"brightcyan":    96,
	"brightwhite":   97,
}

func icsFormat(ICSColorFG string, ICSColorBG string, bold string) string {
	format := "%{"
	if ICSColorFG != "" {
		if ICSColorFG == "clear" {
			// Clear foreground color
			format += "%f"
		} else {
			// Set foreground color
			format += "%F{" + ICSColorFG + "}"
		}
	}
	if ICSColorBG != "" {
		if ICSColorBG == "clear" {
			// Clear background color
			format += "%k"
		} else {
			// Set background color
			format += "%K{" + ICSColorBG + "}"
		}
	}
	if bold != "" {
		if bold == "clear" {
			// Clear bold
			format += "%b"
		} else {
			// Set bold
			format += "%B"
		}
	}
	format += "%}"
	return format
}

// Convert boolean into a string suitable for passing to the `bold` argument
// of icsFormat().
func icsBoldBoolToString(bold bool) string {
	if bold {
		return "bold"
	}
	return "clear"
}

func icsFormatClearAll() string {
	return icsFormat("clear", "clear", "clear")
}

func icsRenderPrompt(icsText string, colorMode int, shell string) string {
	// return fmt.Sprintf("icsRenderPrompt()  %d %s", colorMode, shell)
	result := icsText
	regexFormat := regexp.MustCompile(`%{(.*?)%}`)
	formatMatches := regexFormat.FindAllStringSubmatch(icsText, -1)
	// fmt.Printf("%#v\n", formatMatches)
	for _, match := range formatMatches {
		formatRaw := match[0]
		formatPayload := match[1]
		if shell == "bash" {
			// ------------------
			// Bash
			// ------------------
			escapeText := "\\[" + icsToEscapeCodes(formatPayload, colorMode) + "\\]"
			result = strings.Replace(result, formatRaw, escapeText, 1)
		} else {
			// ------------------
			// Zsh
			// ------------------
			escapeText := icsToZshPromptCodes(formatPayload, colorMode)
			result = strings.Replace(result, formatRaw, escapeText, 1)
		}
	}
	return result
}

func icsRenderDisplay(icsText string, colorMode int) string {
	// fmt.Printf("icsRenderDisplay()  %d", colorMode)
	result := icsText
	regexFormat := regexp.MustCompile(`%{(.*?)%}`)
	formatMatches := regexFormat.FindAllStringSubmatch(icsText, -1)
	// fmt.Printf("%#v\n", formatMatches)
	for _, match := range formatMatches {
		formatRaw := match[0]
		formatPayload := match[1]
		escapeText := icsToEscapeCodes(formatPayload, colorMode)
		result = strings.Replace(result, formatRaw, escapeText, 1)
	}
	return result
}

func icsToEscapeCodes(formatPayload string, colorMode int) string {
	result := formatPayload

	// -------------------
	// Foreground
	// -------------------
	regexForeground := regexp.MustCompile(`%F{(.*?)}`)
	foregroundMatches := regexForeground.FindAllStringSubmatch(formatPayload, -1)
	// fmt.Printf("%#v\n", foregroundMatches)
	for _, match := range foregroundMatches {
		formatRaw := match[0]
		formatPayload := match[1]
		escapeText := icsColorToEscapeCodesFG(formatPayload, colorMode)
		result = strings.Replace(result, formatRaw, escapeText, 1)
	}

	// -------------------
	// Background
	// -------------------
	regexBackground := regexp.MustCompile(`%K{(.*?)}`)
	backgroundMatches := regexBackground.FindAllStringSubmatch(formatPayload, -1)
	// fmt.Printf("%#v\n", backgroundMatches)
	for _, match := range backgroundMatches {
		formatRaw := match[0]
		formatPayload := match[1]
		escapeText := icsColorToEscapeCodesBG(formatPayload, colorMode)
		result = strings.Replace(result, formatRaw, escapeText, 1)
	}

	// -------------------
	// Bold
	// -------------------
	result = strings.Replace(result, "%B", ESCAPE_BOLD, -1)

	// -------------------
	// Reset Bold
	// -------------------
	result = strings.Replace(result, "%b", ESCAPE_RESET_BOLD, -1)

	// -------------------
	// Reset Foreground
	// -------------------
	result = strings.Replace(result, "%f", ESCAPE_RESET_FOREGROUND, -1)

	// -------------------
	// Reset Background
	// -------------------
	result = strings.Replace(result, "%k", ESCAPE_RESET_BACKGROUND, -1)

	return result
}

func icsColorToEscapeCodesFG(icsColor string, colorMode int) string {
	colors := strings.Split(icsColor, ":")
	if len(colors) >= 3 && colorMode == ColorMode16m {
		// TODO: Handle error
		colorRGBA, _ := ParseHexColor(colors[2])
		return FGEscape16m(colorRGBA)
		// return
	}
	if len(colors) >= 2 && (colorMode == ColorMode256 || colorMode == ColorMode16m) {
		// TODO: Handle error
		color16, _ := strconv.Atoi(colors[1])
		return FGEscape256(uint8(color16))
	}
	if len(colors) >= 1 && colorMode != ColorModeNone {
		return FGEscape16(colors[0])
	}
	// No color
	return ""
}

func icsColorToEscapeCodesBG(icsColor string, colorMode int) string {
	colors := strings.Split(icsColor, ":")
	if len(colors) >= 3 && colorMode == ColorMode16m {
		// TODO: Handle error
		colorRGBA, _ := ParseHexColor(colors[2])
		return BGEscape16m(colorRGBA)
		// return
	}
	if len(colors) >= 2 && (colorMode == ColorMode256 || colorMode == ColorMode16m) {
		// TODO: Handle error
		color16, _ := strconv.Atoi(colors[1])
		return BGEscape256(uint8(color16))
	}
	if len(colors) >= 1 && colorMode != ColorModeNone {
		return BGEscape16(colors[0])
	}
	// No color
	return ""
}

func FGEscape16m(c color.RGBA) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

func BGEscape16m(c color.RGBA) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", c.R, c.G, c.B)
}

func FGEscape256(color256 uint8) string {
	return fmt.Sprintf("\033[38;5;%dm", color256)
}

func BGEscape256(color256 uint8) string {
	return fmt.Sprintf("\033[48;5;%dm", color256)
}

func FGEscape16(colorName16 string) string {
	colorCode16, ok := COLOR_CODES_16[colorName16]
	if ok {
		return fmt.Sprintf("\033[%dm", colorCode16)
	}

	return ""
}

func BGEscape16(colorName16 string) string {
	colorCode16, ok := COLOR_CODES_16[colorName16]
	if ok {
		return fmt.Sprintf("\033[%dm", colorCode16+10)
	}

	return ""
}

func icsToZshPromptCodes(formatPayload string, colorMode int) string {
	result := formatPayload

	// -------------------
	// Foreground
	// -------------------
	regexForeground := regexp.MustCompile(`%F{(.*?)}`)
	foregroundMatches := regexForeground.FindAllStringSubmatch(formatPayload, -1)
	// fmt.Printf("%#v\n", foregroundMatches)
	for _, match := range foregroundMatches {
		formatRaw := match[0]
		formatPayload := match[1]
		zshColor := icsColorToZshColor(formatPayload, colorMode, "F")
		result = strings.Replace(result, formatRaw, zshColor, 1)
	}

	// -------------------
	// Background
	// -------------------
	regexBackground := regexp.MustCompile(`%K{(.*?)}`)
	backgroundMatches := regexBackground.FindAllStringSubmatch(formatPayload, -1)
	// fmt.Printf("%#v\n", backgroundMatches)
	for _, match := range backgroundMatches {
		formatRaw := match[0]
		formatPayload := match[1]
		zshColor := icsColorToZshColor(formatPayload, colorMode, "K")
		result = strings.Replace(result, formatRaw, zshColor, 1)
	}

	if colorMode == ColorModeNone {
		// -------------------
		// Remove Foreground and Background reset directives
		// -------------------
		result = strings.Replace(result, "%f", "", -1)
		result = strings.Replace(result, "%k", "", -1)
	}

	if result != "" {
		result = "%{" + result + "%}"
	}
	return result
}

func icsColorToZshColor(icsColor string, colorMode int, zshPrefix string) string {
	colors := strings.Split(icsColor, ":")
	if len(colors) >= 3 && colorMode == ColorMode16m {
		return wrapZshColor(zshPrefix, colors[2])
	}
	if len(colors) >= 2 && (colorMode == ColorMode256 || colorMode == ColorMode16m) {
		return wrapZshColor(zshPrefix, colors[1])
	}
	if len(colors) >= 1 && colorMode != ColorModeNone {
		return wrapZshColor(zshPrefix, colors[0])
	}
	// No color
	return ""
}

func wrapZshColor(zshPrefix string, zshColor string) string {
	return "%" + zshPrefix + "{" + zshColor + "}"
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
