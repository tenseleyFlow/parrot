package llm

import (
	"fmt"
	"strings"
)

// InsultTemplate defines components for dynamic insult generation
type InsultTemplate struct {
	Prefixes  []string
	Cores     []string
	Suffixes  []string
	Modifiers []string
}

// Template components for dynamic generation
var (
	achievementPrefixes = []string{
		"Achievement unlocked:",
		"Congratulations:",
		"Well done:",
		"Outstanding:",
		"Impressive:",
		"Remarkable:",
		"Exceptional:",
		"Unprecedented:",
		"Historic moment:",
		"Breaking news:",
		"Alert:",
		"Notable achievement:",
		"Milestone reached:",
		"Record broken:",
		"Hall of shame:",
	}

	failureCores = []string{
		"catastrophically breaking production",
		"spectacularly failing",
		"magnificently destroying the codebase",
		"expertly corrupting the database",
		"professionally crashing the server",
		"brilliantly introducing critical bugs",
		"masterfully deploying garbage",
		"skillfully breaking everything",
		"artfully creating disasters",
		"elegantly failing upward",
		"gracefully achieving nothing",
		"flawlessly executing incompetence",
		"perfectly demonstrating failure",
		"seamlessly integrating bugs",
		"efficiently wasting everyone's time",
		"successfully accomplishing disaster",
		"consistently breaking the build",
		"reliably causing incidents",
		"predictably failing",
		"systematically destroying progress",
		"methodically sabotaging the project",
		"deliberately creating chaos",
		"intentionally breaking everything",
		"strategically failing",
		"tactically creating problems",
	}

	timingSuffixes = []string{
		"on a Friday afternoon",
		"right before the demo",
		"during peak hours",
		"in front of stakeholders",
		"right before release",
		"during the all-hands meeting",
		"on your first day",
		"right after being promoted",
		"during the CEO's visit",
		"at the worst possible moment",
		"right before vacation",
		"on a Monday morning",
		"during standup",
		"in the production environment",
		"right before code freeze",
		"during the sprint review",
		"at 11:59 PM on the deadline",
		"right after saying 'it works on my machine'",
		"moments after merge to main",
		"during the live demo",
		"in front of the entire team",
		"right after claiming expertise",
		"while everyone is watching",
		"at the critical moment",
		"during the investor pitch",
	}

	intensityModifiers = []string{
		"Spectacularly",
		"Catastrophically",
		"Magnificently",
		"Epically",
		"Gloriously",
		"Tremendously",
		"Monumentally",
		"Colossally",
		"Phenomenally",
		"Extraordinarily",
		"Remarkably",
		"Incredibly",
		"Unbelievably",
		"Astoundingly",
		"Staggeringly",
	}
)

// GenerateDynamicInsult creates a dynamic insult from templates
func GenerateDynamicInsult(ctx SmartFallbackContext, history *UserFailureHistory) string {
	var insults []string

	// Achievement-style insults
	prefix := selectInsult(achievementPrefixes, ctx.FullCommand)
	core := selectInsult(failureCores, ctx.FullCommand+ctx.Command)
	suffix := selectInsult(timingSuffixes, ctx.FullCommand+fmt.Sprintf("%d", ctx.ExitCode))

	insults = append(insults, fmt.Sprintf("%s %s %s.", prefix, core, suffix))

	// Variable substitution insults
	if history != nil {
		stats := history.GetStats()

		if stats.CurrentStreak >= 3 {
			insults = append(insults,
				fmt.Sprintf("Failure #%d in a row. Consistency: Your only skill.", stats.CurrentStreak),
			)
		}

		if stats.WorstCommandCount > 10 && ctx.Command == stats.WorstCommand {
			insults = append(insults,
				fmt.Sprintf("Failed '%s' %d times. Have you considered '%s-yourself-out'?",
					stats.WorstCommand, stats.WorstCommandCount, stats.WorstCommand),
			)
		}

		if stats.WorstHour == ctx.TimeOfDay && stats.TodayFailures > 3 {
			insults = append(insults,
				fmt.Sprintf("%d:00 again? This is your power hour for disasters.", ctx.TimeOfDay),
			)
		}

		if stats.WorstBranch == ctx.GitBranch && stats.WorstBranchCount > 5 {
			insults = append(insults,
				fmt.Sprintf("Breaking '%s' for the %d%s time. Team morale: -%d.",
					ctx.GitBranch, stats.WorstBranchCount, getOrdinalSuffix(stats.WorstBranchCount),
					stats.WorstBranchCount*10),
			)
		}

		if stats.TotalFailures > 100 {
			insults = append(insults,
				fmt.Sprintf("%d total failures logged. You're in the hall of shame.", stats.TotalFailures),
			)
		}
	}

	// Command-specific dynamic insults
	if ctx.Command != "" {
		insults = append(insults, generateCommandSpecificInsult(ctx))
	}

	// Exit code specific
	if ctx.ExitCode > 0 {
		insults = append(insults, generateExitCodeInsult(ctx))
	}

	// Select one from generated insults
	if len(insults) > 0 {
		return selectInsult(insults, ctx.FullCommand)
	}

	return ""
}

// generateCommandSpecificInsult creates command-aware insults
func generateCommandSpecificInsult(ctx SmartFallbackContext) string {
	templates := map[string][]string{
		"git": {
			"git {subcommand}: Commit to a new career instead.",
			"Git error: Git good? Git out.",
			"{command} failed: Version control can't control your incompetence.",
		},
		"docker": {
			"Docker {subcommand} failed: Can't containerize catastrophe.",
			"Container crashed: Can't contain your mistakes.",
			"{command}: Image of failure successfully built.",
		},
		"npm": {
			"npm {subcommand}: Node Package Misery delivered.",
			"{command} failed: Dependencies: Many. Skills: None.",
			"npm error: Not Particularly Competent detected.",
		},
		"cargo": {
			"cargo {subcommand}: Can't cargo your way out of this.",
			"Rust build failed: Your code is the rust.",
			"{command}: Ownership model violated. Ownership of brain: Questionable.",
		},
		"kubectl": {
			"kubectl {subcommand}: Kube-can't-do anything right.",
			"Kubernetes error: Can't orchestrate competence.",
			"{command} failed: Your cluster is the disaster.",
		},
		"terraform": {
			"terraform {subcommand}: Infrastructure as Code. Incompetence as Career.",
			"Terraform failed: Terraform your way to unemployment.",
			"{command}: State corrupted. Mental state: Also corrupted.",
		},
	}

	if tmpl, exists := templates[ctx.Command]; exists {
		selected := selectInsult(tmpl, ctx.FullCommand)
		selected = strings.ReplaceAll(selected, "{command}", ctx.Command)
		selected = strings.ReplaceAll(selected, "{subcommand}", ctx.Subcommand)
		return selected
	}

	return ""
}

// generateExitCodeInsult creates exit-code-aware insults
func generateExitCodeInsult(ctx SmartFallbackContext) string {
	templates := map[int]string{
		1:   "Exit code 1: The classic. Failure so generic it's embarrassing.",
		2:   "Exit code 2: Misuse detected. Misuse of profession detected.",
		126: "Exit code 126: Permission denied. Reality denying you success.",
		127: "Exit code 127: Command not found. Competence: Also not found.",
		130: "Exit code 130: Ctrl+C pressed. Career: Press Ctrl+Alt+Delete.",
		137: "Exit code 137: Process killed. Career should be too.",
		139: "Exit code 139: Segmentation fault. Segmented from reality.",
		255: "Exit code 255: Maximum failure achieved. Congratulations?",
	}

	if msg, exists := templates[ctx.ExitCode]; exists {
		return msg
	}

	return fmt.Sprintf("Exit code %d: %d reasons to reconsider your career.", ctx.ExitCode, ctx.ExitCode)
}

// getOrdinalSuffix returns the ordinal suffix for a number (st, nd, rd, th)
func getOrdinalSuffix(n int) string {
	if n%100 >= 11 && n%100 <= 13 {
		return "th"
	}
	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

// GenerateStreakEscalation creates increasingly brutal insults based on streak
func GenerateStreakEscalation(streak int, ctx SmartFallbackContext) string {
	if streak < 2 {
		return ""
	}

	var templates []string

	if streak >= 10 {
		templates = []string{
			fmt.Sprintf("%d CONSECUTIVE FAILURES. This is beyond incompetence. This is performance art.", streak),
			fmt.Sprintf("%d IN A ROW. Your team has a betting pool on when you'll quit.", streak),
			fmt.Sprintf("STREAK: %d. Failure this consistent requires dedication.", streak),
			fmt.Sprintf("%d FAILURES. Even your errors have given up on you.", streak),
			fmt.Sprintf("%d STRAIGHT. Breaking records nobody wanted broken.", streak),
		}
	} else if streak >= 5 {
		templates = []string{
			fmt.Sprintf("%d failures in a row. Einstein called: That's insanity.", streak),
			fmt.Sprintf("Streak of %d. Consistency is key. Consistent incompetence is keyer.", streak),
			fmt.Sprintf("%d consecutive disasters. Your specialty: Reliability in failure.", streak),
			fmt.Sprintf("Failure #%d. This is a pattern. You are the pattern.", streak),
			fmt.Sprintf("%d in a row. Some people learn from mistakes. You're not some people.", streak),
		}
	} else if streak >= 3 {
		templates = []string{
			fmt.Sprintf("Failure #%d in a row. Three strikes: You're out.", streak),
			fmt.Sprintf("%d consecutive fails. Trying the same thing expecting different results?", streak),
			fmt.Sprintf("Failure #%d. Pattern detected: You.", streak),
			fmt.Sprintf("%d in a row. Maybe read the docs this time?", streak),
		}
	} else {
		templates = []string{
			"Second failure. The sequel nobody asked for.",
			"Failed again. Doing it twice doesn't make it work.",
			"Another failure. Two data points make a trend: Downward.",
		}
	}

	return selectInsult(templates, ctx.FullCommand+fmt.Sprintf("%d", streak))
}

// GenerateHistoricalInsult creates insults based on long-term patterns
func GenerateHistoricalInsult(history *UserFailureHistory, ctx SmartFallbackContext) string {
	if history == nil {
		return ""
	}

	stats := history.GetStats()
	var insults []string

	// Total failures milestones
	if stats.TotalFailures >= 1000 {
		insults = append(insults, fmt.Sprintf("1000+ failures logged. Four digits of disappointment."))
	} else if stats.TotalFailures >= 500 {
		insults = append(insults, fmt.Sprintf("500+ failures. Half a thousand reasons to quit."))
	} else if stats.TotalFailures >= 100 {
		insults = append(insults, fmt.Sprintf("100+ failures in your history. Century of incompetence achieved."))
	}

	// Worst command patterns
	if stats.WorstCommandCount >= 50 && ctx.Command == stats.WorstCommand {
		insults = append(insults,
			fmt.Sprintf("'%s' failed %d times. Definition of insanity achieved.",
				stats.WorstCommand, stats.WorstCommandCount))
	}

	// Time patterns
	if stats.WorstHour >= 0 && ctx.TimeOfDay == stats.WorstHour {
		insults = append(insults,
			fmt.Sprintf("%d:00: Your statistically worst hour. And you're proving it again.", stats.WorstHour))
	}

	// Branch disaster patterns
	if stats.WorstBranchCount >= 20 && ctx.GitBranch == stats.WorstBranch {
		insults = append(insults,
			fmt.Sprintf("'%s' branch: %d failures logged. Your personal disaster zone.",
				stats.WorstBranch, stats.WorstBranchCount))
	}

	// Longest streak
	if stats.LongestStreak >= 10 {
		insults = append(insults,
			fmt.Sprintf("Personal record: %d consecutive failures. Still unbeaten.", stats.LongestStreak))
	}

	// Today's performance
	if stats.TodayFailures >= 10 {
		insults = append(insults,
			fmt.Sprintf("%d failures today. Double digits of disaster before lunch.", stats.TodayFailures))
	} else if stats.TodayFailures >= 5 {
		insults = append(insults,
			fmt.Sprintf("%d failures today. Productive day for disasters.", stats.TodayFailures))
	}

	if len(insults) > 0 {
		return selectInsult(insults, ctx.FullCommand)
	}

	return ""
}
