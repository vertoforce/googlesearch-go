package google

import "testing"

import "context"

func TestQuery(t *testing.T) {
	results, err := Query(context.Background(), &Search{Q: `Test`})
	if err != nil {
		t.Error(err)
	}
	if len(results) == 0 {
		t.Errorf("No results")
	}
}
