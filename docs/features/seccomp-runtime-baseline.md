---
type: feature
status: active
owner: rotem@armosec.io
scope: repo
related_code:
  - armotypes/seccomp.go
---

# Seccomp runtime baseline syscalls (`SeccompRuntimeBaselineSyscalls`)

The set of syscalls a container needs in order to be created and reach its
entrypoint. Every generated seccomp profile must allow them, on top of whatever
the runtime tracer observed.

## Why it exists

A generated profile used to be `observed syscalls + MandatorySeccompSyscalls`, and
`MandatorySeccompSyscalls` holds three entries. That makes the profile only as
complete as the observation.

Container startup syscalls run once, before or during container init. A tracer that
attaches to an already running container can miss them. Workloads that exec many
short-lived processes keep re-emitting them and look fine; a workload that starts
once and then idles does not. Its profile ends up without `execve`, `mmap`, `read`
or `fstatfs`, and applying that profile stops the container from starting at all:

```
runc create failed: unable to start container process: error during container init:
error closing exec fds: get handle to /proc/thread-self/fd: fstatfs ... operation not permitted
```

The failure is `RunContainerError` with exit code 128, so it happens before any
application code runs. The workload crashloops and the deployment rollout wedges.

## Contract for generators

Union `SeccompRuntimeBaselineSyscalls` into the allow list of every generated
profile. Do not make it conditional on what was observed: the point is that these
syscalls cannot be relied on to appear in the observation.

`MandatorySeccompSyscalls` keeps its separate meaning. Those three are
nondeterministically observed, so they are added to the profile and are also
excluded when scoring a profile as overly permissive.

## Security tradeoff

The baseline is a startup floor, not a permissive profile. It targets the architectures generated profiles declare today (SCMP_ARCH_X86_64/X86/X32); names that do not resolve on the running kernel are skipped by the deployed runtime rather than rejected (measured: a profile carrying an unresolvable name applied and the container started). If generated profiles ever declare arm64, the baseline must become architecture-keyed. It is roughly 70 syscalls against about 300 in the container runtime default profile, so an
enforced profile stays meaningfully tighter than unconfined while never being able
to prevent a workload from starting. A test asserts the list stays well under the
runtime default size, so growth is a deliberate decision rather than drift.

## Groups in the list

| Group | Purpose |
|---|---|
| runtime init and process setup | runc init, credential and namespace setup, `execve` |
| memory | loader and allocator setup |
| file descriptors and metadata | includes `fstatfs`, needed when runc closes exec fds |
| filesystem access | loader and libc startup reads, including pread64 for ELF headers and readlink for path resolution |
| signals | libc signal setup |
| identity and misc startup queries | uid/gid and environment queries at startup |
| scheduling and synchronization | primitives used from the first instruction |
