package llm

// Tier 6 Utility Functions - Shared across new intelligent systems
// This prevents function redeclaration issues

// min3 returns minimum of three integers
func min3Tier6(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}

// minTier6 returns minimum of two integers
func minTier6(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// maxTier6 returns maximum of two integers
func maxTier6(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// toLowerTier6 converts string to lowercase
func toLowerTier6(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}

// containsWordTier6 checks if text contains word
func containsWordTier6(text, word string) bool {
	i := 0
	wordLower := toLowerTier6(word)
	textLower := toLowerTier6(text)

	for i+len(wordLower) <= len(textLower) {
		if textLower[i:i+len(wordLower)] == wordLower {
			// Check word boundaries
			beforeOk := i == 0 || !isAlphaNumTier6(rune(textLower[i-1]))
			afterOk := i+len(wordLower) == len(textLower) || !isAlphaNumTier6(rune(textLower[i+len(wordLower)]))

			if beforeOk && afterOk {
				return true
			}
		}
		i++
	}

	return false
}

// isAlphaNumTier6 checks if rune is alphanumeric
func isAlphaNumTier6(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

// findSubstringTier6 finds substring in string
func findSubstringTier6(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// levenshteinDistanceTier6 calculates edit distance between two strings
func levenshteinDistanceTier6(s1, s2 string) int {
	len1, len2 := len(s1), len(s2)

	// Create matrix
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
	}

	// Initialize first row and column
	for i := 0; i <= len1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min3Tier6(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len1][len2]
}

// hasTagTier6 checks if tags contain target
func hasTagTier6(tags []InsultTag, target InsultTag) bool {
	for _, tag := range tags {
		if tag == target {
			return true
		}
	}
	return false
}

// errorCategoriesToTagsTier6 converts error categories to insult tags
func errorCategoriesToTagsTier6(categories []ErrorCategory) []InsultTag {
	tags := make([]InsultTag, 0, len(categories))

	for _, cat := range categories {
		switch cat {
		case ErrorPermission:
			tags = append(tags, TagPermission)
		case ErrorSyntax:
			tags = append(tags, TagSyntax)
		case ErrorNetwork:
			tags = append(tags, TagNetwork)
		case ErrorMergeConflict:
			tags = append(tags, TagMergeConflict)
		case ErrorTestFailure:
			tags = append(tags, TagTest)
		case ErrorBuildFailure:
			tags = append(tags, TagBuild)
		case ErrorTimeout:
			tags = append(tags, TagTimeout)
		}
	}

	return tags
}

// replaceWordTier6 replaces word in text
func replaceWordTier6(text, old, new string) string {
	result := ""
	i := 0

	for i < len(text) {
		if i+len(old) <= len(text) && text[i:i+len(old)] == old {
			// Check word boundaries
			beforeOk := i == 0 || !isAlphaNumTier6(rune(text[i-1]))
			afterOk := i+len(old) == len(text) || !isAlphaNumTier6(rune(text[i+len(old)]))

			if beforeOk && afterOk {
				result += new
				i += len(old)
				continue
			}
		}
		result += string(text[i])
		i++
	}

	return result
}
