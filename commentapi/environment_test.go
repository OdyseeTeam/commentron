package commentapi

import "testing"

func TestListCacheKeysSeparateMembershipEnvironments(t *testing.T) {
	empty, production, staging := "", "production", "staging"
	for name, key := range map[string]func(*string) (string, error){
		"comments":   func(environment *string) (string, error) { return (ListArgs{Environment: environment}).Key() },
		"superchats": func(environment *string) (string, error) { return (SuperListArgs{Environment: environment}).Key() },
	} {
		t.Run(name, func(t *testing.T) {
			seen := make(map[string]bool)
			for _, environment := range []*string{nil, &empty, &production, &staging} {
				value, err := key(environment)
				if err != nil {
					t.Fatal(err)
				}
				if seen[value] {
					t.Fatal("membership environments share a list cache entry")
				}
				seen[value] = true
			}
		})
	}
}
