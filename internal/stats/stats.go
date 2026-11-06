package stats

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/jmeiracorbal/mac-optimizer/internal/runner"
)

func Print(r *runner.Runner) {
	fmt.Println("\033[0;36m=== System Statistics ===\033[0m")
	fmt.Printf("\033[0;34mMemory free:  %s\033[0m\n", memoryFree(r))
	fmt.Printf("\033[0;34mDisk used:    %s\033[0m\n", diskUsed())
	fmt.Printf("\033[0;34mUptime:       %s\033[0m\n", uptime(r))
	fmt.Println()
}

func memoryFree(r *runner.Runner) string {
	out, err := r.Query("memory_pressure")
	if err != nil {
		return "unavailable"
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "System-wide memory free percentage") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
		}
	}
	return "unavailable"
}

func diskUsed() string {
	var st syscall.Statfs_t
	if err := syscall.Statfs("/", &st); err != nil {
		return "unavailable"
	}
	total := st.Blocks * uint64(st.Bsize)
	free := st.Bfree * uint64(st.Bsize)
	pct := int(float64(total-free) / float64(total) * 100)
	return fmt.Sprintf("%d%%", pct)
}

func uptime(r *runner.Runner) string {
	out, err := r.Query("uptime")
	if err != nil {
		return "unavailable"
	}
	s := strings.TrimSpace(string(out))
	if idx := strings.Index(s, "up "); idx >= 0 {
		rest := s[idx+3:]
		if end := strings.Index(rest, ", load"); end >= 0 {
			return strings.TrimSpace(rest[:end])
		}
		if end := strings.Index(rest, " user"); end >= 0 {
			return strings.TrimSpace(rest[:end])
		}
	}
	return s
}
