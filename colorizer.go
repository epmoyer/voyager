package main

type colorizerT struct {
	shell string
}

func (colorizer colorizerT) colorize(text string, colorHexFG string, colorHexBG string) string {
	var formatFG string
	var formatBG string
	if colorHexFG != "" {
		formatFG = "%F{" + colorHexFG + "}"
	}
	if colorHexBG != "" {
		formatBG = "%K{" + colorHexBG + "}"
	}
	result := "%{" + formatFG + formatBG + "%}" + text
	return result
}

func (colorizer colorizerT) reset() string {
	return "%{%f%b%k%}"
}
