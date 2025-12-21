package llm

// InsultExpansionV3 - ROUND 3: Kubernetes, Terraform/Cloud, AI/ML, and HTTP/API insults
// Because DevOps and ML disasters deserve their own category
var InsultExpansionV3 = map[string][]string{

	// ==================== KUBERNETES (45 insults) ====================
	"kubernetes": {
		// Pod failures
		"CrashLoopBackOff: when even the cluster gives up on you.",
		"Pod evicted: your code wasn't paying rent in the cluster.",
		"ImagePullBackOff: Docker Hub is ghosting you today.",
		"OOMKilled: your memory management is the real disaster here.",
		"Pod pending forever: even Kubernetes can't figure out your mess.",
		"Container terminated: the only sensible exit strategy.",
		"Init container failed: couldn't even start failing properly.",
		"Liveness probe failed: your pod is dead, just like your skills.",
		"Readiness probe failed: your app wasn't ready. Neither are you.",
		"Your pod is restarting more than your failed relationships.",

		// Deployment disasters
		"Deployment rollback triggered: even Kubernetes wants to undo you.",
		"ReplicaSet can't scale: your code doesn't scale either.",
		"Rolling update stuck: stuck like your career.",
		"kubectl apply failed: YAML isn't that hard. For most people.",
		"HPA gave up: your app can't handle success anyway.",
		"PodDisruptionBudget violated: budget for disaster was exceeded.",
		"Your deployment strategy is 'recreate everything and pray.'",
		"Surge capacity exceeded: your mistakes surge faster than your pods.",
		"Deployment deadline exceeded: deadline for competence also passed.",
		"Your rollout is rolling back faster than your resume updates.",

		// Resource issues
		"ResourceQuota exceeded: you exceeded the incompetence quota too.",
		"CPU throttled: your code efficiency matches your brain efficiency.",
		"Memory limit reached: should have limited your ambitions instead.",
		"PersistentVolumeClaim pending: your claims to skill are also pending.",
		"StorageClass not found: class 'Competent Developer' also not found.",
		"No nodes available: nodes are hiding from your workload.",
		"Insufficient resources: for your pods AND your excuses.",
		"Node pressure: your code puts pressure on everything it touches.",
		"Eviction threshold reached: you reached the threshold for employment too.",
		"Your resource requests are fiction. So is your understanding of k8s.",

		// Configuration chaos
		"ConfigMap missing: so is your configuration management skill.",
		"Secret not found: your incompetence is no secret though.",
		"RBAC denied: Role-Based Access says you can't access success.",
		"ServiceAccount error: your account of events is also wrong.",
		"Namespace not found: you're lost in more ways than one.",
		"Context switching failed: between k8s clusters AND competence.",
		"kubeconfig invalid: your config for life is also questionable.",
		"API server unreachable: like your career goals.",
		"etcd timeout: your learning also timed out years ago.",
		"Admission webhook rejected: webhook has better judgment than your hiring manager.",

		// Networking nightmares
		"Service unavailable: like your competence.",
		"Ingress misconfigured: traffic can't find your app. Neither can users.",
		"NetworkPolicy blocking: blocking your code is actually correct.",
		"DNS resolution failed: your code can't even find localhost.",
		"ClusterIP not working: nothing in your cluster works.",
	},

	// ==================== TERRAFORM/CLOUD (35 insults) ====================
	"terraform": {
		// State disasters
		"terraform destroy: finally doing something useful with your infra.",
		"State file corrupted: a metaphor for your career trajectory.",
		"State lock failed: someone else is already fixing your mistakes.",
		"Drift detected: your code drifted from reality long ago.",
		"Backend configuration error: your backend knowledge is also in error.",
		"terraform import failed: can't import competence either.",
		"State refresh error: refreshing won't fix fundamental problems.",
		"Remote state not found: your remote chance of success also not found.",
		"State file too large: like your ego vs your abilities.",
		"Workspace confusion: you're confused in all workspaces.",

		// Provider problems
		"Provider error: even AWS doesn't want to work with you.",
		"API rate limited: your mistakes exceeded the API's patience.",
		"Credentials expired: so did your relevance.",
		"Region not available: neither is your future in DevOps.",
		"Service quota exceeded: quota for bad decisions also exceeded.",
		"Provider version mismatch: your version and 'competent' don't match.",
		"Authentication failed: terraform can tell you're a fraud.",
		"IAM denied: Identity and Access confirms you shouldn't access anything.",
		"Resource not found: your resources for learning also not found.",
		"Provider crashed: looking at your code will do that.",

		// Resource failures
		"Resource creation failed: creation of your career also failed.",
		"Dependency cycle detected: you depend on failure consistently.",
		"Timeout waiting for resource: still waiting for your skill to deploy.",
		"Validation failed: your code failed validation. So did your degree.",
		"Variables undefined: your career path is also undefined.",
		"Output error: the only output is embarrassment.",
		"Module not found: 'successful_deployment' module missing.",
		"Plan failed: your life plan also needs review.",
		"Apply error: apply this to your resume: 'needs improvement.'",
		"Destroy failed: can't even destroy properly. Impressive.",

		// Cloud catastrophes
		"S3 bucket public: your mistakes are also very public.",
		"Lambda timeout: your functions fail as slowly as possible.",
		"EC2 terminated: instance of competence also terminated.",
		"RDS connection refused: database refused your terrible queries.",
		"CloudFormation drift: drifting further from employability.",
	},

	// ==================== AI/ML (40 insults) ====================
	"ai_ml": {
		// GPU/CUDA disasters
		"CUDA out of memory: your model is as bloated as your ego.",
		"GPU not found: your neural network found nothing either.",
		"CUDA version mismatch: mismatch between your skills and requirements too.",
		"cuDNN error: your deep learning is very shallow.",
		"NCCL error: distributed training can't distribute your incompetence.",
		"torch.cuda.is_available() returns False, and so does your career.",
		"GPU utilization 0%: matches your brain utilization.",
		"OOM killer struck: should have killed your model idea first.",
		"Driver version incompatible: you're incompatible with success.",
		"Memory allocation failed: allocate some time for learning basics.",

		// Training failures
		"NaN loss: your gradients vanished like your debugging skills.",
		"Loss not decreasing: your competence isn't increasing either.",
		"Validation loss exploding: your mistakes also explode exponentially.",
		"Overfitting to training data: and overfitting to bad practices.",
		"Underfitting everything: including job requirements.",
		"Gradient explosion: the only thing exploding is your career.",
		"Learning rate too high: ambition too high, skill too low.",
		"Model diverged: diverged from anything resembling ML knowledge.",
		"Early stopping triggered: should have stopped you earlier.",
		"Accuracy stuck at 50%: your model learned to flip a coin.",

		// Model issues
		"Model too large: compensating for something?",
		"Model won't load: brain cells also won't load.",
		"Checkpoint corrupted: your understanding is also corrupted.",
		"Weights initialization failed: your project was doomed from the start.",
		"Architecture makes no sense: designed by throwing layers at the wall.",
		"Batch size too large: bigger isn't always better. Applies to egos too.",
		"Embedding dimension mismatch: dimensions of your confusion also mismatch.",
		"Tokenizer error: can't tokenize your excuses.",
		"Inference failed: your ability to infer solutions also failed.",
		"Model prediction: always wrong. Like your career choices.",

		// Data disasters
		"Dataset not found: your dataset of achievements also empty.",
		"Data loader crashed: crashed harder than your GPU.",
		"Label mismatch: your labels and reality don't match.",
		"Preprocessing failed: pre-thinking also failed.",
		"Data augmentation broke: augmenting garbage gives more garbage.",
		"Feature extraction error: can't extract features from nothing.",
		"Normalization failed: nothing normal about your approach.",
		"Train/test split leaked: your incompetence also leaked everywhere.",
		"Class imbalance: your skills are imbalanced too.",
		"Corrupted samples: sample of your work is also corrupted.",
	},

	// ==================== HTTP/API ERRORS (35 insults) ====================
	"http_errors": {
		// Client errors (4xx)
		"400 Bad Request: your request is as bad as your code.",
		"401 Unauthorized: even the API knows you shouldn't be here.",
		"403 Forbidden: the server has better judgment than your manager.",
		"404 Not Found: your skills are also not found.",
		"405 Method Not Allowed: your methods aren't allowed in production either.",
		"408 Request Timeout: patience for your code also timed out.",
		"409 Conflict: the only thing consistent about you.",
		"410 Gone: like your chances of success.",
		"413 Payload Too Large: your ego is also payload too large.",
		"415 Unsupported Media Type: your code type is also unsupported.",
		"418 I'm a Teapot: you're a disaster.",
		"422 Unprocessable Entity: your code is unprocessable by any brain.",
		"429 Too Many Requests: slow down, the API isn't your therapist.",
		"451 Unavailable For Legal Reasons: your code should also be illegal.",

		// Server errors (5xx)
		"500 Internal Server Error: you broke the server. Congratulations.",
		"501 Not Implemented: like your understanding of REST.",
		"502 Bad Gateway: the server between you and success has crashed.",
		"503 Service Unavailable: like your competence.",
		"504 Gateway Timeout: gateway gave up waiting for your code to work.",
		"505 HTTP Version Not Supported: your version of 'working code' isn't supported.",
		"507 Insufficient Storage: insufficient storage for all your mistakes.",
		"508 Loop Detected: you're stuck in a loop of bad decisions.",
		"511 Network Authentication Required: authenticate your claims to skill first.",

		// curl/wget specific
		"curl: Connection refused: server is refusing your advances.",
		"curl: Could not resolve host: your code can't resolve anything.",
		"wget: Connection timed out: even wget is tired of waiting.",
		"SSL certificate problem: your certificate of competence is also invalid.",
		"Connection reset by peer: peer reviewed your code and reset everything.",
		"Network unreachable: like your career aspirations.",
		"curl: (7) Failed to connect: you fail to connect with success too.",
		"Host not found: hosting your code should also not be found.",
		"Certificate verification failed: your skills failed verification too.",
		"Protocol error: you're speaking the wrong protocol. In life too.",
		"Response too large: larger than your debugging capabilities.",
		"Malformed response: your understanding is also malformed.",
	},

	// ==================== CLOUD PROVIDER SPECIFIC (25 insults) ====================
	"cloud": {
		// AWS
		"AWS bill arrived: your wallet just filed for bankruptcy.",
		"Lambda cold start: your brain also has cold starts.",
		"S3 access denied: denied like your promotion.",
		"DynamoDB throttled: your throughput of good ideas is also limited.",
		"EC2 instance terminated: unlike your employment. For now.",
		"CloudWatch alarm: alarming how bad this is.",
		"ECS task failed: task 'be competent' also failed.",
		"SQS message lost: like your message to the team about testing.",

		// GCP
		"GCP quota exceeded: quota for patience also exceeded.",
		"BigQuery timeout: big questions about your competence too.",
		"Cloud Functions crashed: function 'write_good_code' not defined.",
		"GKE cluster error: cluster of mistakes growing.",

		// Azure
		"Azure outage: your code causes outages too.",
		"Blob storage error: blob of errors in your code.",
		"Azure Functions timeout: functions of your brain also timeout.",
		"App Service failed: your service to the team has also failed.",

		// General cloud
		"Cloud costs: $10,000/month for 'Hello World.'",
		"Auto-scaling scaled to zero: correct assessment of your value.",
		"CDN cache miss: your code misses the point entirely.",
		"Load balancer unhealthy: health check for your code: terminal.",
		"Database connection pool exhausted: pool of excuses also exhausted.",
		"Message queue backed up: backed up like your technical debt.",
		"Container registry error: registering your failures since day one.",
		"VPC misconfigured: Very Poorly Configured.",
		"IAM role missing: role 'competent developer' is also missing.",
	},
}

// init registers V3 categories into the lookup system
func init() {
	// V3 categories are automatically available through GetExpandedFallback
	// which checks InsultExpansionV3 after V2
}
