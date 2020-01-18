package google

import (
	"context"
	"fmt"
)

func ExampleQuery() {
	search := Search{
		Q: "Test",
	}

	results, _ := Query(context.Background(), &search)
	if len(results) > 0 {
		fmt.Println(results[0].Title)
	}

	// Output: Speedtest by Ookla - The Global Broadband Speed Test
}
