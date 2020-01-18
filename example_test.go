package google

import (
	"context"
	"fmt"
)

func ExampleQuery() {
	results, _ := Query(context.Background(), &Search{Q: "Test"})
	if len(results) > 0 {
		fmt.Println(results[0].Title)
	}

	// Output: Speedtest by Ookla - The Global Broadband Speed Test
}

func ExampleQuery_tryHard() {
	results, _ := Query(context.Background(), &Search{Q: "Test", TryHard: true})
	if len(results) > 0 {
		fmt.Println(results[0].Title)
	}

	// Output: Speedtest by Ookla - The Global Broadband Speed Test
}
