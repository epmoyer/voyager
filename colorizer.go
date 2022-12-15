package main

type colorizerT struct {
	shell string
}

func (colorizer colorizerT) colorize(text string, colorHexFG string, colorHexBG string, bold bool) string {
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

func (colorizer colorizerT) reset() string {
	return "%{%f%b%k%}"
}
