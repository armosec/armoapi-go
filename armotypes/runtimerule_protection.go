package armotypes

// OpenMatchers is the by-kind union of `opens` matchers aggregated across a set
// of rules' profileDataRequired. It is the input the storage layer turns into
// rule-aware collapse protection: sensitive prefixes (and their ancestors) are
// pinned to literal during profile generalisation, so anomaly rules such as
// R0010 ("unexpected /etc/shadow access") can still distinguish a never-seen
// sensitive path from a generalised wildcard like /etc/⋯. Both storage
// environments consume it — the in-cluster apiserver (fed from the rules CRD via
// config) and the backend (fed from rules persisted in MongoDB).
type OpenMatchers struct {
	Exact    []string
	Prefix   []string
	Suffix   []string
	Contains []string
}

// Empty reports whether no opens matcher was contributed by any rule.
func (m OpenMatchers) Empty() bool {
	return len(m.Exact)+len(m.Prefix)+len(m.Suffix)+len(m.Contains) == 0
}

// UnionOpenProtection aggregates and de-duplicates the `opens` matchers across
// all rules' profileDataRequired into a single OpenMatchers. Rules that declare
// no profileDataRequired, no opens field, or `opens: all` contribute nothing:
// "all" keeps every entry and therefore needs no prefix pinning.
func UnionOpenProtection(rules []RuntimeRule) OpenMatchers {
	var m OpenMatchers
	seen := map[string]struct{}{}
	add := func(dst *[]string, kind byte, v string) {
		if v == "" {
			return
		}
		key := string(kind) + v
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		*dst = append(*dst, v)
	}
	for i := range rules {
		pdr := rules[i].ProfileDataRequired
		if pdr == nil || pdr.Opens == nil || pdr.Opens.All {
			continue
		}
		for _, p := range pdr.Opens.Patterns {
			add(&m.Exact, 'e', p.Exact)
			add(&m.Prefix, 'p', p.Prefix)
			add(&m.Suffix, 's', p.Suffix)
			add(&m.Contains, 'c', p.Contains)
		}
	}
	return m
}
