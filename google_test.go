package google

import (
	"context"
	"os"
	"testing"
)

func TestQueryLive(t *testing.T) {
	results, err := Query(context.Background(), &Search{Q: `Test`})
	if err != nil {
		t.Error(err)
	}
	if len(results) == 0 {
		t.Errorf("No results")
	}
}

func TestQueryLocal(t *testing.T) {
	localFile, err := os.Open("out2.html")
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
