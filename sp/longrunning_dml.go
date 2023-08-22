package sp

import (
	"context"
	"log"
	"os"
	"sync/atomic"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func LongRunningMutation(args SpannerArguments) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, args.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	query, err := os.ReadFile(args.File)
	if err != nil {
		log.Fatal(err)
	}
	stmt := spanner.Statement{
		SQL: string(query),
	}
	var affectedRows int64 = -1
	// for monitoring
	loopCount := 0
	// as long as there are rows affected...
	for affectedRows != 0 {
		loopCount++
		// the function in the transaction could be called multiple times
		// so we need to collect the total count first
		var txRows int64
		_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, rwt *spanner.ReadWriteTransaction) error {
			iter := rwt.QueryWithStats(ctx, stmt)
			// although we do not read the rows we need to drain the iterator properly
			defer iter.Stop()
			for {
				_, err := iter.Next()
				atomic.AddInt64(&txRows, iter.RowCount)
				if err == iterator.Done {
					break
				}
				// if the query fails then abort
				if err != nil {
					return err
				}
			}
			return nil
		})
		// if the query fails then abort
		if err != nil && err != iterator.Done {
			return err
		} else {
			log.Printf("loop count:%d affected rows:%v\n", loopCount, txRows)
			// update the loop condition var
			affectedRows = txRows
		}
	}
	return nil
}
