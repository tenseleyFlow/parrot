package llm

import (
	"regexp"
	"strings"
	"time"
)

// CommandIntent represents what the user was trying to accomplish
type CommandIntent struct {
	PrimaryIntent   string   // Main action: push, build, test, deploy, etc.
	SecondaryIntents []string // Additional actions detected
	Targets         []string // What they were acting on: files, branches, containers
	Complexity      string   // simple, moderate, complex
	RiskLevel       string   // low, medium, high (based on destructiveness)
}

// IntentParser extracts semantic meaning from commands
type IntentParser struct {
	intentPatterns map[string]*regexp.Regexp
	complexityIndicators []string
	highRiskPatterns []*regexp.Regexp
}

// NewIntentParser creates a new intent parser
func NewIntentParser() *IntentParser {
	ip := &IntentParser{
		intentPatterns: make(map[string]*regexp.Regexp),
		complexityIndicators: []string{
			"|", "&&", "||", "xargs", "awk", "sed",
		},
	}
	ip.initializePatterns()
	return ip
}

func (ip *IntentParser) initializePatterns() {
	// Intent patterns
	ip.intentPatterns["push"] = regexp.MustCompile(`(git\s+push|docker\s+push|npm\s+publish)`)
	ip.intentPatterns["pull"] = regexp.MustCompile(`(git\s+pull|docker\s+pull)`)
	ip.intentPatterns["commit"] = regexp.MustCompile(`git\s+commit`)
	ip.intentPatterns["merge"] = regexp.MustCompile(`git\s+(merge|rebase)`)
	ip.intentPatterns["build"] = regexp.MustCompile(`(make|cargo\s+build|npm\s+run\s+build|go\s+build|mvn\s+compile|gradle\s+build|docker\s+build)`)
	ip.intentPatterns["test"] = regexp.MustCompile(`(test|jest|pytest|cargo\s+test|go\s+test|npm\s+test|mvn\s+test)`)
	ip.intentPatterns["deploy"] = regexp.MustCompile(`(deploy|kubectl\s+apply|helm\s+install|terraform\s+apply)`)
	ip.intentPatterns["install"] = regexp.MustCompile(`(install|npm\s+i|pip\s+install|cargo\s+add|go\s+get|apt\s+install|brew\s+install)`)
	ip.intentPatterns["configure"] = regexp.MustCompile(`(config|configure|setup|init)`)
	ip.intentPatterns["debug"] = regexp.MustCompile(`(debug|gdb|lldb|strace|lsof)`)
	ip.intentPatterns["refactor"] = regexp.MustCompile(`(refactor|rename|move\s+.*\.(js|ts|py|rs|go))`)
	ip.intentPatterns["lint"] = regexp.MustCompile(`(lint|eslint|pylint|clippy|golint|rubocop|prettier)`)
	ip.intentPatterns["format"] = regexp.MustCompile(`(format|prettier|black|rustfmt|gofmt)`)
	ip.intentPatterns["start"] = regexp.MustCompile(`(start|run|serve|up)`)
	ip.intentPatterns["stop"] = regexp.MustCompile(`(stop|kill|down)`)
	ip.intentPatterns["clean"] = regexp.MustCompile(`(clean|prune|rm|remove)`)
	ip.intentPatterns["revert"] = regexp.MustCompile(`git\s+(revert|reset|checkout)`)

	// High risk patterns
	ip.highRiskPatterns = []*regexp.Regexp{
		regexp.MustCompile(`--force`),
		regexp.MustCompile(`-f\s`),
		regexp.MustCompile(`rm\s+-rf`),
		regexp.MustCompile(`git\s+reset\s+--hard`),
		regexp.MustCompile(`drop\s+(database|table)`),
		regexp.MustCompile(`kubectl\s+delete`),
		regexp.MustCompile(`terraform\s+destroy`),
		regexp.MustCompile(`docker\s+system\s+prune`),
		regexp.MustCompile(`sudo\s+rm`),
		regexp.MustCompile(`chmod\s+777`),
	}
}

// ParseIntent analyzes a command to extract intent
func (ip *IntentParser) ParseIntent(command string) CommandIntent {
	intent := CommandIntent{
		PrimaryIntent: "unknown",
		SecondaryIntents: []string{},
		Targets: []string{},
		Complexity: "simple",
		RiskLevel: "low",
	}

	commandLower := strings.ToLower(command)

	// Detect intents
	intentsFound := make(map[string]bool)
	for intentName, pattern := range ip.intentPatterns {
		if pattern.MatchString(commandLower) {
			intentsFound[intentName] = true
		}
	}

	// Set primary and secondary intents
	if len(intentsFound) > 0 {
		// Priority order for primary intent
		priorityOrder := []string{
			"deploy", "push", "build", "test", "merge",
			"commit", "install", "configure", "start",
			"stop", "clean", "revert", "debug",
		}

		for _, priority := range priorityOrder {
			if intentsFound[priority] {
				intent.PrimaryIntent = priority
				delete(intentsFound, priority)
				break
			}
		}

		// Remaining are secondary
		for secondary := range intentsFound {
			intent.SecondaryIntents = append(intent.SecondaryIntents, secondary)
		}
	}

	// Extract targets (file names, branch names, etc.)
	intent.Targets = ip.extractTargets(command)

	// Determine complexity
	intent.Complexity = ip.assessComplexity(command)

	// Determine risk level
	intent.RiskLevel = ip.assessRisk(command)

	return intent
}

func (ip *IntentParser) extractTargets(command string) []string {
	targets := []string{}

	// Extract git branches
	branchPattern := regexp.MustCompile(`(origin/|refs/heads/)?(main|master|develop|feature/[\w-]+|bugfix/[\w-]+)`)
	if matches := branchPattern.FindAllString(command, -1); len(matches) > 0 {
		targets = append(targets, matches...)
	}

	// Extract file paths
	filePattern := regexp.MustCompile(`[\w/.-]+\.(js|ts|py|rs|go|java|cpp|c|rb|php|json|yaml|yml|toml|md)`)
	if matches := filePattern.FindAllString(command, -1); len(matches) > 0 {
		targets = append(targets, matches...)
	}

	// Extract container/image names
	containerPattern := regexp.MustCompile(`[\w.-]+:[\w.-]+`)
	if matches := containerPattern.FindAllString(command, -1); len(matches) > 0 {
		targets = append(targets, matches...)
	}

	return targets
}

func (ip *IntentParser) assessComplexity(command string) string {
	// Check for complexity indicators
	complexityScore := 0

	for _, indicator := range ip.complexityIndicators {
		if strings.Contains(command, indicator) {
			complexityScore++
		}
	}

	// Length also indicates complexity
	if len(command) > 100 {
		complexityScore++
	}

	// Multiple commands
	if strings.Contains(command, "&&") || strings.Contains(command, ";") {
		complexityScore++
	}

	if complexityScore == 0 {
		return "simple"
	} else if complexityScore <= 2 {
		return "moderate"
	}
	return "complex"
}

func (ip *IntentParser) assessRisk(command string) string {
	for _, pattern := range ip.highRiskPatterns {
		if pattern.MatchString(command) {
			return "high"
		}
	}

	// Moderate risk patterns
	moderateRiskKeywords := []string{
		"delete", "remove", "drop", "destroy",
		"force", "hard", "production", "master", "main",
	}

	commandLower := strings.ToLower(command)
	for _, keyword := range moderateRiskKeywords {
		if strings.Contains(commandLower, keyword) {
			return "medium"
		}
	}

	return "low"
}

// ContextualTags generates semantic tags based on context
func ContextualTags(ctx *SmartFallbackContext, intent CommandIntent) []InsultTag {
	tags := []InsultTag{}

	// Add error-based tags
	if ctx.ErrorPattern != "" {
		switch ctx.ErrorPattern {
		case "permission_denied":
			tags = append(tags, TagPermission)
		case "merge_conflict":
			tags = append(tags, TagMergeConflict)
		case "syntax_error":
			tags = append(tags, TagSyntax)
		case "network_error":
			tags = append(tags, TagNetwork)
		case "timeout":
			tags = append(tags, TagTimeout)
		}
	}

	// Add command type tags
	switch ctx.CommandType {
	case "git":
		tags = append(tags, TagGit)
	case "docker":
		tags = append(tags, TagDocker)
	case "kubernetes":
		tags = append(tags, TagKubernetes)
	case "nodejs":
		tags = append(tags, TagNode)
	case "python":
		tags = append(tags, TagPython)
	case "rust":
		tags = append(tags, TagRust)
	case "golang":
		tags = append(tags, TagGolang)
	}

	// Add intent tags
	switch intent.PrimaryIntent {
	case "push":
		tags = append(tags, TagPush)
	case "pull":
		tags = append(tags, TagPull)
	case "commit":
		tags = append(tags, TagCommit)
	case "build":
		tags = append(tags, TagBuild)
	case "test":
		tags = append(tags, TagTest)
	case "deploy":
		tags = append(tags, TagDeploy)
	case "install":
		tags = append(tags, TagInstall)
	case "revert":
		tags = append(tags, TagRevert)
	}

	// Add time-based tags
	hour := ctx.TimeOfDay
	if hour >= 22 || hour <= 4 {
		tags = append(tags, TagLateNight)
	} else if hour >= 5 && hour <= 7 {
		tags = append(tags, TagEarlyMorning)
	}

	// Weekend check
	if time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday {
		tags = append(tags, TagWeekend)
	}

	// Add context tags
	if ctx.IsCI {
		tags = append(tags, TagCI)
	}

	if ctx.GitBranch == "main" || ctx.GitBranch == "master" {
		tags = append(tags, TagMainBranch)
	}

	if ctx.IsRepeatedFailure {
		tags = append(tags, TagRepeated)
	}

	// Complexity tags
	switch intent.Complexity {
	case "simple":
		tags = append(tags, TagSimple)
	case "complex":
		tags = append(tags, TagComplex)
	}

	// Risk tags
	if intent.RiskLevel == "high" {
		tags = append(tags, TagProduction)
	}

	return tags
}
