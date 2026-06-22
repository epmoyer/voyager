package main

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

type gitInfoT struct {
	IsRepo      bool
	Branch      string
	IsDetached  bool
	IsDirty     bool
	IsStaged    bool
	IsModified  bool
	IsUntracked bool
	HasUnpushed bool // Local commits not on the upstream (ahead)
	HasUnpulled bool // Upstream commits not local (behind)
}

func (gitInfo *gitInfoT) update(path string) {
	var cmd *exec.Cmd
	var e bytes.Buffer
	var out []byte
	var err error
	var branchName string

	cmd = exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Stderr = &e
	cmd.Dir = path
	out, err = cmd.Output()
	if err != nil {
		// This is not a git repo
		return
	}
	if strings.TrimSpace(string(out)) != "true" {
		// This is not a git repo
		return
	}

	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Stderr = &e
	cmd.Dir = path
	out, err = cmd.Output()
	if err == nil {
		branchName = strings.TrimSpace(string(out))
	}
	if branchName == "HEAD" {
		// This is a detached HEAD
		// Get the commit hash instead
		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
		var e bytes.Buffer
		cmd.Stderr = &e
		cmd.Dir = path
		out, err = cmd.Output()
		if err != nil {
			// This is not a git repo
			return
		}
		branchName = "(" + strings.TrimSpace(string(out)) + ")"
		gitInfo.IsDetached = true
	}
	if branchName == "" {
		// This is not a git repo
		return
	}
	gitInfo.IsRepo = true
	gitInfo.Branch = branchName

	// ---------------------------
	// Git Status
	// ---------------------------
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Stderr = &e
	cmd.Dir = path
	out, err = cmd.Output()
	if err == nil {
		result := string(out)
		statusIndex, statusWorking := extractPorcelainStatus(result)
		switch statusIndex {
		case " ":
			// Nothing Staged
			break
		case "?":
			gitInfo.IsUntracked = true
		default:
			gitInfo.IsStaged = true
		}
		switch statusWorking {
		case " ":
			// Nothing Modified
			break
		case "?":
			gitInfo.IsUntracked = true
		default:
			gitInfo.IsModified = true
		}
	}
	gitInfo.IsDirty = (gitInfo.IsUntracked ||
		gitInfo.IsStaged ||
		gitInfo.IsModified)

	// ---------------------------
	// Upstream divergence (ahead/behind)
	// ---------------------------
	// Counts commits the upstream has that we don't (behind/unpulled) and
	// commits we have that the upstream doesn't (ahead/unpushed). If there's
	// no upstream configured (or detached HEAD), this command errors and we
	// simply leave both indicators off.
	cmd = exec.Command("git", "rev-list", "--count", "--left-right", "@{upstream}...HEAD")
	cmd.Stderr = &e
	cmd.Dir = path
	out, err = cmd.Output()
	if err == nil {
		behind, ahead := extractDivergence(string(out))
		gitInfo.HasUnpulled = behind > 0
		gitInfo.HasUnpushed = ahead > 0
	}
}

func extractDivergence(out string) (int, int) {
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return 0, 0
	}
	behind, _ := strconv.Atoi(fields[0])
	ahead, _ := strconv.Atoi(fields[1])
	return behind, ahead
}

func extractPorcelainStatus(line string) (string, string) {
	if len(line) < 2 {
		return " ", " "
	}
	statusIndex := line[0:1]
	statusWorking := line[1:2]
	return statusIndex, statusWorking
}

func (git gitInfoT) render(isPowerline bool) string {
	var symbols map[string]string
	if isPowerline {
		symbols = SYMBOLS_POWERLINE
	} else {
		symbols = SYMBOLS_TEXT
	}

	text := ""
	indicator := "branch"
	if git.IsDetached {
		indicator = "detached"
	}
	text = symbols[indicator] + " " + git.Branch

	// Upstream divergence indicator, placed at the end of the branch name
	// (bobthefish-style): unpushed (+), unpulled (-), or both (±).
	divergence := ""
	switch {
	case git.HasUnpushed && git.HasUnpulled:
		divergence = symbols["diverged"]
	case git.HasUnpushed:
		divergence = symbols["unpushed"]
	case git.HasUnpulled:
		divergence = symbols["unpulled"]
	}
	if divergence != "" {
		text += " " + divergence
	}

	status := ""
	if git.IsStaged {
		status += symbols["staged"]
	}
	if git.IsModified {
		status += symbols["modified"]
	}
	if git.IsUntracked {
		status += symbols["untracked"]
	}
	if status != "" {
		text += " " + status
	}
	return text
}
