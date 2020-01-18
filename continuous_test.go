package google

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestQueryContinuous(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	results := QueryContinuous(ctx, &Search{Q: "Test"}, time.Second*5, time.Second*5)

	numberOfWantedResults := 200
	numberOfResults := 0
	for result := range results {
		fmt.Println(result.Title)
		numberOfResults++
		if numberOfResults == numberOfWantedResults {
			break
		}
	}

	cancel()
}

func TestTimeVariance(t *testing.T) {
	rand.Seed(time.Now().Unix())
	fmt.Printf("%v\n", rand.Float32()*2-1)
	fmt.Printf("%v\n", time.Duration(float32(time.Second.Nanoseconds())*(rand.Float32()*2-1)))
}
