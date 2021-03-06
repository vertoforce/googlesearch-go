package google

import (
	"context"
	"os"
	"testing"
)

func TestQueryLive(t *testing.T) {
	results, err := Query(context.Background(), &Search{Q: `Test`, Num: 50, Start: 200})
	if err != nil {
		t.Error(err)
	}
	if len(results) <= 5 {
		t.Errorf("No results")
	}
}

func TestQueryLiveTryHard(t *testing.T) {
	results, err := Query(context.Background(), &Search{Q: `Test`, TryHard: true})
	if err != nil {
		t.Error(err)
	}
	if len(results) == 0 {
		t.Errorf("No results")
	}
}

func TestQueryLocal(t *testing.T) {
	localFile, err := os.Open("out.html")
	if err != nil {
		t.Error(err)
		return
	}
	results, err := parse(localFile)
	if err != nil {
		t.Error(err)
	}
	if len(results) == 0 {
		t.Errorf("No results")
	}
}
