package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// =======================
// Entry Point
// =======================

func main() {
	rFlag := flag.Int("r", -1, "generate N random integers (N >= 10)")
	iFlag := flag.String("i", "", "input file with integers")
	dFlag := flag.String("d", "", "directory with .txt files")
	flag.Parse()

	switch {
	case *rFlag != -1:
		if err := runRandom(*rFlag); err != nil {
			log.Fatal(err)
		}
	case *iFlag != "":
		if err := runInputFile(*iFlag); err != nil {
			log.Fatal(err)
		}
	case *dFlag != "":
		if err := runDirectory(*dFlag); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Usage: gosort -r N | -i input.txt | -d directory")
	}
}

// =======================
// -r Mode
// =======================

func runRandom(n int) error {
	if n < 10 {
		return errors.New("N must be >= 10")
	}

	numbers := generateRandomNumbers(n)

	fmt.Println("Original numbers:")
	fmt.Println(numbers)

	processAndPrint(numbers)
	return nil
}

// =======================
// -i Mode
// =======================

func runInputFile(filename string) error {
	numbers, err := readNumbersFromFile(filename)
	if err != nil {
		return err
	}

	if len(numbers) < 10 {
		return errors.New("input file must contain at least 10 valid integers")
	}

	fmt.Println("Original numbers:")
	fmt.Println(numbers)

	processAndPrint(numbers)
	return nil
}

// =======================
// -d Mode
// =======================

func runDirectory(dir string) error {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return errors.New("invalid directory")
	}

	// Using your name and student ID in the output directory name
	outDir := dir + "_sorted_Nandana_Subhash_241ADB029"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}

		inputPath := filepath.Join(dir, e.Name())
		numbers, err := readNumbersFromFile(inputPath)
		if err != nil || len(numbers) < 10 {
			continue
		}

		sorted := process(numbers)

		outputPath := filepath.Join(outDir, e.Name())
		if err := writeNumbersToFile(outputPath, sorted); err != nil {
			return err
		}
	}

	fmt.Printf("Sorted files saved in directory: %s\n", outDir)
	return nil
}

// =======================
// Core Processing
// =======================

func processAndPrint(numbers []int) {
	chunks := splitIntoChunks(numbers)

	fmt.Println("\nChunks before sorting:")
	printChunks(chunks)

	sortedChunks := sortChunksConcurrently(chunks)

	fmt.Println("\nChunks after sorting:")
	printChunks(sortedChunks)

	result := mergeSortedChunks(sortedChunks)

	fmt.Println("\nFinal sorted result:")
	fmt.Println(result)
}

func process(numbers []int) []int {
	chunks := splitIntoChunks(numbers)
	sortedChunks := sortChunksConcurrently(chunks)
	return mergeSortedChunks(sortedChunks)
}

// =======================
// Chunking
// =======================

func splitIntoChunks(numbers []int) [][]int {
	n := len(numbers)

	numChunks := int(math.Ceil(math.Sqrt(float64(n))))
	if numChunks < 4 {
		numChunks = 4
	}

	base := n / numChunks
	extra := n % numChunks

	chunks := make([][]int, 0, numChunks)
	index := 0

	for i := 0; i < numChunks; i++ {
		size := base
		if i < extra {
			size++
		}
		if size == 0 {
			continue
		}
		chunks = append(chunks, numbers[index:index+size])
		index += size
	}

	return chunks
}

// =======================
// Concurrent Sorting
// =======================

func sortChunksConcurrently(chunks [][]int) [][]int {
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for i := range chunks {
		go func(i int) {
			defer wg.Done()
			sort.Ints(chunks[i])
		}(i)
	}

	wg.Wait()
	return chunks
}

// =======================
// Merge Logic
// =======================

func mergeSortedChunks(chunks [][]int) []int {
	if len(chunks) == 0 {
		return []int{}
	}
	if len(chunks) == 1 {
		return chunks[0]
	}

	result := make([]int, 0)
	indices := make([]int, len(chunks))
	
	// Calculate total length for efficiency
	totalLength := 0
	for _, chunk := range chunks {
		totalLength += len(chunk)
	}
	result = make([]int, 0, totalLength)

	for {
		minIdx := -1
		minVal := 0
		first := true

		for i := range chunks {
			if indices[i] >= len(chunks[i]) {
				continue
			}
			val := chunks[i][indices[i]]
			if first || val < minVal {
				minVal = val
				minIdx = i
				first = false
			}
		}

		if minIdx == -1 {
			break
		}

		result = append(result, minVal)
		indices[minIdx]++
	}

	return result
}

// =======================
// Helpers
// =======================

func generateRandomNumbers(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rand.Intn(1000) // range: 0–999
	}
	return nums
}

func readNumbersFromFile(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var nums []int
	scanner := bufio.NewScanner(file)
	line := 0

	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		var n int
		_, err := fmt.Sscanf(text, "%d", &n)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at line %d", line)
		}
		nums = append(nums, n)
	}

	return nums, scanner.Err()
}

func writeNumbersToFile(filename string, numbers []int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, n := range numbers {
		fmt.Fprintln(w, n)
	}
	return w.Flush()
}

func printChunks(chunks [][]int) {
	for i, c := range chunks {
		fmt.Printf("Chunk %d: %v\n", i, c)
	}
}
