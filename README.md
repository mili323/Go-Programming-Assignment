# gosort

Concurrent integer sorting in Go using goroutines.

## Run

```bash
go mod init gosort
go build
./gosort 

## Modes

### Random (`-r`)

./gosort -r N

* `N` ≥ 10
* Generates random integers (0–999)
* Prints original numbers, chunks (before/after), and final result

### Input File (`-i`)

./gosort -i input.txt


* One integer per line
* Empty lines ignored
* Invalid lines cause error
* ≥ 10 integers required

### Directory (`-d`)

./gosort -d incoming

* Processes all `.txt` files
* Each file sorted independently
* Output directory:


incoming_sorted_firstname_surname_studentID

## Rules

* Chunks = `max(4, ceil(sqrt(n)))`
* Chunk sizes differ by at most 1
* Each chunk sorted in its own goroutine
* Sorted chunks merged manually (no re-sort)
