package main

import (
	"fmt"
	"parrot/internal/llm"
)

func main() {
	fmt.Println("Parrot Insult System Benchmark")
	fmt.Println("================================")

	// Create benchmark
	benchmark := llm.NewBenchmark()

	fmt.Printf("Loading benchmark with %d samples...\n\n", len(benchmark.Samples))

	// Initialize ensemble system
	db := llm.NewInsultDatabase()
	scorer := llm.NewInsultScorer(db)
	hist := llm.NewInsultHistory(20)
	ensemble := llm.NewEnsembleSystem(db, scorer, hist)

	fmt.Println("Training ensemble system...")
	ensemble.Train()
	fmt.Println("Training complete!")

	// Run benchmark
	fmt.Println("Running benchmark...")
	results := benchmark.EvaluateSystem(ensemble)

	// Print results
	fmt.Println()
	results.Print()

	// Print detailed sample results
	fmt.Println("\nDetailed Sample Results:")
	fmt.Println("========================")

	for i, score := range results.DetailedScores {
		if i >= 10 { // Show first 10
			fmt.Printf("... and %d more samples\n", len(results.DetailedScores)-10)
			break
		}

		sample := benchmark.Samples[i]
		fmt.Printf("Sample: %s (%s)\n", sample.ID, sample.Description)
		fmt.Printf("  Command: %s\n", sample.Command)
		fmt.Printf("  Generated: %s\n", score.GeneratedInsult)
		fmt.Printf("  Relevance: %.3f | Latency: %v | Method: %s\n",
			score.Relevance, score.Latency, score.Method)
		fmt.Println()
	}

	// Summary statistics
	fmt.Println("\nAnalysis:")
	fmt.Println("=========")

	if results.AvgRelevance < 0.6 {
		fmt.Println("⚠️  Low relevance score - need better context matching")
	} else if results.AvgRelevance < 0.75 {
		fmt.Println("⚡ Moderate relevance - room for improvement")
	} else {
		fmt.Println("✅ Good relevance scores!")
	}

	if results.FallbackRate > 0.3 {
		fmt.Println("⚠️  High Markov fallback rate - database may need expansion")
	} else {
		fmt.Println("✅ Low fallback rate - good database coverage")
	}

	if results.DiversityScore < 0.8 {
		fmt.Println("⚠️  Low diversity - seeing too many similar insults")
	} else {
		fmt.Println("✅ Good diversity in selections")
	}

	fmt.Println("\nBenchmark complete!")
}
