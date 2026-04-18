package env

import (
	"strings"
	"testing"
)

func TestSort_AlphaAscending(t *testing.T) {
	env := map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}
	res := Sort(env, SortOptions{})
	if res.Entries[0].Key != "APPLE" || res.Entries[2].Key != "ZEBRA" {
		t.Errorf("expected ascending order, got %v", res.Entries)
	}
}

func TestSort_AlphaDescending(t *testing.T) {
	env := map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}
	res := Sort(env, SortOptions{Reverse: true})
	if res.Entries[0].Key != "ZEBRA" {
		t.Errorf("expected descending order, got %v", res.Entries)
	}
}

func TestSort_ByValue(t *testing.T) {
	env := map[string]string{"A": "zebra", "B": "apple", "C": "mango"}
	res := Sort(env, SortOptions{ByValue: true})
	if res.Entries[0].Value != "apple" {
		t.Errorf("expected value sort, got %v", res.Entries)
	}
}

func TestSort_SecretsLast(t *testing.T) {
	env := map[string]string{"API_SECRET": "s", "APP_NAME": "myapp", "DB_PASSWORD": "pass"}
	res := Sort(env, SortOptions{SecretsLast: true})
	last := res.Entries[len(res.Entries)-1]
	if !isSecret(last.Key) {
		t.Errorf("expected secret key last, got %s", last.Key)
	}
}

func TestSort_Total(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	res := Sort(env, SortOptions{})
	if res.Total != 3 {
		t.Errorf("expected total 3, got %d", res.Total)
	}
}

func TestSortResult_Format_MasksSecrets(t *testing.T) {
	env := map[string]string{"API_SECRET": "supersecret", "APP_NAME": "myapp"}
	res := Sort(env, SortOptions{})
	out := res.Format(true)
	if strings.Contains(out, "supersecret") {
		t.Error("expected secret to be masked")
	}
	if !strings.Contains(out, "myapp") {
		t.Error("expected non-secret value to be visible")
	}
}
