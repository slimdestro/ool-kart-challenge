package coupon

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Index struct {
	mu sync.RWMutex
	// tested with uint64, Processing all three files consumes like 10+ GB or memory.
	// with uint32, it can save upto 3BG but there are chances of hash collision. On Prod, will always go with Redis or similar service
	index map[uint32]byte
	count int64
}

func NewIndex() *Index {
	return &Index{
		index: make(map[uint32]byte),
	}
}

func (ci *Index) Add(hash uint32, fileBit byte) {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	prev := ci.index[hash]
	if prev == 0 {
		ci.count++
	}
	ci.index[hash] = prev | fileBit
}

func (ci *Index) Get(hash uint32) (byte, bool) {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	v, ok := ci.index[hash]
	return v, ok
}

func (ci *Index) Size() int64 {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	return ci.count
}

func hash32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

var couponRegexp = regexp.MustCompile(`(?i)\b[A-Z0-9]{8,10}\b`)

func streamAndIndex(filePath string, fileBit byte, ci *Index) (time.Duration, error) {
	start := time.Now()
	//log.Printf("Indexing file %d started: %s", fileBit, filePath)

	f, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("could not open coupon file %s: %w", filePath, err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return 0, fmt.Errorf("could not create gzip reader for %s: %w", filePath, err)
	}
	defer gz.Close()

	reader := bufio.NewReaderSize(gz, 4*1024*1024)
	scanner := bufio.NewScanner(reader)
	const maxCapacity = 10 * 1024 * 1024
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Text()
		matches := couponRegexp.FindAllString(line, -1)
		for _, raw := range matches {
			code := strings.ToUpper(strings.TrimSpace(raw))
			if len(code) < 8 || len(code) > 10 {
				continue
			}
			h := hash32(code)
			ci.Add(h, fileBit)
		}
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		return 0, fmt.Errorf("scanner error on %s: %w", filePath, err)
	}

	elapsed := time.Since(start)
	//log.Printf("Indexing file %d finished in %s. Current unique coupon count: %d", fileBit, elapsed.String(), ci.Size())
	return elapsed, nil
}

func BuildIndex(filePaths []string) (*Index, time.Duration) {
	ci := NewIndex()
	totalStreamingDuration := time.Duration(0)

	for i, path := range filePaths {
		if i >= 3 {
			log.Fatalf("Fatal: BuildIndex received more than 3 file paths. Logic supports only 3.")
		}

		fileBit := byte(1 << i)

		dur, err := streamAndIndex(path, fileBit, ci)
		if err != nil {
			log.Printf("warning: %v", err)
			continue
		}
		totalStreamingDuration += dur
	}

	return ci, totalStreamingDuration
}

func IsCouponValid(ci *Index, cache *LRUCache, raw string) bool {
	code := strings.ToUpper(strings.TrimSpace(raw))
	if code == "" {
		return true
	}
	if len(code) < 8 || len(code) > 10 {
		return false
	}

	if v, ok := cache.Get(code); ok {
		return v
	}

	h := hash32(code)
	mask, ok := ci.Get(h)
	if !ok {
		cache.Set(code, false)
		return false
	}

	count := 0
	for i := 0; i < 3; i++ {
		if (mask>>i)&1 == 1 {
			count++
		}
	}
	valid := count >= 2

	cache.Set(code, valid)
	return valid
}
