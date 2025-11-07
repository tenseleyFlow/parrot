package llm

// InsultTag represents a semantic tag for categorizing insults
type InsultTag string

const (
	// Error type tags
	TagPermission     InsultTag = "permission"
	TagSyntax         InsultTag = "syntax"
	TagNetwork        InsultTag = "network"
	TagDependency     InsultTag = "dependency"
	TagMergeConflict  InsultTag = "merge_conflict"
	TagTestFailure    InsultTag = "test_failure"
	TagBuildFailure   InsultTag = "build_failure"
	TagTimeout        InsultTag = "timeout"
	TagAuthentication InsultTag = "authentication"
	TagDiskSpace      InsultTag = "disk_space"
	TagMemory         InsultTag = "memory"
	TagSegfault       InsultTag = "segfault"
	TagLinting        InsultTag = "linting"
	TagTyping         InsultTag = "typing"
	TagDeprecated     InsultTag = "deprecated"

	// Command type tags
	TagGit        InsultTag = "git"
	TagDocker     InsultTag = "docker"
	TagKubernetes InsultTag = "kubernetes"
	TagNode       InsultTag = "node"
	TagPython     InsultTag = "python"
	TagRust       InsultTag = "rust"
	TagGolang     InsultTag = "golang"
	TagJava       InsultTag = "java"

	// Intent tags (what user was trying to do)
	TagPush       InsultTag = "push"
	TagPull       InsultTag = "pull"
	TagCommit     InsultTag = "commit"
	TagBuild      InsultTag = "build"
	TagTest       InsultTag = "test"
	TagDeploy     InsultTag = "deploy"
	TagInstall    InsultTag = "install"
	TagConfigure  InsultTag = "configure"
	TagDebug      InsultTag = "debug"
	TagRefactor   InsultTag = "refactor"
	TagRevert     InsultTag = "revert"

	// Context tags
	TagLateNight  InsultTag = "late_night"
	TagEarlyMorning InsultTag = "early_morning"
	TagWeekend    InsultTag = "weekend"
	TagCI         InsultTag = "ci"
	TagProduction InsultTag = "production"
	TagMainBranch InsultTag = "main_branch"
	TagStreak     InsultTag = "streak"
	TagBeginner   InsultTag = "beginner"
	TagExpert     InsultTag = "expert"

	// Severity tags
	TagMild     InsultTag = "mild"
	TagSarcastic InsultTag = "sarcastic"
	TagSavage   InsultTag = "savage"

	// Situational tags
	TagDeadline  InsultTag = "deadline"
	TagDemo      InsultTag = "demo"
	TagRepeated  InsultTag = "repeated"
	TagSimple    InsultTag = "simple"
	TagComplex   InsultTag = "complex"
	TagCopyPaste InsultTag = "copy_paste"
)

// TaggedInsult represents an insult with semantic metadata
type TaggedInsult struct {
	Text     string
	Tags     []InsultTag
	Severity int // 1-10, for personality filtering
	Weight   float64 // Base weight for scoring, 0.0-1.0
}

// InsultDatabase stores tagged insults organized by categories
type InsultDatabase struct {
	Insults []TaggedInsult
}

// NewInsultDatabase creates a database with all tagged insults
func NewInsultDatabase() *InsultDatabase {
	db := &InsultDatabase{
		Insults: make([]TaggedInsult, 0, 500),
	}
	db.loadTaggedInsults()
	return db
}

func (db *InsultDatabase) loadTaggedInsults() {
	// Permission errors
	db.addInsults([]TaggedInsult{
		{Text: "chmod 777 isn't the answer this time, though I admire your optimism", Tags: []InsultTag{TagPermission, TagSarcastic, TagSimple}, Severity: 4, Weight: 0.9},
		{Text: "Ah yes, permission denied. Have you tried asking nicely?", Tags: []InsultTag{TagPermission, TagSarcastic}, Severity: 3, Weight: 0.8},
		{Text: "sudo make me a sandwich... but you can't even sudo properly", Tags: []InsultTag{TagPermission, TagSavage, TagSimple}, Severity: 7, Weight: 0.85},
		{Text: "Permission denied. The computer has decided you're not ready for this level of responsibility", Tags: []InsultTag{TagPermission, TagMild}, Severity: 2, Weight: 0.75},
		{Text: "Access denied. Even your own computer doesn't trust you anymore", Tags: []InsultTag{TagPermission, TagSavage, TagStreak}, Severity: 8, Weight: 0.8},
		{Text: "The system has revoked your privileges. Can't say I blame it", Tags: []InsultTag{TagPermission, TagSarcastic, TagRepeated}, Severity: 6, Weight: 0.85},
	})

	// Git merge conflicts
	db.addInsults([]TaggedInsult{
		{Text: "Merge conflict? More like 'I forgot to pull again' conflict", Tags: []InsultTag{TagMergeConflict, TagGit, TagSarcastic, TagRepeated}, Severity: 5, Weight: 0.95},
		{Text: "Git is trying to tell you something. Maybe listen for once?", Tags: []InsultTag{TagMergeConflict, TagGit, TagSarcastic}, Severity: 4, Weight: 0.8},
		{Text: "Congratulations on creating a merge conflict that will haunt your team for days", Tags: []InsultTag{TagMergeConflict, TagGit, TagSavage, TagMainBranch}, Severity: 7, Weight: 0.9},
		{Text: "<<<<<<< HEAD is not a valid merge resolution strategy", Tags: []InsultTag{TagMergeConflict, TagGit, TagSarcastic, TagBeginner}, Severity: 6, Weight: 0.85},
		{Text: "This merge conflict has more drama than a soap opera", Tags: []InsultTag{TagMergeConflict, TagGit, TagMild}, Severity: 3, Weight: 0.7},
		{Text: "Your branch and origin have diverged. So has your common sense", Tags: []InsultTag{TagMergeConflict, TagGit, TagSavage}, Severity: 7, Weight: 0.8},
		{Text: "Maybe try communicating with your team? Novel concept, I know", Tags: []InsultTag{TagMergeConflict, TagGit, TagSarcastic, TagRepeated}, Severity: 5, Weight: 0.85},
	})

	// Test failures
	db.addInsults([]TaggedInsult{
		{Text: "Tests failed. Shocking absolutely no one who read your code", Tags: []InsultTag{TagTestFailure, TagTest, TagSavage}, Severity: 7, Weight: 0.9},
		{Text: "Red is a lovely color, but not in your test suite", Tags: []InsultTag{TagTestFailure, TagTest, TagSarcastic}, Severity: 4, Weight: 0.8},
		{Text: "Your code failed more tests than a student who didn't study", Tags: []InsultTag{TagTestFailure, TagTest, TagSarcastic}, Severity: 5, Weight: 0.85},
		{Text: "The tests are failing to hide their disappointment", Tags: []InsultTag{TagTestFailure, TagTest, TagMild}, Severity: 3, Weight: 0.75},
		{Text: "Expected: working code. Received: whatever this is", Tags: []InsultTag{TagTestFailure, TagTest, TagSarcastic, TagTyping}, Severity: 6, Weight: 0.9},
		{Text: "Your test coverage is great! Too bad your code isn't", Tags: []InsultTag{TagTestFailure, TagTest, TagSavage}, Severity: 7, Weight: 0.85},
		{Text: "Did you test this before committing? Oh wait, that's what the CI is for, right?", Tags: []InsultTag{TagTestFailure, TagTest, TagCI, TagSarcastic}, Severity: 6, Weight: 0.95},
		{Text: "Failing tests on a Friday afternoon? Bold strategy", Tags: []InsultTag{TagTestFailure, TagTest, TagSarcastic, TagDeadline}, Severity: 5, Weight: 0.85},
	})

	// Build failures
	db.addInsults([]TaggedInsult{
		{Text: "Build failed. Have you considered a career in gardening?", Tags: []InsultTag{TagBuildFailure, TagBuild, TagSavage}, Severity: 7, Weight: 0.85},
		{Text: "Your build failed faster than my faith in this project", Tags: []InsultTag{TagBuildFailure, TagBuild, TagSavage, TagStreak}, Severity: 8, Weight: 0.9},
		{Text: "Compilation error: your code doesn't even build, let alone run", Tags: []InsultTag{TagBuildFailure, TagBuild, TagMild}, Severity: 3, Weight: 0.8},
		{Text: "The compiler is personally offended by what you've written", Tags: []InsultTag{TagBuildFailure, TagBuild, TagSarcastic}, Severity: 5, Weight: 0.8},
		{Text: "It doesn't compile. Shocking twist in this debugging saga", Tags: []InsultTag{TagBuildFailure, TagBuild, TagSarcastic, TagRepeated}, Severity: 6, Weight: 0.85},
		{Text: "Your build broke so hard it took down the CI server's will to live", Tags: []InsultTag{TagBuildFailure, TagBuild, TagCI, TagSavage}, Severity: 8, Weight: 0.95},
		{Text: "'It works on my machine' - Famous last words before this build failure", Tags: []InsultTag{TagBuildFailure, TagBuild, TagSarcastic}, Severity: 5, Weight: 0.85},
	})

	// Syntax errors
	db.addInsults([]TaggedInsult{
		{Text: "Syntax error. Did your cat walk across the keyboard again?", Tags: []InsultTag{TagSyntax, TagSarcastic, TagSimple}, Severity: 4, Weight: 0.8},
		{Text: "Unexpected token. The only thing unexpected is that you thought this would work", Tags: []InsultTag{TagSyntax, TagSarcastic}, Severity: 5, Weight: 0.85},
		{Text: "Parse error: this code is literally unreadable. Figuratively too", Tags: []InsultTag{TagSyntax, TagSavage}, Severity: 7, Weight: 0.8},
		{Text: "Missing semicolon. Or was it a missing brain cell?", Tags: []InsultTag{TagSyntax, TagSavage, TagSimple}, Severity: 6, Weight: 0.9},
		{Text: "Syntax error on line 1. That's impressively bad", Tags: []InsultTag{TagSyntax, TagSarcastic, TagBeginner}, Severity: 5, Weight: 0.85},
		{Text: "Your linter is having an existential crisis trying to parse this", Tags: []InsultTag{TagSyntax, TagLinting, TagSarcastic}, Severity: 6, Weight: 0.85},
	})

	// Network/timeout errors
	db.addInsults([]TaggedInsult{
		{Text: "Connection timeout. Unlike your patience, which ran out hours ago", Tags: []InsultTag{TagTimeout, TagNetwork, TagSarcastic, TagLateNight}, Severity: 5, Weight: 0.9},
		{Text: "Network unreachable. Kind of like your project deadline", Tags: []InsultTag{TagNetwork, TagSarcastic, TagDeadline}, Severity: 6, Weight: 0.9},
		{Text: "Could not resolve host. Have you tried turning it off and leaving it off?", Tags: []InsultTag{TagNetwork, TagSavage}, Severity: 7, Weight: 0.8},
		{Text: "Connection refused. Even the server wants nothing to do with this", Tags: []InsultTag{TagNetwork, TagSarcastic}, Severity: 5, Weight: 0.85},
		{Text: "Request timed out. So did everyone's patience with your debugging", Tags: []InsultTag{TagTimeout, TagNetwork, TagSavage, TagStreak}, Severity: 7, Weight: 0.85},
		{Text: "Network error. Have you considered sending a carrier pigeon instead?", Tags: []InsultTag{TagNetwork, TagMild}, Severity: 3, Weight: 0.75},
	})

	// Dependency errors
	db.addInsults([]TaggedInsult{
		{Text: "Module not found. Much like your understanding of package management", Tags: []InsultTag{TagDependency, TagInstall, TagSarcastic}, Severity: 5, Weight: 0.9},
		{Text: "Did you forget to npm install? That's what the error message literally says", Tags: []InsultTag{TagDependency, TagNode, TagInstall, TagSarcastic, TagSimple}, Severity: 6, Weight: 0.95},
		{Text: "Dependency hell is real, and you're living in it", Tags: []InsultTag{TagDependency, TagInstall, TagSavage}, Severity: 7, Weight: 0.85},
		{Text: "Package not found. Try reading the docs? Revolutionary, I know", Tags: []InsultTag{TagDependency, TagInstall, TagSarcastic}, Severity: 5, Weight: 0.8},
		{Text: "Your package.json has more issues than a therapist's waiting room", Tags: []InsultTag{TagDependency, TagNode, TagSavage}, Severity: 8, Weight: 0.85},
		{Text: "Cargo couldn't find that crate. Maybe check the spelling this time?", Tags: []InsultTag{TagDependency, TagRust, TagInstall, TagSarcastic}, Severity: 4, Weight: 0.9},
	})

	// Authentication errors
	db.addInsults([]TaggedInsult{
		{Text: "Authentication failed. Did you forget your password or your dignity?", Tags: []InsultTag{TagAuthentication, TagSarcastic}, Severity: 5, Weight: 0.85},
		{Text: "401 Unauthorized. Even the API doesn't want to talk to you", Tags: []InsultTag{TagAuthentication, TagSavage}, Severity: 7, Weight: 0.9},
		{Text: "Token expired. Unlike your enthusiasm for reading documentation", Tags: []InsultTag{TagAuthentication, TagSarcastic}, Severity: 5, Weight: 0.8},
		{Text: "Invalid credentials. Have you tried using the right ones?", Tags: []InsultTag{TagAuthentication, TagSarcastic, TagSimple}, Severity: 4, Weight: 0.85},
		{Text: "Access token rejected. The server has standards, apparently", Tags: []InsultTag{TagAuthentication, TagSavage}, Severity: 6, Weight: 0.8},
	})

	// Docker/container errors
	db.addInsults([]TaggedInsult{
		{Text: "Container failed to start. Much like your understanding of Docker", Tags: []InsultTag{TagDocker, TagSarcastic}, Severity: 5, Weight: 0.85},
		{Text: "Docker daemon not running. Because even daemons need a break from your code", Tags: []InsultTag{TagDocker, TagSarcastic}, Severity: 6, Weight: 0.8},
		{Text: "Image not found. Try docker pull? Or pull yourself together?", Tags: []InsultTag{TagDocker, TagSarcastic, TagSimple}, Severity: 5, Weight: 0.9},
		{Text: "Your Dockerfile has more layers than an onion, and makes me cry just as much", Tags: []InsultTag{TagDocker, TagSavage}, Severity: 7, Weight: 0.8},
		{Text: "Port already in use. Because apparently you can't keep track of your own containers", Tags: []InsultTag{TagDocker, TagSarcastic, TagRepeated}, Severity: 5, Weight: 0.85},
	})

	// Late night / time-based
	db.addInsults([]TaggedInsult{
		{Text: "It's 2 AM. The bugs aren't the only thing that needs fixing", Tags: []InsultTag{TagLateNight, TagMild}, Severity: 3, Weight: 0.9},
		{Text: "Late night debugging? Tomorrow-you is going to hate today-you", Tags: []InsultTag{TagLateNight, TagSarcastic}, Severity: 4, Weight: 0.85},
		{Text: "Caffeine and bad decisions: a programmer's autobiography, Chapter ", Tags: []InsultTag{TagLateNight, TagSarcastic}, Severity: 5, Weight: 0.8},
		{Text: "Working at 3 AM? Even your rubber duck has clocked out", Tags: []InsultTag{TagLateNight, TagSavage}, Severity: 6, Weight: 0.9},
	})

	// CI/Production failures
	db.addInsults([]TaggedInsult{
		{Text: "Failed in CI. Congratulations, everyone on the team just got your shame notification", Tags: []InsultTag{TagCI, TagSavage}, Severity: 8, Weight: 0.95},
		{Text: "Breaking the build? That's a paddlin'", Tags: []InsultTag{TagCI, TagBuildFailure, TagSarcastic}, Severity: 6, Weight: 0.9},
		{Text: "Red pipeline. That's the color of your team's disappointment", Tags: []InsultTag{TagCI, TagSavage, TagMainBranch}, Severity: 8, Weight: 0.95},
		{Text: "This failed in production? Bold move, Cotton", Tags: []InsultTag{TagProduction, TagSavage, TagMainBranch}, Severity: 9, Weight: 0.95},
		{Text: "CI failed. Did you test locally? Silly question, I know you didn't", Tags: []InsultTag{TagCI, TagSarcastic, TagTest}, Severity: 7, Weight: 0.9},
	})

	// Streak-based (repeated failures)
	db.addInsults([]TaggedInsult{
		{Text: "Again? At this point it's just impressive", Tags: []InsultTag{TagRepeated, TagStreak, TagSarcastic}, Severity: 6, Weight: 0.95},
		{Text: "Failing the same command repeatedly is the definition of insanity", Tags: []InsultTag{TagRepeated, TagStreak, TagSavage}, Severity: 8, Weight: 0.95},
		{Text: "This is the third time. Three strikes and you're... still trying, apparently", Tags: []InsultTag{TagRepeated, TagStreak, TagSarcastic}, Severity: 7, Weight: 0.9},
		{Text: "Are you speed-running bad decisions?", Tags: []InsultTag{TagRepeated, TagStreak, TagSavage}, Severity: 8, Weight: 0.9},
		{Text: "Persistent in failure. That's one way to build character", Tags: []InsultTag{TagRepeated, TagStreak, TagMild}, Severity: 4, Weight: 0.85},
	})

	// Type errors
	db.addInsults([]TaggedInsult{
		{Text: "Type error. Have you met TypeScript? Apparently not", Tags: []InsultTag{TagTyping, TagNode, TagSarcastic}, Severity: 5, Weight: 0.85},
		{Text: "Type mismatch. Like your skills and your job requirements", Tags: []InsultTag{TagTyping, TagSavage}, Severity: 7, Weight: 0.8},
		{Text: "The compiler is confused. To be fair, so am I after reading your code", Tags: []InsultTag{TagTyping, TagSarcastic}, Severity: 6, Weight: 0.8},
		{Text: "String expected, received chaos", Tags: []InsultTag{TagTyping, TagSarcastic}, Severity: 5, Weight: 0.85},
	})

	// Python specific
	db.addInsults([]TaggedInsult{
		{Text: "IndentationError. Python: where whitespace is apparently too difficult", Tags: []InsultTag{TagPython, TagSyntax, TagSarcastic}, Severity: 5, Weight: 0.9},
		{Text: "NameError: name 'competence' is not defined", Tags: []InsultTag{TagPython, TagSavage}, Severity: 7, Weight: 0.8},
		{Text: "ImportError. Did you activate your venv? Don't answer, I know you didn't", Tags: []InsultTag{TagPython, TagDependency, TagSarcastic, TagRepeated}, Severity: 6, Weight: 0.95},
	})

	// Rust specific
	db.addInsults([]TaggedInsult{
		{Text: "Borrow checker says no. And honestly, it has a point", Tags: []InsultTag{TagRust, TagSarcastic}, Severity: 5, Weight: 0.85},
		{Text: "Fighting the borrow checker? The borrow checker always wins", Tags: []InsultTag{TagRust, TagMild}, Severity: 3, Weight: 0.8},
		{Text: "rustc rejected your lifetime annotations faster than I'm losing faith", Tags: []InsultTag{TagRust, TagSavage}, Severity: 7, Weight: 0.8},
	})

	// Git push failures
	db.addInsults([]TaggedInsult{
		{Text: "Push rejected. Did you pull first? Novel concept", Tags: []InsultTag{TagGit, TagPush, TagSarcastic, TagRepeated}, Severity: 5, Weight: 0.95},
		{Text: "Failed to push. The remote branch has standards", Tags: []InsultTag{TagGit, TagPush, TagSavage}, Severity: 7, Weight: 0.85},
		{Text: "Push rejected: branch protection. Thank goodness someone protected main from you", Tags: []InsultTag{TagGit, TagPush, TagMainBranch, TagSavage}, Severity: 8, Weight: 0.95},
		{Text: "git push --force? Are we learning nothing from our mistakes?", Tags: []InsultTag{TagGit, TagPush, TagSarcastic, TagMainBranch}, Severity: 7, Weight: 0.9},
	})

	// Kubernetes/k8s specific
	db.addInsults([]TaggedInsult{
		{Text: "Pod failed to start. Like your understanding of Kubernetes", Tags: []InsultTag{TagKubernetes, TagSarcastic}, Severity: 5, Weight: 0.85},
		{Text: "CrashLoopBackOff: when even the cluster gives up on you", Tags: []InsultTag{TagKubernetes, TagSavage}, Severity: 7, Weight: 0.9},
		{Text: "ImagePullBackOff. Did you push the image? Silly question", Tags: []InsultTag{TagKubernetes, TagDocker, TagSarcastic}, Severity: 6, Weight: 0.9},
		{Text: "Your YAML is invalid. Tabs in YAML? Really?", Tags: []InsultTag{TagKubernetes, TagSyntax, TagSarcastic, TagSimple}, Severity: 6, Weight: 0.9},
	})

	// Generic but context-aware
	db.addInsults([]TaggedInsult{
		{Text: "Command failed successfully... wait, no, just failed", Tags: []InsultTag{TagSarcastic}, Severity: 4, Weight: 0.7},
		{Text: "Error 404: Success not found", Tags: []InsultTag{TagSarcastic}, Severity: 5, Weight: 0.75},
		{Text: "Congratulations, you've discovered a new way to break things", Tags: []InsultTag{TagSarcastic, TagComplex}, Severity: 5, Weight: 0.75},
		{Text: "This error message is clearer than your code comments", Tags: []InsultTag{TagSarcastic}, Severity: 5, Weight: 0.7},
		{Text: "If at first you don't succeed, fail, fail again", Tags: []InsultTag{TagStreak, TagMild}, Severity: 3, Weight: 0.8},
	})
}

func (db *InsultDatabase) addInsults(insults []TaggedInsult) {
	db.Insults = append(db.Insults, insults...)
}

// FilterByTags returns insults that match any of the provided tags
func (db *InsultDatabase) FilterByTags(tags []InsultTag) []TaggedInsult {
	var matching []TaggedInsult
	for _, insult := range db.Insults {
		if hasAnyTag(insult.Tags, tags) {
			matching = append(matching, insult)
		}
	}
	return matching
}

// FilterBySeverity returns insults within the severity range
func (db *InsultDatabase) FilterBySeverity(minSeverity, maxSeverity int) []TaggedInsult {
	var matching []TaggedInsult
	for _, insult := range db.Insults {
		if insult.Severity >= minSeverity && insult.Severity <= maxSeverity {
			matching = append(matching, insult)
		}
	}
	return matching
}

func hasAnyTag(insultTags []InsultTag, searchTags []InsultTag) bool {
	for _, searchTag := range searchTags {
		for _, insultTag := range insultTags {
			if insultTag == searchTag {
				return true
			}
		}
	}
	return false
}
