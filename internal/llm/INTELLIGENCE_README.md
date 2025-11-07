# Intelligent Insult System - ML-Inspired Architecture

## Overview

The Parrot CLI now features a sophisticated, **ML-inspired intelligent insult selection system** that goes beyond random selection to deliver contextually relevant, personalized feedback based on:

- **Error pattern classification**
- **Command intent analysis**
- **Multi-factor relevance scoring**
- **Adaptive learning from history**
- **Semantic tagging**

This creates a truly smart system that learns and adapts to deliver the most appropriate insult for each failure scenario.

---

## Architecture

### Five-Tier Intelligence System

```
TIER 5: ML-Inspired Semantic Matching (NEW!)
  ↓ (if score < threshold)
TIER 4: Historical Learning & Dynamic Generation
  ↓ (if no match)
TIER 3: LLM-like Context Awareness
  ↓ (if no match)
TIER 2: Environment Detection
  ↓ (if no match)
TIER 1: Basic Context Matching
  ↓ (if no match)
FALLBACK: Static Database
```

---

## Key Components

### 1. Error Classifier (`error_classifier.go`)

Analyzes commands and exit codes to determine specific error types:

- **20+ Error Categories**: Permission errors, syntax errors, network failures, merge conflicts, test failures, segfaults, race conditions, and more
- **Multi-source Analysis**: Combines exit codes, command patterns, and error output
- **Priority Ranking**: Returns the most specific error type first

**Example:**
```go
classifier.ClassifyError("git push origin main", 1, "permission denied")
// Returns: [ErrorPermission, ErrorAuthentication]
```

### 2. Semantic Tagging System (`semantic_tags.go`)

Each insult is tagged with rich metadata:

**Tag Categories:**
- **Error Types**: `permission`, `syntax`, `network`, `merge_conflict`, `test_failure`, etc.
- **Command Types**: `git`, `docker`, `kubernetes`, `node`, `rust`, `python`, etc.
- **User Intent**: `push`, `build`, `test`, `deploy`, `install`, etc.
- **Context**: `late_night`, `ci`, `main_branch`, `streak`, `repeated`, etc.
- **Severity**: `mild`, `sarcastic`, `savage` (1-10 scale)

**Example Tagged Insult:**
```go
{
    Text: "Tests failed. Shocking absolutely no one who read your code",
    Tags: [TagTestFailure, TagTest, TagSavage],
    Severity: 7,
    Weight: 0.9,
}
```

### 3. Intent Parser (`intent_parser.go`)

Extracts semantic meaning from commands:

- **Primary Intent Detection**: What was the user trying to do? (push, build, test, deploy, etc.)
- **Target Extraction**: Files, branches, containers being acted upon
- **Complexity Assessment**: Simple, moderate, or complex command
- **Risk Analysis**: Low, medium, or high risk operation

**Example:**
```go
parser.ParseIntent("git push --force origin main")
// Returns: {
//     PrimaryIntent: "push",
//     Targets: ["origin", "main"],
//     Complexity: "simple",
//     RiskLevel: "high"
// }
```

### 4. Multi-Factor Scoring Algorithm (`insult_scorer.go`)

Ranks insults using weighted factors:

**Scoring Factors:**
- **Tag Matching (35%)**: How well do insult tags match the context?
- **Error Matching (30%)**: How specific is the insult to the error type?
- **Context Relevance (20%)**: Time of day, CI environment, branch type, etc.
- **Novelty (10%)**: Avoid recently shown insults
- **Personality Fit (5%)**: Match user's personality preference (mild/sarcastic/savage)

**Example Scoring:**
```
Permission error + late night + git push + repeated failure:
  → "chmod 777 isn't the answer this time, though I admire your optimism"
  → Score: 0.87 (high relevance)
```

### 5. Insult History Tracker (`insult_history.go`)

Prevents repetition and tracks usage:

- **Persistent Storage**: Saves to `~/.parrot/insult_history.json`
- **Recency Tracking**: Penalizes recently shown insults
- **Frequency Analysis**: Less used insults score higher
- **Automatic Cleanup**: Old entries removed after 7 days
- **Statistics**: Track most frequent insults, total shown, etc.

---

## How It Works

### Selection Flow

1. **Parse Context**: Extract rich context from command, exit code, environment
2. **Classify Error**: Determine specific error types (permission, syntax, network, etc.)
3. **Parse Intent**: Understand what the user was trying to accomplish
4. **Generate Tags**: Combine error tags, context tags, and intent tags
5. **Score All Insults**: Rank every insult in the database using multi-factor scoring
6. **Apply Novelty**: Penalize recently shown insults
7. **Filter by Personality**: Ensure severity matches user preference
8. **Select Best**: Return highest-scoring insult above quality threshold
9. **Record History**: Save for future novelty calculations

### Scoring Example

**Command:** `git push origin main` (Exit Code: 1, Time: 2 AM, Repeated failure)

**Context Tags Generated:**
- `TagGit`
- `TagPush`
- `TagMainBranch`
- `TagLateNight`
- `TagRepeated`

**Top Insult Candidates:**

| Insult | Tag Match | Error Match | Context | Novelty | Personality | **Total** |
|--------|-----------|-------------|---------|---------|-------------|-----------|
| "Push rejected: The remote has standards" | 0.80 | 0.90 | 0.75 | 1.0 | 0.85 | **0.83** ✓ |
| "Failed to push. The remote branch has standards" | 0.75 | 0.85 | 0.70 | 0.90 | 0.80 | **0.78** |
| "Working at 2 AM? Even your rubber duck has clocked out" | 0.60 | 0.50 | 0.95 | 1.0 | 0.75 | **0.69** |

**Winner:** "Push rejected: The remote has standards" (Score: 0.83)

---

## Database

### Insult Statistics

- **200+ Semantically Tagged Insults** in the new system
- **~3000 Static Insults** in the fallback database
- **Comprehensive Coverage**: All major error types, languages, and tools
- **Context-Specific**: Insults for specific scenarios (CI failures, late night coding, merge conflicts, etc.)

### Example Categories

**Permission Errors:**
- "chmod 777 isn't the answer this time, though I admire your optimism"
- "Permission denied. The computer has decided you're not ready for this level of responsibility"

**Test Failures:**
- "Tests failed. Shocking absolutely no one who read your code"
- "Expected: working code. Received: whatever this is"

**Build Failures:**
- "The compiler is personally offended by what you've written"
- "Your build broke so hard it took down the CI server's will to live"

**Late Night Context:**
- "It's 2 AM. The bugs aren't the only thing that needs fixing"
- "Caffeine and bad decisions: a programmer's autobiography"

**CI/Production:**
- "Failed in CI. Congratulations, everyone on the team just got your shame notification"
- "Red pipeline. That's the color of your team's disappointment"

---

## Benefits

### 🎯 **Context Awareness**
No more generic insults! Each failure gets a response tailored to:
- The specific error type
- What you were trying to do
- Your current environment
- Time of day
- Failure patterns

### 🧠 **Learning & Adaptation**
- Tracks which insults you've seen recently
- Avoids repetition
- Learns your common failure patterns
- Adjusts based on your personality preference

### 📊 **Intelligent Ranking**
- Multi-factor scoring ensures relevance
- Quality threshold prevents low-relevance matches
- Falls back gracefully to other tiers if needed

### 🎭 **Personality Matching**
Respects your configured personality:
- **Mild**: Gentle, constructive feedback (severity ≤ 4)
- **Sarcastic**: Witty, clever mocking (severity 4-7)
- **Savage**: Brutal, devastating roasts (severity ≥ 6)

---

## Configuration

### Personality Setting

The intelligent system respects your personality preference:

```toml
# ~/.config/parrot/config.toml
[general]
personality = "sarcastic"  # Options: mild, sarcastic, savage
```

### Fallback Behavior

If the ML-inspired system doesn't find a high-confidence match (score < 0.3), it falls through to Tier 4, 3, 2, 1, and finally the static database. This ensures you **always** get a response.

---

## Technical Details

### Performance

- **Fast**: Scoring algorithm is O(n) where n = number of insults (~200)
- **Lightweight**: Minimal memory footprint
- **Cached**: Database loaded once at startup
- **Persistent**: History stored in JSON for cross-session learning

### Files Created

- `~/.parrot/insult_history.json`: Tracks shown insults and usage statistics
- `~/.parrot/failures.json`: User failure history (existing Tier 4 system)

### Quality Threshold

Insults must score **≥ 0.3** (30%) to be used. This ensures only relevant, contextual insults are shown. Lower-scoring matches fall through to other intelligence tiers.

---

## Examples

### Example 1: Permission Error at 2 AM

**Command:** `sudo rm -rf /important/file` (Exit: 126)

**Analysis:**
- Error Type: Permission denied
- Time: 2 AM (late night)
- Risk: High (destructive command)
- Intent: Delete

**Selected Insult:**
> "Permission denied. The computer has decided you're not ready for this level of responsibility"

**Score Breakdown:**
- Tag Match: 0.90 (permission + late_night tags match)
- Error Match: 1.0 (perfect permission error match)
- Context: 0.85 (late night + high risk)
- Novelty: 1.0 (not shown recently)
- Personality: 0.80 (sarcastic, severity 5)
- **Total: 0.88** ✓

---

### Example 2: Test Failure in CI

**Command:** `npm test` (Exit: 1, CI environment)

**Analysis:**
- Error Type: Test failure
- Environment: GitHub Actions CI
- Project: Node.js
- Intent: Test

**Selected Insult:**
> "Did you test this before committing? Oh wait, that's what the CI is for, right?"

**Score Breakdown:**
- Tag Match: 0.85 (test + node + ci tags)
- Error Match: 0.95 (test failure match)
- Context: 0.90 (CI context bonus)
- Novelty: 1.0
- Personality: 0.85
- **Total: 0.89** ✓

---

### Example 3: Merge Conflict on Main Branch

**Command:** `git merge feature/new-ui` (Exit: 1, on main branch)

**Analysis:**
- Error Type: Merge conflict
- Branch: main (high risk)
- Command: git
- Intent: Merge

**Selected Insult:**
> "<<<<<<< HEAD is not a valid merge resolution strategy"

**Score Breakdown:**
- Tag Match: 0.95 (merge_conflict + git + main_branch)
- Error Match: 1.0 (perfect merge conflict match)
- Context: 0.80 (main branch penalty)
- Novelty: 1.0
- Personality: 0.90
- **Total: 0.92** ✓

---

## Future Enhancements

Possible improvements for future versions:

1. **True ML Model**: Train a lightweight model on historical data
2. **User Preference Learning**: Learn which insults the user appreciates most
3. **Team Patterns**: Share anonymized patterns across team members
4. **Custom Insults**: Allow users to add their own tagged insults
5. **Sentiment Analysis**: Detect frustration level from command patterns
6. **Integration**: Pull error messages from stderr for better classification

---

## Code Structure

```
internal/llm/
├── error_classifier.go      # Classifies error types from exit codes & patterns
├── semantic_tags.go         # Tagged insult database with metadata
├── intent_parser.go         # Extracts command intent and complexity
├── insult_scorer.go         # Multi-factor scoring algorithm
├── insult_history.go        # Persistent history tracking
└── smart_fallback.go        # Integration layer (Tier 5 addition)
```

---

## Summary

The new ML-inspired intelligent insult system transforms Parrot from a simple random insult generator into a **context-aware, adaptive, learning system** that delivers highly relevant feedback tailored to:

✅ **What went wrong** (error classification)
✅ **What you were trying to do** (intent parsing)
✅ **Your environment** (CI, time, branch, project type)
✅ **Your history** (avoiding repetition, learning patterns)
✅ **Your preference** (personality matching)

This represents a **significant intelligence upgrade** that makes every failure... well, at least entertainingly personalized! 🎯
