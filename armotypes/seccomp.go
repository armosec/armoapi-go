package armotypes

type SeccompStatus int

const (
	SeccompStatusUnknown            SeccompStatus = 0
	SeccompStatusMissingRuntimeInfo SeccompStatus = 1
	SeccompStatusMissing            SeccompStatus = 2
	SeccompStatusOverlyPermissive   SeccompStatus = 3
	SeccompStatusOptimized          SeccompStatus = 4
	SeccompStatusMisconfigured      SeccompStatus = 5
)

// MandatorySeccompSyscalls are syscalls that are nondeterministically observed by the
// runtime tracer, so they are always added to a generated profile and are not counted
// as unused when scoring one.
var MandatorySeccompSyscalls = []string{"epoll_wait", "tgkill", "sched_yield"}

// SeccompRuntimeBaselineSyscalls are the syscalls a container needs in order to be
// created and reach its entrypoint at all: the container runtime's init work (runc),
// the dynamic loader, and libc process startup.
//
// A generated profile must always allow these. They execute once, before or during
// container startup, so a tracer that attaches to a running container can miss them
// entirely. A workload that then starts once and idles produces a profile without
// them, and applying that profile makes the container unable to start:
//
//	runc create failed: unable to start container process: error during container init:
//	error closing exec fds: get handle to /proc/thread-self/fd: fstatfs ... operation not permitted
//
// Allowing this set trades a small amount of tightness for the guarantee that a
// generated profile can never prevent a workload from starting. It is far narrower
// than the container runtime's default profile, which allows roughly 300 syscalls.
//
// Architecture scope: this is the baseline for the architectures generated profiles
// declare today (SCMP_ARCH_X86_64/X86/X32, see the seccomp profile generator), so
// x86-specific names such as arch_prctl and getpgrp are correct here. Names that do
// not resolve on the running kernel are skipped by the deployed runtime rather than
// rejected (measured: profiles carrying an unresolvable name applied and started).
// If generated profiles ever declare arm64, key this baseline by architecture.
var SeccompRuntimeBaselineSyscalls = []string{
	// container runtime init and process setup
	"arch_prctl", "capget", "capset", "chdir", "clone", "clone3", "execve", "execveat",
	"exit", "exit_group", "prctl", "set_robust_list", "set_tid_address", "setgid",
	"setgroups", "setuid", "umask", "vfork", "wait4",
	// memory
	"brk", "madvise", "mmap", "mprotect", "mremap", "munmap",
	// file descriptors and metadata, used by runc when closing exec fds
	"close", "close_range", "dup2", "dup3", "fcntl", "fstat", "fstatfs", "lseek",
	"newfstatat", "stat", "statfs", "statx",
	// filesystem access performed by the loader and libc startup; the dynamic
	// loader reads ELF and program headers with pread64 and resolves paths with
	// readlink before the entrypoint runs
	"access", "faccessat2", "getcwd", "getdents64", "open", "openat", "openat2",
	"pread64", "read", "readlink", "readlinkat", "write",
	// signals
	"rt_sigaction", "rt_sigprocmask", "rt_sigreturn", "sigaltstack",
	// identity and misc startup queries
	"getegid", "geteuid", "getgid", "getpgrp", "getpid", "getppid", "getrandom",
	"gettid", "getuid", "ioctl", "sysinfo", "uname",
	// scheduling and synchronization primitives used from the first instruction;
	// glibc 2.35+ registers every thread with rseq during startup
	"futex", "nanosleep", "pipe2", "rseq",
}

type SeccompWorkload struct {
	Name                     string                   `json:"name"`
	Kind                     string                   `json:"kind"`
	Namespace                string                   `json:"namespace"`
	ClusterName              string                   `json:"clusterName"`
	K8sResourceHash          string                   `json:"k8sResourceHash"`
	ProfileStatus            SeccompStatus            `json:"profileStatus"`
	SyscallsUsedCount        int                      `json:"syscallsUsedCount"`
	SyscallsUnusedCount      int                      `json:"syscallsUnusedCount"`
	SyscallsUsed             []string                 `json:"syscallsUsed"`
	SyscallUnused            []string                 `json:"syscallsUnused"`
	MissingRuntimeInfoReason MissingRuntimeInfoReason `json:"missingRuntimeInfoReason"`
}
