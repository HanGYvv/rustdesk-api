package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var sectionOrder = []string{
	"Added",
	"Changed",
	"Deprecated",
	"Removed",
	"Fixed",
	"Security",
}

type releaseNotesArgs struct {
	tag         string
	repoURL     string
	cwd         string
	output      string
	releaseDate string
}

func main() {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	notes, err := buildReleaseNotes(args.cwd, args.tag, args.repoURL, args.releaseDate)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if args.output == "-" {
		fmt.Println(notes)
		return
	}

	if err := os.WriteFile(args.output, []byte(notes+"\n"), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseArgs(argv []string) (releaseNotesArgs, error) {
	fs := flag.NewFlagSet("release-notes", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var args releaseNotesArgs
	fs.StringVar(&args.tag, "tag", "", "release tag")
	fs.StringVar(&args.repoURL, "repo-url", "", "repository url")
	fs.StringVar(&args.cwd, "cwd", "", "working directory")
	fs.StringVar(&args.output, "output", "-", "output file")
	fs.StringVar(&args.releaseDate, "release-date", todayUTC(), "release date")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: release-notes --tag TAG --repo-url URL [--cwd DIR] [--output FILE] [--release-date YYYY-MM-DD]")
	}

	if err := fs.Parse(argv); err != nil {
		return releaseNotesArgs{}, err
	}
	if args.tag == "" {
		return releaseNotesArgs{}, errors.New("missing required argument: --tag")
	}
	if args.repoURL == "" {
		return releaseNotesArgs{}, errors.New("missing required argument: --repo-url")
	}
	if args.cwd == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return releaseNotesArgs{}, err
		}
		args.cwd = cwd
	}
	return args, nil
}

func todayUTC() string {
	return time.Now().UTC().Format("2006-01-02")
}

func buildReleaseNotes(cwd, tag, repoURL, releaseDate string) (string, error) {
	previousTag, err := resolvePreviousTag(cwd, tag)
	if err != nil {
		return "", err
	}

	commits, err := resolveCommitSubjects(cwd, tag, previousTag)
	if err != nil {
		return "", err
	}

	return renderReleaseNotes(tag, previousTag, releaseDate, repoURL, commits), nil
}

func resolvePreviousTag(cwd, currentTag string) (string, error) {
	out, err := gitOutput(cwd, "tag", "--sort=-creatordate")
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(out, "\n") {
		tag := strings.TrimSpace(line)
		if tag != "" && tag != currentTag {
			return tag, nil
		}
	}
	return "", nil
}

func resolveCommitSubjects(cwd, tag, previousTag string) ([]string, error) {
	rng := tag
	if previousTag != "" {
		rng = previousTag + ".." + tag
	}

	out, err := gitOutput(cwd, "log", "--no-merges", "--format=%s", rng)
	if err != nil {
		return nil, err
	}

	var commits []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			commits = append(commits, line)
		}
	}
	return commits, nil
}

func gitOutput(cwd string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = filepath.Clean(cwd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git command failed: %s", strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func renderReleaseNotes(tag, previousTag, releaseDate, repoURL string, commits []string) string {
	sections := map[string][]string{}
	for _, name := range sectionOrder {
		sections[name] = nil
	}

	for _, subject := range commits {
		section := classifySubject(subject)
		cleaned := normalizeSubject(subject)
		if !contains(sections[section], cleaned) {
			sections[section] = append(sections[section], cleaned)
		}
	}

	lines := []string{
		"# Changelog",
		"",
		fmt.Sprintf("## [%s] - %s", strings.TrimPrefix(tag, "v"), releaseDate),
		"",
	}

	for _, section := range sectionOrder {
		items := sections[section]
		if len(items) == 0 {
			continue
		}
		lines = append(lines, "### "+section, "")
		for _, item := range items {
			lines = append(lines, "- "+item)
		}
		lines = append(lines, "")
	}

	if previousTag != "" {
		lines = append(lines,
			"### Full Changelog",
			"",
			"<details>",
			"<summary>Show full changelog</summary>",
			"",
			fmt.Sprintf("[Compare changes](%s/compare/%s...%s)", strings.TrimRight(repoURL, "/"), previousTag, tag),
			"",
			"</details>",
		)
	} else {
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
	}

	return strings.Join(lines, "\n")
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func normalizeSubject(subject string) string {
	subject = stripPRSuffix(subject)
	subject = stripConventionalPrefix(subject)
	return stripLeadingAction(subject)
}

func stripPRSuffix(subject string) string {
	trimmed := strings.TrimSpace(subject)
	start := strings.LastIndex(trimmed, " (#")
	if start >= 0 && strings.HasSuffix(trimmed, ")") {
		id := trimmed[start+3 : len(trimmed)-1]
		if id != "" {
			ok := true
			for _, r := range id {
				if r < '0' || r > '9' {
					ok = false
					break
				}
			}
			if ok {
				return strings.TrimSpace(trimmed[:start])
			}
		}
	}
	return trimmed
}

func stripConventionalPrefix(subject string) string {
	colon := strings.Index(subject, ":")
	if colon < 0 {
		return strings.TrimSpace(subject)
	}

	prefix := subject[:colon]
	ok := len(prefix) > 0
	for _, r := range prefix {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && r != '(' && r != ')' {
			ok = false
			break
		}
	}
	if ok {
		return strings.TrimSpace(subject[colon+1:])
	}
	return strings.TrimSpace(subject)
}

func stripLeadingAction(subject string) string {
	prefixes := []string{
		"add ", "added ", "fix ", "fixed ", "update ", "updated ",
		"remove ", "removed ", "delete ", "deleted ",
	}
	for _, prefix := range prefixes {
		if len(subject) >= len(prefix) && strings.EqualFold(subject[:len(prefix)], prefix) {
			return strings.TrimSpace(subject[len(prefix):])
		}
	}
	return strings.TrimSpace(subject)
}

func classifySubject(subject string) string {
	lower := strings.ToLower(subject)
	switch {
	case strings.Contains(lower, "security"):
		return "Security"
	case strings.HasPrefix(lower, "feat") || strings.HasPrefix(lower, "add"):
		return "Added"
	case strings.HasPrefix(lower, "fix") || strings.HasPrefix(lower, "bugfix"):
		return "Fixed"
	case strings.HasPrefix(lower, "remove") || strings.HasPrefix(lower, "delete"):
		return "Removed"
	case strings.HasPrefix(lower, "deprecate"):
		return "Deprecated"
	default:
		return "Changed"
	}
}
