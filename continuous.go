package google

import "context"

import "time"

import "math/rand"

// QueryContinuous Keeps fetching pages of google,
// WaitTimeBetweenRequests is the sleep time between each request +/- VariationTime.
// Recommended > 10 second wait time with > 5 second variation
func QueryContinuous(ctx context.Context, search *Search, WaitTimeBetweenRequests time.Duration, VariationTime time.Duration) chan Result {
	results := make(chan Result)

	go func() {
		defer close(results)

		// Loop over pages however long we can
		var curResultNum = 0
		for {
			s := *search
			s.Start = curResultNum
			s.Num = 50

			// Do query
			thisResults, err := Query(ctx, &s)
			if err != nil {
				return
			}

			// Check if no more results
			if len(thisResults) == 0 {
				return
			}

			// Send off results
			for _, result := range thisResults {
				select {
				case results <- result:
				case <-ctx.Done():
					return
				}
			}

			// Prepare for next request
			curResultNum += s.Num

			// Sleep
			select {
			case <-ctx.Done():
				return
			case <-time.After(WaitTimeBetweenRequests + time.Duration(float32(VariationTime.Nanoseconds())*(rand.Float32()*2-1))):
			}
		}
	}()

	return results
}
