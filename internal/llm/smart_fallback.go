package llm

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SmartFallbackContext contains all context needed for intelligent fallback generation
type SmartFallbackContext struct {
	Command      string
	CommandType  string
	Subcommand   string
	Arguments    []string
	ExitCode     int
	FullCommand  string

	// Tier 1 Intelligence
	Environment      map[string]string
	WorkingDir       string
	TimeOfDay        int  // Hour 0-23
	FileExtensions   []string
	CommandLength    int
	NumericArgs      []int
	Shell            string
	HasPipes         bool
	HasChaining      bool

	// Tier 2 Intelligence
	GitBranch      string
	ProjectType    string // "node", "rust", "go", "python", "java", etc.
	ProjectFiles   []string // List of project files found
}

// ParseCommandContext extracts context from a command for intelligent fallback
func ParseCommandContext(command string, commandType string, exitCode string) SmartFallbackContext {
	ctx := SmartFallbackContext{
		FullCommand: command,
		CommandType: commandType,
		Environment: make(map[string]string),
	}

	// Parse exit code
	if ec, err := strconv.Atoi(exitCode); err == nil {
		ctx.ExitCode = ec
	}

	// Parse command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ctx
	}

	ctx.Command = parts[0]
	if len(parts) > 1 {
		ctx.Subcommand = parts[1]
		ctx.Arguments = parts[2:]
	}

	// Tier 1 Intelligence Gathering
	ctx.CommandLength = len(command)
	ctx.TimeOfDay = time.Now().Hour()
	ctx.HasPipes = strings.Contains(command, "|")
	ctx.HasChaining = strings.Contains(command, "&&") || strings.Contains(command, "||")

	// Get working directory
	if wd, err := os.Getwd(); err == nil {
		ctx.WorkingDir = wd
	}

	// Get shell type
	ctx.Shell = filepath.Base(os.Getenv("SHELL"))

	// Collect interesting environment variables
	envVars := []string{"CI", "DEBUG", "NODE_ENV", "PROD", "PRODUCTION", "STAGING", "DEVELOPMENT",
		"GITHUB_ACTIONS", "GITLAB_CI", "JENKINS_HOME", "CIRCLECI", "TRAVIS"}
	for _, envVar := range envVars {
		if val := os.Getenv(envVar); val != "" {
			ctx.Environment[envVar] = val
		}
	}

	// Extract file extensions from command
	extRegex := regexp.MustCompile(`\.\w{1,6}\b`)
	exts := extRegex.FindAllString(command, -1)
	ctx.FileExtensions = exts

	// Extract numeric arguments (ports, chmod values, PIDs, etc.)
	numRegex := regexp.MustCompile(`\b\d+\b`)
	nums := numRegex.FindAllString(command, -1)
	for _, num := range nums {
		if n, err := strconv.Atoi(num); err == nil {
			ctx.NumericArgs = append(ctx.NumericArgs, n)
		}
	}

	// Tier 2 Intelligence Gathering

	// Detect git branch if in a git repo
	ctx.GitBranch = detectGitBranch()

	// Detect project type by checking for common project files
	ctx.ProjectType, ctx.ProjectFiles = detectProjectType()

	return ctx
}

// GenerateSmartFallback generates a context-aware insult
func GenerateSmartFallback(ctx SmartFallbackContext) string {
	// Tier 2 Intelligence - Highest Priority

	// 1. Git branch awareness
	if insult := getGitBranchInsult(ctx); insult != "" {
		return insult
	}

	// 2. Project type detection
	if insult := getProjectTypeInsult(ctx); insult != "" {
		return insult
	}

	// Tier 1 Intelligence - Priority Order

	// 3. Environment-specific insults
	if insult := getEnvironmentInsult(ctx); insult != "" {
		return insult
	}

	// 2. Time-sensitive insults
	if insult := getTimeOfDayInsult(ctx); insult != "" {
		return insult
	}

	// 3. Working directory insults
	if insult := getWorkingDirInsult(ctx); insult != "" {
		return insult
	}

	// 4. File extension insults
	if insult := getFileExtensionInsult(ctx); insult != "" {
		return insult
	}

	// 5. Numeric argument insults
	if insult := getNumericArgumentInsult(ctx); insult != "" {
		return insult
	}

	// 6. Shell-specific insults
	if insult := getShellInsult(ctx); insult != "" {
		return insult
	}

	// 7. Command complexity insults
	if insult := getCommandComplexityInsult(ctx); insult != "" {
		return insult
	}

	// Existing intelligence layers

	// 8. Exit code specific insults
	if insult := getExitCodeInsult(ctx); insult != "" {
		return insult
	}

	// 9. Command-specific patterns
	if insult := getCommandPatternInsult(ctx); insult != "" {
		return insult
	}

	// 10. Argument-aware insults
	if insult := getArgumentAwareInsult(ctx); insult != "" {
		return insult
	}

	// 11. Fall back to expanded database
	return GetExpandedFallback(ctx.CommandType, ctx.FullCommand)
}

// getExitCodeInsult returns insults specific to common exit codes
func getExitCodeInsult(ctx SmartFallbackContext) string {
	exitCodeInsults := map[int][]string{
		1: {
			"Exit code 1: One failure, infinite disappointment.",
			"Generic error for a generic developer.",
			"Exit 1: First step to unemployment.",
			"Error level 1: Your competence level 0.",
			"Failed with distinction: Exit code 1.",
		},
		2: {
			"Exit 2: Misuse of command. Misuse of developer title.",
			"Built-in syntax error: You're a built-in failure.",
			"Wrong arguments. Wrong career.",
			"Command misuse detected. Life misuse detected.",
			"Exit 2: Two brain cells, neither working.",
		},
		126: {
			"Permission denied: Can't execute what you wrote anyway.",
			"Not executable: Neither are your plans.",
			"126: Command found, competence not found.",
			"Can't execute: Your logic isn't executable either.",
		},
		127: {
			"Command not found: Neither is your skill.",
			"127: Path to success not found.",
			"Command doesn't exist. Your competence doesn't exist.",
			"Not in PATH: You're not on the path to success.",
			"Command not found: Story of your career.",
		},
		128: {
			"Invalid exit argument. Invalid life argument.",
			"Exit 128: One-two-eight steps to failure.",
			"Invalid signal: Your brain sends invalid signals.",
		},
		130: {
			"Ctrl+C? Can't escape your mistakes that easily.",
			"Interrupted: Like your thought process.",
			"SIGINT: Significance? Interrupted.",
			"Killed by keyboard: At least something stopped you.",
		},
		137: {
			"SIGKILL: Something had to stop you forcefully.",
			"Exit 137: Murdered by the system. Deservedly.",
			"Killed with prejudice: OOMKiller knows best.",
			"137: You were killed. Your code was mercy-killed.",
		},
		139: {
			"Segmentation fault: Your logic segfaulted first.",
			"139: Memory violation. Logic violation. Everything violation.",
			"Segfault: Your brain segfaulted at compile time.",
			"Core dumped: Your core competencies dumped earlier.",
		},
		143: {
			"SIGTERM: Terminated for being terrible.",
			"Graceful termination: More graceful than your code.",
			"143: Terminated by common sense.",
		},
		255: {
			"Exit 255: Overflow of failure.",
			"Maximum exit code: Maximum incompetence.",
			"255: You maxed out the failure counter.",
			"Exit code overflow: Like your error overflow.",
		},
	}

	if insults, exists := exitCodeInsults[ctx.ExitCode]; exists {
		return selectInsult(insults, ctx.FullCommand)
	}

	return ""
}

// getCommandPatternInsult returns insults based on command + subcommand patterns
func getCommandPatternInsult(ctx SmartFallbackContext) string {
	pattern := ctx.Command + " " + ctx.Subcommand
	pattern = strings.TrimSpace(pattern)

	commandPatterns := map[string][]string{
		// Git operations
		"git push": {
			"Push rejected: The remote has standards.",
			"Git push failed: Even version control rejects you.",
			"Rejected by remote: Story of your life.",
			"Push denied: Your code is unpushable.",
			"Remote said no: Listen to the remote.",
			"Can't push incompetence to production.",
			"Git push failed: Your career trajectory in command form.",
		},
		"git pull": {
			"Pull failed: Can't pull competence from thin air.",
			"Merge conflicts incoming: Your code vs. reality.",
			"Pull rejected: Your branch diverged from sanity.",
			"Can't pull: You're already pulling everyone down.",
			"Fetch failed: Can't fetch what doesn't exist.",
		},
		"git commit": {
			"Commit failed: Even git won't commit to your code.",
			"Nothing to commit: Nothing worth committing.",
			"Pre-commit hook failed: Your code failed harder.",
			"Can't commit disaster: Wait, git tried and failed.",
			"Commit message empty: Like your understanding.",
		},
		"git merge": {
			"Merge failed: Can't merge competence into chaos.",
			"Conflict resolution required: Start with your career.",
			"Merge aborted: Smart choice by git.",
			"Auto-merge failed: Manual merge won't help you either.",
		},
		"git clone": {
			"Clone failed: Repository running away from you.",
			"Can't clone competence: It doesn't exist to clone.",
			"Permission denied: Even public repos protect themselves.",
			"Clone timeout: Repository chose death over your attention.",
		},
		"git rebase": {
			"Rebase failed: Can't rebase on a foundation of failure.",
			"Interactive rebase: Interactively watching you fail.",
			"Rebase conflict: Your existence conflicts with success.",
			"Can't rewrite history to hide your incompetence.",
		},
		"git checkout": {
			"Checkout failed: Can't check out from reality.",
			"Branch not found: Neither is your competence.",
			"Detached HEAD: Matches your detachment from reality.",
			"Already on that branch: Already on the failure branch.",
		},
		"git reset": {
			"Reset failed: Can't reset your mistakes that easily.",
			"Hard reset won't fix soft skills.",
			"Resetting to a previous commit won't fix current you.",
		},
		"git stash": {
			"Stash failed: Can't stash your incompetence away.",
			"Nothing to stash: Nothing worth saving.",
			"Stash apply failed: Your problems can't be applied away.",
		},
		"git branch": {
			"Branch creation failed: Branching into more failure.",
			"Can't create branch: Too many failure branches already.",
			"Branch diverged: You diverged from competence long ago.",
		},
		"git fetch": {
			"Fetch failed: Can't fetch common sense.",
			"Nothing to fetch: Nothing to learn from you either.",
			"Remote unreachable: Like your career goals.",
		},

		// Docker operations
		"docker build": {
			"Build failed: Can't dockerize disaster.",
			"Dockerfile syntax error: Your syntax is always wrong.",
			"Build context too large: Your mistakes are infinite.",
			"Layer failed: All your layers are failures.",
			"FROM scratch: You are scratch.",
			"Build arg undefined: Like your competence.",
		},
		"docker run": {
			"Container exited immediately: Smart container.",
			"Run failed: Nothing wants to run for you.",
			"Port binding failed: Can't bind success to you.",
			"Volume mount error: Can't mount your chaos.",
			"Container crashed on startup: Your code in container form.",
		},
		"docker push": {
			"Push denied: Registry has standards.",
			"Authentication failed: You're not authenticated as competent.",
			"Image push rejected: Your image is not production-ready.",
			"Manifest invalid: Your competence manifest is invalid.",
		},
		"docker pull": {
			"Pull failed: Can't pull what doesn't work.",
			"Image not found: Your skill image doesn't exist.",
			"Digest invalid: Can't digest your code.",
		},
		"docker-compose up": {
			"Compose failed: Can't compose order from chaos.",
			"Service unhealthy: You're the unhealthy service.",
			"Network creation failed: Your networking is broken too.",
			"Volume error: Can't volume-ize your mistakes.",
		},
		"docker exec": {
			"Exec failed: Can't exec into disaster.",
			"Container not running: Your competence isn't running either.",
			"No such container: No such developer.",
		},

		// NPM operations
		"npm install": {
			"Install failed: NPM refuses to install for you.",
			"Dependency hell: You're the dependency from hell.",
			"Package not found: Neither is your talent.",
			"ERESOLVE: Can't resolve your incompetence.",
			"Peer dependency conflict: Your existence is a conflict.",
			"Funding request: Fund your education first.",
		},
		"npm start": {
			"Start script failed: Can't start what's broken.",
			"Port already in use: By someone competent.",
			"Module not found: Neither is your ability.",
		},
		"npm run": {
			"Script not found: Neither is your skill.",
			"Build failed: Can't build on a foundation of failure.",
			"Test failed: Your code is the test, reality failed you.",
		},
		"npm test": {
			"Tests failed: 0% passing, 100% crying.",
			"Test suite disaster: Every assertion asserts your failure.",
			"Coverage 0%: Covered in incompetence though.",
		},

		// Python operations
		"python": {
			"Python execution failed: The snake bit back.",
			"ModuleNotFoundError: Module 'brain' not found.",
			"SyntaxError: Invalid syntax, invalid developer.",
			"IndentationError: Your career is misaligned too.",
		},
		"pip install": {
			"Pip install failed: Package manager managing disappointment.",
			"Requirements not met: Competence requirement not met.",
			"Dependency resolution impossible: Like resolving to make you competent.",
		},

		// Rust operations
		"cargo build": {
			"Build failed: Rust compiles. You don't.",
			"Borrow checker says no: You can't borrow competence.",
			"Lifetime error: Your career lifetime is expiring.",
			"Type mismatch: Expected developer, found disaster.",
		},
		"cargo run": {
			"Run failed: Panic in main thread.",
			"Binary execution failed: Your execution is always flawed.",
		},

		// Database operations
		"mysql": {
			"MySQL error: My SQL, your hell.",
			"Connection refused: Database has self-respect.",
			"Access denied: Denied access to success.",
		},
		"psql": {
			"Postgres error: Post-gres, pre-disaster.",
			"Connection failed: Can't connect competence to you.",
		},

		// Make/Build operations
		"make": {
			"Make failed: Make better choices.",
			"Target not found: Your target of competence not found.",
			"Recipe failed: Recipe for disaster succeeded though.",
		},
		"cmake": {
			"CMake error: Can't make sense of you.",
			"Configuration failed: You're misconfigured.",
		},

		// SSH operations
		"ssh": {
			"Connection refused: Server protecting itself.",
			"Permission denied: Your credentials are insufficient.",
			"Host key verification failed: Host doesn't trust you.",
			"Timeout: Server chose silence over your presence.",
		},
	}

	if insults, exists := commandPatterns[pattern]; exists {
		return selectInsult(insults, ctx.FullCommand)
	}

	// Try just the command without subcommand
	if insults, exists := commandPatterns[ctx.Command]; exists {
		return selectInsult(insults, ctx.FullCommand)
	}

	return ""
}

// getArgumentAwareInsult returns insults based on command arguments
func getArgumentAwareInsult(ctx SmartFallbackContext) string {
	// Check for specific argument patterns
	fullCmd := strings.ToLower(ctx.FullCommand)

	// Force operations
	if strings.Contains(fullCmd, " -f") || strings.Contains(fullCmd, "--force") {
		return selectInsult([]string{
			"Force flag detected: Forcing failure down everyone's throat.",
			"--force won't force competence into you.",
			"Force push to production: Force push to unemployment.",
			"Forcing it won't make it work. Like your career.",
		}, ctx.FullCommand)
	}

	// Sudo operations
	if strings.HasPrefix(fullCmd, "sudo") {
		return selectInsult([]string{
			"Sudo failed: Superuser can't grant super-competence.",
			"Even root privileges can't fix your code.",
			"Sudo make me a developer: Permission denied.",
			"With great power comes great responsibility. You have neither.",
		}, ctx.FullCommand)
	}

	// Recursive operations
	if strings.Contains(fullCmd, " -r") || strings.Contains(fullCmd, "--recursive") {
		return selectInsult([]string{
			"Recursive fail: Failing recursively at every level.",
			"-r flag: Recursively destroying everything.",
			"Recursion depth exceeded: By your incompetence.",
		}, ctx.FullCommand)
	}

	// Verbose operations
	if strings.Contains(fullCmd, " -v") || strings.Contains(fullCmd, "--verbose") {
		return selectInsult([]string{
			"Verbose mode: More output, same failure.",
			"--verbose showing verbose failure details.",
			"Verbose mode: Because watching you fail in detail is entertaining.",
		}, ctx.FullCommand)
	}

	// Help flags
	if strings.Contains(fullCmd, " --help") || strings.Contains(fullCmd, " -h") {
		return selectInsult([]string{
			"Reading help? That ship sailed long ago.",
			"--help can't help you now.",
			"Even the help documentation gave up on you.",
			"RTFM: Read The Failed Manual you just failed.",
		}, ctx.FullCommand)
	}

	// Version checks
	if strings.Contains(fullCmd, "--version") || strings.Contains(fullCmd, "-v") {
		return selectInsult([]string{
			"Checking version: Your version is deprecated.",
			"--version: v0.0.0-incompetent",
			"Software version: Current. Developer version: Obsolete.",
		}, ctx.FullCommand)
	}

	return ""
}

// Tier 1 Intelligence Detection Functions

// getEnvironmentInsult returns insults based on environment variables
func getEnvironmentInsult(ctx SmartFallbackContext) string {
	envInsults := map[string][]string{
		"CI": {
			"Breaking CI? Breaking everyone's day.",
			"CI failure: Continuous Incompetence detected.",
			"Failed in CI: Failing Continuously and Immediately.",
			"CI pipeline broken: Your career pipeline next.",
		},
		"GITHUB_ACTIONS": {
			"GitHub Actions failed: Your actions speak louder than words.",
			"Actions workflow broken: Action item: Find new career.",
			"GitHub runner quit: Running from your code.",
		},
		"GITLAB_CI": {
			"GitLab CI failed: Lab results show terminal incompetence.",
			"Pipeline failed: Pipe down, you're done.",
		},
		"JENKINS_HOME": {
			"Jenkins build failed: Job security failed too.",
			"Jenkins says no: Automated rejection system working.",
		},
		"DEBUG": {
			"Debug mode active: Can't debug your brain.",
			"Debugging? You ARE the bug.",
		},
		"PRODUCTION": {
			"Production error: Producing only failures.",
			"PROD failure: Professional Regression Of Development.",
			"Testing in production? Testing everyone's patience.",
		},
		"NODE_ENV": {
			"Node environment error: Environment of incompetence.",
		},
	}

	for env, val := range ctx.Environment {
		if insults, exists := envInsults[env]; exists && val != "" {
			return selectInsult(insults, ctx.FullCommand)
		}
	}

	return ""
}

// getTimeOfDayInsult returns insults based on time of day
func getTimeOfDayInsult(ctx SmartFallbackContext) string {
	hour := ctx.TimeOfDay

	switch {
	case hour >= 0 && hour < 6:
		return selectInsult([]string{
			"3 AM debugging? Tomorrow won't fix today's code.",
			"Coding at 3 AM? Your code is as tired as you.",
			"Late night failure: Sleep won't fix this.",
			"Midnight coding: Both your code and judgment are impaired.",
		}, ctx.FullCommand)
	case hour >= 6 && hour < 9:
		return selectInsult([]string{
			"Morning failure sets the tone for the day.",
			"Failed before breakfast: Hungry for failure.",
			"Early bird gets the worm: Early coder gets the bugs.",
		}, ctx.FullCommand)
	case hour >= 17 && hour < 20:
		return selectInsult([]string{
			"Evening failure: Overtime making more bugs.",
			"5 PM deploy failed: Weekend ruined.",
			"After hours coding: After competence hours too.",
		}, ctx.FullCommand)
	case hour >= 20 && hour < 24:
		return selectInsult([]string{
			"Late night commit: Commit to quitting instead.",
			"Coding past 8 PM? Desperation detected.",
			"Night owl? More like night fail.",
		}, ctx.FullCommand)
	}

	return ""
}

// getWorkingDirInsult returns insults based on working directory
func getWorkingDirInsult(ctx SmartFallbackContext) string {
	wd := strings.ToLower(ctx.WorkingDir)

	dirPatterns := map[string][]string{
		"/tmp": {
			"Coding in /tmp? That's where your code belongs: temporary.",
			"Temp directory for temp solution for temp developer.",
			"/tmp: Temporary directory, permanent failure.",
		},
		"downloads": {
			"Coding in Downloads? Your career is downloading too.",
			"Downloads folder: Downloaded failure.",
		},
		"/var/log": {
			"Working in logs? You belong in error logs.",
			"/var/log: Logging your mistakes for posterity.",
		},
		"/root": {
			"Running as root? Root of all problems.",
			"Root directory? Rooted in incompetence.",
		},
		"desktop": {
			"Desktop coding? Desktop disaster.",
			"Desktop folder: Where careers go to die.",
		},
		"/opt": {
			"Optional directory for optional competence.",
			"/opt: Opted out of skill.",
		},
	}

	for pattern, insults := range dirPatterns {
		if strings.Contains(wd, pattern) {
			return selectInsult(insults, ctx.FullCommand)
		}
	}

	return ""
}

// getFileExtensionInsult returns insults based on file extensions in command
func getFileExtensionInsult(ctx SmartFallbackContext) string {
	if len(ctx.FileExtensions) == 0 {
		return ""
	}

	ext := strings.ToLower(ctx.FileExtensions[0])

	extInsults := map[string][]string{
		".rs": {
			"Rust file failed: Rust in peace, code.",
			".rs: Rust? Your skills are corroded.",
			"Rust compile failed: Oxidized incompetence.",
		},
		".go": {
			"Go file failed: Stop. Don't go. Just don't.",
			".go error: Should've Go-ne into another career.",
			"Go build failed: Go away.",
		},
		".java": {
			"Java failed: Needs more than coffee beans.",
			".java compile error: Java the Hutt-level bloat.",
			"Java exception: You're the exception to competence.",
		},
		".cpp": {
			"C++ failed: C++ you later, career.",
			".cpp segfault: C++ more like C-- --.",
			"C++ error: Can't ++ your skill level.",
		},
		".c": {
			"C compilation failed: C you don't understand C.",
			".c file error: C-riously incompetent.",
		},
		".py": {
			"Python failed: Snake bit back.",
			".py error: Python crying from your code.",
			"Python IndentationError: Career misaligned too.",
		},
		".js": {
			"JavaScript failed: Just awful Script.",
			".js error: Java-Script? Neither Java nor scripted competence.",
		},
		".ts": {
			"TypeScript failed: Type: Disaster.",
			".ts error: TypeScript can't type your chaos.",
		},
		".sh": {
			"Shell script failed: Script kiddie confirmed.",
			".sh error: Shell-shocked by incompetence.",
		},
		".rb": {
			"Ruby failed: More like Rub-y wounds in codebase.",
			".rb error: Ruby gem? Cubic zirconia skill.",
		},
		".php": {
			"PHP failed: Probably Horrible Programming.",
			".php error: PHP stands for Please Help Professional.",
		},
	}

	if insults, exists := extInsults[ext]; exists {
		return selectInsult(insults, ctx.FullCommand)
	}

	return ""
}

// getNumericArgumentInsult returns insults based on numeric arguments
func getNumericArgumentInsult(ctx SmartFallbackContext) string {
	if len(ctx.NumericArgs) == 0 {
		return ""
	}

	num := ctx.NumericArgs[0]

	// Port numbers
	if num >= 1 && num <= 65535 {
		portInsults := map[int][]string{
			22: {
				"Port 22: 22 ways to fail at SSH.",
				"SSH on port 22: Access denied to competence.",
			},
			80: {
				"Port 80: HTTP status 500 Internal User Error.",
				"Port 80: Gateway to failure.",
			},
			443: {
				"Port 443: HTTPS - Hyper Text Tragic Protocol Stupidity.",
				"Port 443: Secure connection to incompetence.",
			},
			3000: {
				"Port 3000: Three thousand problems detected.",
				"Port 3000: Development port for underdeveloped skills.",
			},
			8080: {
				"Port 8080: Eight-zero-eight-zero errors found.",
				"Port 8080: Alternative HTTP, alternative competence (zero).",
			},
			5432: {
				"Port 5432: PostgreSQL rejecting your queries and you.",
				"Port 5432: Postgres? More like Post-regrets.",
			},
			3306: {
				"Port 3306: MySQL - My Structured Query: Why are you coding?",
				"Port 3306: MySQL rejecting your SQL and existence.",
			},
			27017: {
				"Port 27017: MongoDB - More like MongoDON'T.",
				"Port 27017: NoSQL? No skill either.",
			},
		}

		if insults, exists := portInsults[num]; exists {
			return selectInsult(insults, ctx.FullCommand)
		}
	}

	// Chmod values
	if num == 777 {
		return selectInsult([]string{
			"chmod 777: Maximum permissions, minimum security, zero brains.",
			"777: Jackpot of incompetence.",
			"chmod 777: Triple seven, triple failure.",
		}, ctx.FullCommand)
	}
	if num == 666 {
		return "chmod 666: Devil's permission for devilish code."
	}
	if num == 000 || num == 0 {
		return "chmod 000: Like your access to competence."
	}

	// Kill signals
	if num == 9 {
		return selectInsult([]string{
			"kill -9: Killing process. Can't kill your incompetence.",
			"SIGKILL sent: Signal your career is over.",
		}, ctx.FullCommand)
	}

	return ""
}

// getShellInsult returns insults based on shell type
func getShellInsult(ctx SmartFallbackContext) string {
	shell := strings.ToLower(ctx.Shell)

	shellInsults := map[string][]string{
		"bash": {
			"Bash error: Bash your head against keyboard, same result.",
			"Bourne Again Shell: Borne to fail again.",
		},
		"zsh": {
			"Z shell: Z for Zero competence level.",
			"Zsh failure: Last shell of the alphabet, last in skill.",
		},
		"fish": {
			"Fish shell: You're swimming in failure.",
			"Fish error: Fishing for competence, caught nothing.",
		},
		"sh": {
			"Bourne shell: Born to fail.",
			"sh: Should've stayed in the shell.",
		},
		"ksh": {
			"Korn shell: Your code is corny.",
		},
		"csh": {
			"C shell: See? Shell of incompetence.",
		},
	}

	if insults, exists := shellInsults[shell]; exists {
		return selectInsult(insults, ctx.FullCommand)
	}

	return ""
}

// getCommandComplexityInsult returns insults based on command complexity
func getCommandComplexityInsult(ctx SmartFallbackContext) string {
	length := ctx.CommandLength

	if length < 10 {
		return selectInsult([]string{
			"Short command, short career.",
			"Simple command failed: Simply incompetent.",
		}, ctx.FullCommand)
	}

	if length > 100 {
		return selectInsult([]string{
			"Command longer than your employment prospects.",
			"100+ characters: Complexity hiding incompetence.",
			"Long command: Compensating for short skills.",
		}, ctx.FullCommand)
	}

	if ctx.HasPipes && ctx.HasChaining {
		return selectInsult([]string{
			"Pipes AND chaining? Piping chained disasters together.",
			"Complex piped chain: Complexly incompetent.",
		}, ctx.FullCommand)
	}

	if ctx.HasPipes {
		return selectInsult([]string{
			"Pipe fail: Piping garbage to garbage.",
			"Pipeline broken: Like your career pipeline.",
		}, ctx.FullCommand)
	}

	if ctx.HasChaining {
		return selectInsult([]string{
			"Chaining commands: Chaining failures sequentially.",
			"&& operator: AND you're terrible AND incompetent.",
		}, ctx.FullCommand)
	}

	return ""
}

// Tier 2 Intelligence Detection Functions

// detectGitBranch detects the current git branch
func detectGitBranch() string {
	// Check if we're in a git repo first
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return ""
	}

	// Try to get current branch
	if data, err := os.ReadFile(".git/HEAD"); err == nil {
		head := string(data)
		// Format: "ref: refs/heads/branch-name"
		if strings.HasPrefix(head, "ref: refs/heads/") {
			branch := strings.TrimPrefix(head, "ref: refs/heads/")
			return strings.TrimSpace(branch)
		}
	}

	return ""
}

// detectProjectType detects project type by checking for common project files
func detectProjectType() (string, []string) {
	projectFiles := []string{
		"package.json",      // Node.js
		"Cargo.toml",        // Rust
		"go.mod",            // Go
		"requirements.txt",  // Python
		"Pipfile",           // Python (pipenv)
		"pyproject.toml",    // Python (poetry)
		"pom.xml",           // Java (Maven)
		"build.gradle",      // Java (Gradle)
		"Gemfile",           // Ruby
		"composer.json",     // PHP
		"Makefile",          // C/C++
		"CMakeLists.txt",    // C/C++ (CMake)
		"package.swift",     // Swift
		"mix.exs",           // Elixir
		"Dockerfile",        // Docker
		"docker-compose.yml", // Docker Compose
	}

	var foundFiles []string
	for _, file := range projectFiles {
		if _, err := os.Stat(file); err == nil {
			foundFiles = append(foundFiles, file)
		}
	}

	// Determine project type from found files
	if len(foundFiles) == 0 {
		return "", nil
	}

	// Priority order for type detection
	typeMap := map[string]string{
		"package.json":      "node",
		"Cargo.toml":        "rust",
		"go.mod":            "go",
		"requirements.txt":  "python",
		"Pipfile":           "python",
		"pyproject.toml":    "python",
		"pom.xml":           "java",
		"build.gradle":      "java",
		"Gemfile":           "ruby",
		"composer.json":     "php",
		"package.swift":     "swift",
		"mix.exs":           "elixir",
	}

	for _, file := range foundFiles {
		if projectType, exists := typeMap[file]; exists {
			return projectType, foundFiles
		}
	}

	return "generic", foundFiles
}

// getGitBranchInsult returns insults based on git branch
func getGitBranchInsult(ctx SmartFallbackContext) string {
	if ctx.GitBranch == "" {
		return ""
	}

	branch := strings.ToLower(ctx.GitBranch)

	branchInsults := map[string][]string{
		"main": {
			"Breaking main? Breaking everyone's day.",
			"Main branch failure: Main character of disasters.",
			"Failed on main: Mainly incompetent.",
			"Main branch disaster: You're the main problem.",
		},
		"master": {
			"Master branch failure: Master of disasters.",
			"Breaking master: Mastering incompetence.",
			"Master branch error: You've mastered failure.",
		},
		"develop": {
			"Develop branch failed: Develop your skills first.",
			"Development branch: Under-developed skills.",
			"Develop? More like devolve.",
		},
		"dev": {
			"Dev branch failed: Dev-astating incompetence.",
			"Dev environment: Environment of failure.",
		},
		"staging": {
			"Staging failure: Staging your resignation.",
			"Staging branch: Staging area for disaster.",
		},
		"production": {
			"Production branch failed: Producing unemployment.",
			"Prod branch error: Professionally regressive.",
		},
	}

	// Check exact matches first
	if insults, exists := branchInsults[branch]; exists {
		return selectInsult(insults, ctx.FullCommand)
	}

	// Check for patterns
	if strings.Contains(branch, "feature") || strings.HasPrefix(branch, "feat/") {
		return selectInsult([]string{
			"Feature branch: Featured failure.",
			"New feature: Newly incompetent.",
			"Feature branch failed: Feature: Broken. Developer: Broken.",
		}, ctx.FullCommand)
	}

	if strings.Contains(branch, "hotfix") || strings.HasPrefix(branch, "fix/") {
		return selectInsult([]string{
			"Hotfix branch: You ARE the bug.",
			"Hotfix failed: Can't fix what's fundamentally broken: You.",
			"Hotfix? More like hot mess.",
		}, ctx.FullCommand)
	}

	if strings.Contains(branch, "bugfix") || strings.Contains(branch, "bug/") {
		return selectInsult([]string{
			"Bugfix branch: The bug is you.",
			"Fixing bugs? You ARE the bug.",
		}, ctx.FullCommand)
	}

	if strings.Contains(branch, "release") {
		return selectInsult([]string{
			"Release branch failed: Release your grip on keyboard.",
			"Release branch: Releasing disaster into the world.",
		}, ctx.FullCommand)
	}

	if strings.Contains(branch, "test") {
		return selectInsult([]string{
			"Test branch failed: You're the test, reality failed you.",
			"Testing branch: Test results: FAIL.",
		}, ctx.FullCommand)
	}

	return ""
}

// getProjectTypeInsult returns insults based on detected project type
func getProjectTypeInsult(ctx SmartFallbackContext) string {
	if ctx.ProjectType == "" {
		return ""
	}

	projectInsults := map[string][]string{
		"node": {
			"Node project detected. Dependencies: Many. Skills: None.",
			"package.json found: Package of failures.",
			"Node.js project: Node your way out of this one.",
			"npm detected: Node Package Misery.",
		},
		"rust": {
			"Cargo.toml found: Can't cargo your incompetence.",
			"Rust project detected: Rust in peace, code.",
			"Cargo workspace: Working on failure.",
		},
		"go": {
			"go.mod found: Go away.",
			"Go project detected: Should've gone into another career.",
			"Go modules: Modular incompetence.",
		},
		"python": {
			"requirements.txt found: Requirement for skill: UNMET.",
			"Python project detected: Snake bit you back.",
			"Python dependencies: Depending on incompetence.",
			"Virtual environment detected: Can't isolate stupidity.",
		},
		"java": {
			"pom.xml found: Maven project: Mav-un successful.",
			"Java project detected: Needs more than coffee.",
			"Gradle detected: Grade: F. Project: Failed.",
		},
		"ruby": {
			"Gemfile found: Cubic zirconia skills.",
			"Ruby project: Ruby red with embarrassment.",
			"Bundler detected: Bundle of incompetence.",
		},
		"php": {
			"composer.json found: Can't compose competence.",
			"PHP project: Probably Horrible Programming detected.",
		},
		"swift": {
			"Swift project: Swift path to unemployment.",
			"Package.swift found: Swiftly failing.",
		},
		"elixir": {
			"Elixir project: No elixir can cure this.",
			"mix.exs found: Mixed results: All bad.",
		},
	}

	if insults, exists := projectInsults[ctx.ProjectType]; exists {
		return selectInsult(insults, ctx.FullCommand)
	}

	return ""
}

// selectInsult picks an insult using pseudo-random selection
func selectInsult(insults []string, seed string) string {
	if len(insults) == 0 {
		return ""
	}

	hash := 0
	for _, char := range seed {
		hash = hash*31 + int(char)
	}
	if hash < 0 {
		hash = -hash
	}

	return insults[hash%len(insults)]
}
