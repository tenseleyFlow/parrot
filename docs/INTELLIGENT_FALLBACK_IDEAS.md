# Intelligent Fallback Generator - Design Ideas

## Phase 2: Making Fallbacks Smarter

This document outlines ideas for enhancing the fallback system to be more context-aware and intelligent, moving beyond simple pre-generated insults to command-specific mockery.

---

## Current State (Phase 1)

✅ **Completed:**
- Expanded fallback database with 600+ brutal insults
- Categorized by command type (git, nodejs, docker, python, rust, etc.)
- Pseudo-random selection for variety
- Better command type detection

**Limitations:**
- No awareness of actual command content
- Can't parse error patterns
- Doesn't understand command arguments
- Same insult for different failures of same command type

---

## Phase 2 Goals: Context-Aware Fallbacks

### 1. **Command Pattern Analysis**

Parse the actual command to understand what the user was trying to do.

#### Git Examples:
```go
// Detect specific git operations
"git push" → "Remote rejected you (literally and figuratively)"
"git pull" → "Merge conflicts incoming! Your code vs. reality"
"git commit" → "Even git doesn't want to commit to your code"
"git rebase" → "Rewriting history won't fix your present"
"git merge" → "Can't merge competence into incompetence"
"git clone" → "Cloning failures: Your specialty"
"git checkout" → "Checking out from reality again?"
"git branch" → "Creating branches in your career path to nowhere"
"git reset --hard" → "Wish you could reset your career this easily"
"git stash" → "Stashing problems for later. Classic you."
```

#### Docker Examples:
```bash
"docker build" → "Build failed: Can't dockerize disaster"
"docker run" → "Container refuses to run from you"
"docker-compose up" → "Composing a symphony of failure"
"docker push" → "Registry rejected your image (and your code)"
"docker exec" → "Can't exec into competence"
```

### 2. **Exit Code Intelligence**

Different insults based on common exit codes:

```go
exitCodeInsults := map[int][]string{
    1:   {"Generic failure. Generic developer."},
    2:   {"Misuse of shell command. Misuse of developer title."},
    126: {"Command not executable. Neither are your plans."},
    127: {"Command not found. Neither is your competence."},
    128: {"Invalid exit argument. Invalid career argument."},
    130: {"Ctrl+C? Can't escape your mistakes that easily."},
    137: {"SIGKILL: Someone had to stop you."},
    139: {"Segfault: Your logic segfaulted first."},
    255: {"Exit code overflow: Like your error overflow."},
}
```

### 3. **Error Pattern Recognition**

Parse common error patterns from stderr to provide specific mockery:

#### Permission Errors:
```
"permission denied" → "Even sudo can't give you competence"
"access forbidden" → "Denied access to success"
"insufficient privileges" → "Insufficient everything"
```

#### Network Errors:
```
"connection refused" → "Server ghosted you"
"timeout" → "Timed out waiting for your skill to load"
"host unreachable" → "Your goals are unreachable too"
"certificate error" → "Your certifications are invalid too"
```

#### File/Path Errors:
```
"no such file" → "No such competence either"
"directory not empty" → "Of failures"
"file exists" → "Your file of mistakes exists"
"disk full" → "Full of your mistakes"
```

### 4. **Command Argument Awareness**

Understand what arguments mean for better context:

```go
// Git branch operations
"git push -f" → "Force push? Forcing failure down everyone's throat"
"git commit -m" → "Commit message won't save you from commit mistakes"
"git reset --hard HEAD~1" → "Undo won't undo your career choices"

// Docker operations
"docker run -d" → "Detached mode: Like your detachment from reality"
"docker build --no-cache" → "No cache can save you"
"docker-compose down" → "Down like your employment prospects"

// Package managers
"npm install --production" → "Production? You're not ready for development"
"pip install --upgrade" → "Can't upgrade incompetence"
```

### 5. **Command History Analysis**

Track repeated failures for escalating mockery:

```go
type FailureTracker struct {
    commandCounts map[string]int
    recentFailures []string
}

// First failure
"npm install" → "NPM install failed. Try again?"

// Third failure (same command)
"npm install" → "Third time's the charm? Not for you apparently."

// Tenth failure
"npm install" → "Einstein defined insanity as... oh wait, you're beyond that."
```

### 6. **Time-of-Day Awareness**

Contextual mockery based on when failures happen:

```go
func getTimeBasedInsult() string {
    hour := time.Now().Hour()

    switch {
    case hour >= 0 && hour < 6:
        return "Failing at 3 AM? Sleep won't fix your code."
    case hour >= 9 && hour < 17:
        return "Failing during work hours? At least you're consistent."
    case hour >= 17 && hour < 23:
        return "Still here? Commitment to failure is impressive."
    case hour == 23:
        return "11 PM failure? Tomorrow won't be better."
    }
}
```

### 7. **Contextual Command Chaining**

Detect related command sequences:

```bash
# User runs:
$ npm install
# fails
$ rm -rf node_modules
$ npm install
# fails again

Parrot: "Deleting node_modules won't delete your incompetence."
```

### 8. **Project Context Detection**

Read project files for better context:

```go
// If package.json exists
"npm install" → "634 dependencies. 635 problems."

// If Dockerfile exists
"docker build" → "FROM disaster AS production"

// If .git exists with many branches
"git push" → "42 branches. 0 working."

// If requirements.txt has many packages
"pip install" → "Your dependency list is longer than your employment history."
```

### 9. **Smart Template System**

Template-based generation with variable substitution:

```go
templates := []string{
    "{{command}} failed? {{reason}} not found.",
    "Error in {{command}}: {{error_type}} detected at {{location}}.",
    "{{command}} rejected by {{service}}: {{witty_reason}}.",
}

variables := map[string][]string{
    "reason": {"Competence", "Logic", "Skill", "Talent", "Ability"},
    "error_type": {"Stupidity", "Incompetence", "Chaos", "Disaster"},
    "location": {"keyboard-chair interface", "brain", "neural network"},
    "service": {"reality", "common sense", "basic standards"},
    "witty_reason": {"standards exist", "quality matters", "not today"},
}
```

### 10. **Machine Learning Approach** (Future)

Build a small ML model trained on:
- Common error messages
- Command patterns
- Successful insult patterns
- User response (if they re-run quickly, insult was good!)

Could use a tiny transformer model or Markov chain for:
- Context-aware insult generation
- Pattern recognition in error messages
- Learning what insults work best for which failures

---

## Implementation Priorities

### **High Priority** (Phase 2.1)
1. ✅ Command argument parsing
2. ✅ Exit code specific insults
3. ✅ Basic error pattern recognition
4. ✅ Enhanced command type detection

### **Medium Priority** (Phase 2.2)
5. Command history tracking
6. Time-of-day context
7. Project context detection
8. Template-based generation

### **Low Priority** (Phase 2.3)
9. Command chaining detection
10. ML-based generation

---

## Technical Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  Smart Fallback System                   │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│              Context Analysis Layer                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │  Command     │  │  Error       │  │  Environment │  │
│  │  Parser      │  │  Pattern     │  │  Detector    │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│           Insult Generation Strategy                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │  Pattern     │  │  Template    │  │  Database    │  │
│  │  Match       │  │  Fill        │  │  Lookup      │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│              Response Ranking & Selection                │
└─────────────────────────────────────────────────────────┘
```

---

## Example: Intelligent Fallback Flow

```
User runs: git push origin main
Exit code: 1
Stderr: "rejected (fetch first)"

Step 1: Parse command
- Command: "git"
- Subcommand: "push"
- Args: ["origin", "main"]
- Type: "git-push"

Step 2: Analyze error
- Pattern: "rejected"
- Reason: "fetch first"
- Category: "out-of-sync"

Step 3: Generate context-aware insult
Options:
a) "Git rejected your push. Did you forget to pull? Amateur."
b) "Remote is ahead. So is everyone else in your field."
c) "Fetch first? You should've fetched competence first."

Step 4: Select based on:
- Recent history (if repeated: escalate)
- Time of day (if late: add time mockery)
- Personality setting (savage mode: option C)

Final output: "Fetch first? You should've fetched competence first."
```

---

## Data Structures

```go
type ContextualFallback struct {
    // Command context
    Command     string
    Subcommand  string
    Arguments   []string
    ExitCode    int

    // Error context
    StderrSnippet string
    ErrorPattern  string
    ErrorCategory string

    // Environment context
    ProjectType   string  // "nodejs", "python", "rust", etc.
    TimeOfDay     int
    IsWeekend     bool

    // History context
    FailureCount  int
    LastFailure   time.Time
    SameCommand   bool
}

type InsultGenerator interface {
    Generate(ctx ContextualFallback) string
    Rank(insults []string, ctx ContextualFallback) []string
}
```

---

## Testing Strategy

Create test cases for intelligent fallbacks:

```go
func TestIntelligentFallback(t *testing.T) {
    tests := []struct{
        command  string
        exitCode int
        stderr   string
        expected string  // should contain
    }{
        {
            command:  "git push origin main",
            exitCode: 1,
            stderr:   "rejected",
            expected: "rejected", // insult should mention rejection
        },
        {
            command:  "docker build .",
            exitCode: 1,
            stderr:   "no such file",
            expected: "file", // should mention file issue
        },
    }
}
```

---

## Future Enhancements

1. **User Preferences**: Allow users to configure insult intensity
2. **Learning Mode**: Adapt to user's most common failures
3. **Streaks**: Track failure/success streaks
4. **Achievements**: "100 git push failures in a row!"
5. **Statistics**: `parrot stats` showing failure patterns
6. **Insult API**: External service for fresh insults
7. **Community Insults**: User-submitted insult database
8. **Language Support**: Multilingual mockery

---

## Performance Considerations

- Context analysis should be < 10ms
- Database lookup: O(1) hash-based
- Template generation: Pre-compiled templates
- Error pattern matching: Regex compilation cached
- History tracking: Ring buffer (last 100 failures)

---

## Backwards Compatibility

- Keep expanded database as fallback-of-fallbacks
- Smart fallback is opt-in via config flag
- Graceful degradation if parsing fails
- All new features behind feature flags

```toml
[features]
smart_fallback = true
context_aware_insults = true
command_history_tracking = true
project_detection = true
```

---

## Metrics to Track

- Fallback usage rate
- Context detection success rate
- Most common failure patterns
- Average insult length
- User engagement (re-run speed)

---

## Conclusion

Phase 2 will transform fallbacks from simple pre-generated insults into an intelligent, context-aware mockery system that understands what users are trying to do and why they're failing, delivering brutally specific feedback that's both entertaining and (accidentally) educational.

The foundation is in place with 600+ insults. Now we build intelligence on top.
