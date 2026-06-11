package armotypes

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestUnionOpenProtection(t *testing.T) {
	ruleYAML := `
profileDependency: 0
profileDataRequired:
  opens:
    - prefix: "/etc/shadow"
    - contains: "/.ssh/"
    - exact: "/etc/sudoers"
    - suffix: "_key"
`
	allRuleYAML := `
profileDependency: 0
profileDataRequired:
  opens: all
`
	noPdrYAML := `
profileDependency: 0
`

	var r1, r2, r3 RuntimeRule
	for _, tc := range []struct {
		y string
		r *RuntimeRule
	}{{ruleYAML, &r1}, {allRuleYAML, &r2}, {noPdrYAML, &r3}} {
		if err := yaml.Unmarshal([]byte(tc.y), tc.r); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
	}

	m := UnionOpenProtection([]RuntimeRule{r1, r2, r3})

	eq := func(name string, got, want []string) {
		if len(got) != len(want) {
			t.Fatalf("%s = %v, want %v", name, got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("%s = %v, want %v", name, got, want)
			}
		}
	}
	eq("Exact", m.Exact, []string{"/etc/sudoers"})
	eq("Prefix", m.Prefix, []string{"/etc/shadow"})
	eq("Suffix", m.Suffix, []string{"_key"})
	eq("Contains", m.Contains, []string{"/.ssh/"})

	// "all" and no-profileDataRequired rules contribute nothing.
	if UnionOpenProtection([]RuntimeRule{r2, r3}).Empty() != true {
		t.Errorf("expected empty union from 'all' + no-pdr rules")
	}
}
