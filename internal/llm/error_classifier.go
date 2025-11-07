package llm

import (
	"regexp"
	"strings"
)

// ErrorCategory represents the type of error that occurred
type ErrorCategory string

const (
	ErrorPermission      ErrorCategory = "permission"
	ErrorSyntax          ErrorCategory = "syntax"
	ErrorNetwork         ErrorCategory = "network"
	ErrorDependency      ErrorCategory = "dependency"
	ErrorConfiguration   ErrorCategory = "config"
	ErrorMergeConflict   ErrorCategory = "merge_conflict"
	ErrorTestFailure     ErrorCategory = "test_failure"
	ErrorBuildFailure    ErrorCategory = "build_failure"
	ErrorTimeout         ErrorCategory = "timeout"
	ErrorNotFound        ErrorCategory = "not_found"
	ErrorAuthentication  ErrorCategory = "authentication"
	ErrorDiskSpace       ErrorCategory = "disk_space"
	ErrorMemory          ErrorCategory = "memory"
	ErrorSegfault        ErrorCategory = "segfault"
	ErrorRaceCondition   ErrorCategory = "race_condition"
	ErrorDeprecated      ErrorCategory = "deprecated"
	ErrorLinting         ErrorCategory = "linting"
	ErrorTypeMismatch    ErrorCategory = "type_mismatch"
	ErrorNullPointer     ErrorCategory = "null_pointer"
	ErrorInfiniteLoop    ErrorCategory = "infinite_loop"
	ErrorGeneric         ErrorCategory = "generic"
)

// ErrorClassifier analyzes commands and exit codes to determine error types
type ErrorClassifier struct {
	exitCodePatterns map[int][]ErrorCategory
	commandPatterns  map[*regexp.Regexp]ErrorCategory
	keywordPatterns  map[string]ErrorCategory
}

// NewErrorClassifier creates a new error classifier
func NewErrorClassifier() *ErrorClassifier {
	ec := &ErrorClassifier{
		exitCodePatterns: make(map[int][]ErrorCategory),
		commandPatterns:  make(map[*regexp.Regexp]ErrorCategory),
		keywordPatterns:  make(map[string]ErrorCategory),
	}
	ec.initializePatterns()
	return ec
}

func (ec *ErrorClassifier) initializePatterns() {
	// Exit code mappings (common Unix exit codes)
	ec.exitCodePatterns[1] = []ErrorCategory{ErrorGeneric}
	ec.exitCodePatterns[2] = []ErrorCategory{ErrorSyntax}
	ec.exitCodePatterns[126] = []ErrorCategory{ErrorPermission}
	ec.exitCodePatterns[127] = []ErrorCategory{ErrorNotFound}
	ec.exitCodePatterns[128] = []ErrorCategory{ErrorSegfault}
	ec.exitCodePatterns[130] = []ErrorCategory{ErrorTimeout}
	ec.exitCodePatterns[137] = []ErrorCategory{ErrorMemory} // SIGKILL
	ec.exitCodePatterns[139] = []ErrorCategory{ErrorSegfault} // SIGSEGV
	ec.exitCodePatterns[143] = []ErrorCategory{ErrorTimeout} // SIGTERM

	// Command pattern mappings
	ec.commandPatterns[regexp.MustCompile(`chmod|chown|sudo|permission`)] = ErrorPermission
	ec.commandPatterns[regexp.MustCompile(`git\s+(merge|rebase|cherry-pick)`)] = ErrorMergeConflict
	ec.commandPatterns[regexp.MustCompile(`(npm|yarn|pip|cargo|go)\s+(install|add)`)] = ErrorDependency
	ec.commandPatterns[regexp.MustCompile(`(pytest|jest|cargo test|go test|npm test)`)] = ErrorTestFailure
	ec.commandPatterns[regexp.MustCompile(`(make|cargo build|npm run build|go build|mvn compile)`)] = ErrorBuildFailure
	ec.commandPatterns[regexp.MustCompile(`curl|wget|fetch|ping|ssh|scp`)] = ErrorNetwork
	ec.commandPatterns[regexp.MustCompile(`docker login|git push.*http|auth|token`)] = ErrorAuthentication
	ec.commandPatterns[regexp.MustCompile(`eslint|pylint|clippy|golint|rubocop`)] = ErrorLinting

	// Keyword patterns (for error messages/command text analysis)
	ec.keywordPatterns["permission denied"] = ErrorPermission
	ec.keywordPatterns["eacces"] = ErrorPermission
	ec.keywordPatterns["eperm"] = ErrorPermission
	ec.keywordPatterns["access denied"] = ErrorPermission
	ec.keywordPatterns["forbidden"] = ErrorPermission

	ec.keywordPatterns["syntax error"] = ErrorSyntax
	ec.keywordPatterns["parse error"] = ErrorSyntax
	ec.keywordPatterns["unexpected token"] = ErrorSyntax
	ec.keywordPatterns["invalid syntax"] = ErrorSyntax

	ec.keywordPatterns["connection refused"] = ErrorNetwork
	ec.keywordPatterns["timeout"] = ErrorTimeout
	ec.keywordPatterns["connection timed out"] = ErrorTimeout
	ec.keywordPatterns["network unreachable"] = ErrorNetwork
	ec.keywordPatterns["could not resolve host"] = ErrorNetwork
	ec.keywordPatterns["econnrefused"] = ErrorNetwork

	ec.keywordPatterns["module not found"] = ErrorDependency
	ec.keywordPatterns["package not found"] = ErrorDependency
	ec.keywordPatterns["cannot find module"] = ErrorDependency
	ec.keywordPatterns["no such file or directory"] = ErrorNotFound
	ec.keywordPatterns["command not found"] = ErrorNotFound

	ec.keywordPatterns["merge conflict"] = ErrorMergeConflict
	ec.keywordPatterns["conflict (content)"] = ErrorMergeConflict
	ec.keywordPatterns["both modified"] = ErrorMergeConflict

	ec.keywordPatterns["authentication failed"] = ErrorAuthentication
	ec.keywordPatterns["401"] = ErrorAuthentication
	ec.keywordPatterns["403"] = ErrorPermission
	ec.keywordPatterns["unauthorized"] = ErrorAuthentication

	ec.keywordPatterns["no space left"] = ErrorDiskSpace
	ec.keywordPatterns["disk full"] = ErrorDiskSpace
	ec.keywordPatterns["enospc"] = ErrorDiskSpace

	ec.keywordPatterns["out of memory"] = ErrorMemory
	ec.keywordPatterns["enomem"] = ErrorMemory
	ec.keywordPatterns["killed"] = ErrorMemory

	ec.keywordPatterns["segmentation fault"] = ErrorSegfault
	ec.keywordPatterns["sigsegv"] = ErrorSegfault
	ec.keywordPatterns["core dumped"] = ErrorSegfault

	ec.keywordPatterns["race detected"] = ErrorRaceCondition
	ec.keywordPatterns["data race"] = ErrorRaceCondition

	ec.keywordPatterns["deprecated"] = ErrorDeprecated
	ec.keywordPatterns["no longer supported"] = ErrorDeprecated

	ec.keywordPatterns["test failed"] = ErrorTestFailure
	ec.keywordPatterns["assertion failed"] = ErrorTestFailure
	ec.keywordPatterns["expected"] = ErrorTestFailure

	ec.keywordPatterns["type error"] = ErrorTypeMismatch
	ec.keywordPatterns["type mismatch"] = ErrorTypeMismatch
	ec.keywordPatterns["cannot convert"] = ErrorTypeMismatch

	ec.keywordPatterns["null pointer"] = ErrorNullPointer
	ec.keywordPatterns["nil pointer"] = ErrorNullPointer
	ec.keywordPatterns["nullptr"] = ErrorNullPointer
}

// ClassifyError analyzes the command and exit code to determine error categories
func (ec *ErrorClassifier) ClassifyError(command string, exitCode int, errorOutput string) []ErrorCategory {
	categories := make(map[ErrorCategory]bool)

	// Check exit code patterns
	if exitCategories, exists := ec.exitCodePatterns[exitCode]; exists {
		for _, cat := range exitCategories {
			categories[cat] = true
		}
	}

	// Check command patterns
	commandLower := strings.ToLower(command)
	for pattern, category := range ec.commandPatterns {
		if pattern.MatchString(commandLower) {
			categories[category] = true
		}
	}

	// Check keyword patterns in command and error output
	combinedText := strings.ToLower(command + " " + errorOutput)
	for keyword, category := range ec.keywordPatterns {
		if strings.Contains(combinedText, keyword) {
			categories[category] = true
		}
	}

	// Convert map to slice
	result := make([]ErrorCategory, 0, len(categories))
	for cat := range categories {
		result = append(result, cat)
	}

	// If no specific category found, return generic
	if len(result) == 0 {
		result = append(result, ErrorGeneric)
	}

	return result
}

// GetPrimaryError returns the most specific error category
func (ec *ErrorClassifier) GetPrimaryError(command string, exitCode int, errorOutput string) ErrorCategory {
	categories := ec.ClassifyError(command, exitCode, errorOutput)

	// Priority order: specific errors first, generic last
	priority := []ErrorCategory{
		ErrorSegfault,
		ErrorMergeConflict,
		ErrorRaceCondition,
		ErrorPermission,
		ErrorAuthentication,
		ErrorDiskSpace,
		ErrorMemory,
		ErrorTimeout,
		ErrorNetwork,
		ErrorDependency,
		ErrorTestFailure,
		ErrorBuildFailure,
		ErrorLinting,
		ErrorSyntax,
		ErrorTypeMismatch,
		ErrorNullPointer,
		ErrorConfiguration,
		ErrorDeprecated,
		ErrorNotFound,
		ErrorInfiniteLoop,
		ErrorGeneric,
	}

	for _, prioCategory := range priority {
		for _, foundCategory := range categories {
			if foundCategory == prioCategory {
				return prioCategory
			}
		}
	}

	return ErrorGeneric
}
