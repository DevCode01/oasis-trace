package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
)

// Matches Apache/Nginx combined log format:
// 127.0.0.1 - - [10/Oct/2026:13:55:36 +0000] "GET /index.html HTTP/1.1" 200 2326
var logLineRe = regexp.MustCompile(
	`^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) (\S+) \S+" (\d{3}) (\d+)`,
)

type stats struct {
	total       int
	statusCount map[string]int
	ipCount     map[string]int
	pathCount   map[string]int
	hourCount   map[string]int
	bytesTotal  int64
}

func newStats() *stats {
	return &stats{
		statusCount: map[string]int{},
		ipCount:     map[string]int{},
		pathCount:   map[string]int{},
		hourCount:   map[string]int{},
	}
}

func (s *stats) ingest(line string) bool {
	m := logLineRe.FindStringSubmatch(line)
	if m == nil {
		return false
	}
	ip := m[1]
	timestamp := m[2]
	path := m[4]
	status := m[5]
	size, _ := strconv.ParseInt(m[6], 10, 64)

	s.total++
	s.ipCount[ip]++
	s.pathCount[path]++
	s.statusCount[status[:1]+"xx"]++
	s.bytesTotal += size

	// timestamp like 10/Oct/2026:13:55:36 +0000 -> hour is chars 13:15
	if len(timestamp) >= 15 {
		hour := timestamp[12:14]
		s.hourCount[hour]++
	}
	return true
}

func topN(m map[string]int, n int) []struct {
	Key   string
	Count int
} {
	type kv struct {
		Key   string
		Count int
	}
	list := make([]kv, 0, len(m))
	for k, v := range m {
		list = append(list, kv{k, v})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Count != list[j].Count {
			return list[i].Count > list[j].Count
		}
		return list[i].Key < list[j].Key
	})
	if len(list) > n {
		list = list[:n]
	}
	result := make([]struct {
		Key   string
		Count int
	}, len(list))
	for i, kv := range list {
		result[i] = struct {
			Key   string
			Count int
		}{kv.Key, kv.Count}
	}
	return result
}

func main() {
	var reader io.Reader = os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open %s: %v\n", os.Args[1], err)
			os.Exit(1)
		}
		defer f.Close()
		reader = f
	}

	s := newStats()
	skipped := 0
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		if !s.ingest(scanner.Text()) {
			skipped++
		}
	}

	if s.total == 0 {
		fmt.Println("No valid log lines found.")
		os.Exit(1)
	}

	fmt.Printf("=== Access Log Report ===\n")
	fmt.Printf("Total requests:    %d (skipped %d unparseable lines)\n", s.total, skipped)
	fmt.Printf("Total bytes sent:  %d\n\n", s.bytesTotal)

	fmt.Println("Status codes:")
	codes := make([]string, 0, len(s.statusCount))
	for c := range s.statusCount {
		codes = append(codes, c)
	}
	sort.Strings(codes)
	for _, c := range codes {
		fmt.Printf("  %-4s %d\n", c, s.statusCount[c])
	}

	fmt.Println("\nTop 5 IPs:")
	for _, kv := range topN(s.ipCount, 5) {
		fmt.Printf("  %-20s %d\n", kv.Key, kv.Count)
	}

	fmt.Println("\nTop 5 paths:")
	for _, kv := range topN(s.pathCount, 5) {
		fmt.Printf("  %-30s %d\n", kv.Key, kv.Count)
	}

	fmt.Println("\nRequests by hour:")
	hours := make([]string, 0, len(s.hourCount))
	for h := range s.hourCount {
		hours = append(hours, h)
	}
	sort.Strings(hours)
	for _, h := range hours {
		count := s.hourCount[h]
		bar := ""
		for i := 0; i < count && i < 50; i++ {
			bar += "#"
		}
		fmt.Printf("  %sh00  %-50s %d\n", h, bar, count)
	}
}
