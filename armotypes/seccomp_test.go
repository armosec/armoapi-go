package armotypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSeccompRuntimeBaselineCoversContainerStart pins the syscalls a container
// cannot start without. A profile generated without these makes the container fail
// in runc init, before the application runs at all.
func TestSeccompRuntimeBaselineCoversContainerStart(t *testing.T) {
	baseline := make(map[string]bool, len(SeccompRuntimeBaselineSyscalls))
	for _, s := range SeccompRuntimeBaselineSyscalls {
		baseline[s] = true
	}

	// fstatfs is the syscall runc was denied when it failed with
	// "error closing exec fds: get handle to /proc/thread-self/fd: fstatfs ... operation not permitted".
	// The rest are the minimum for exec, memory setup, and libc startup.
	for _, required := range []string{
		"fstatfs", "execve", "mmap", "brk", "read", "write", "close",
		"prctl", "capset", "capget", "arch_prctl", "exit_group",
		"set_tid_address", "rt_sigaction", "rt_sigreturn",
	} {
		assert.True(t, baseline[required],
			"%q must be in SeccompRuntimeBaselineSyscalls: a container cannot start without it", required)
	}
}

func TestSeccompRuntimeBaselineHasNoDuplicates(t *testing.T) {
	seen := make(map[string]bool, len(SeccompRuntimeBaselineSyscalls))
	for _, s := range SeccompRuntimeBaselineSyscalls {
		assert.NotEmpty(t, s, "baseline must not contain empty syscall names")
		assert.False(t, seen[s], "duplicate syscall %q in SeccompRuntimeBaselineSyscalls", s)
		seen[s] = true
	}
}

// TestSeccompRuntimeBaselineIsNarrowerThanRuntimeDefault guards the security tradeoff:
// the baseline exists to keep containers startable, not to become a permissive profile.
// The container runtime default profile allows roughly 300 syscalls.
func TestSeccompRuntimeBaselineIsNarrowerThanRuntimeDefault(t *testing.T) {
	assert.Less(t, len(SeccompRuntimeBaselineSyscalls), 120,
		"baseline should stay a startup floor, not approach the runtime default profile")
}
