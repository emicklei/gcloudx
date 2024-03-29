package sp

import (
	"context"
	"log"
	"os"
	"sync/atomic"
	"time"

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
	if args.Verbose {
		log.Println("SQL statement:", string(query))
		log.Println("transaction timeout per loopis set to", args.Timeout)
	}
	var affectedRows int64 = -1
	var totalRows int64 = 0
	// for monitoring
	loopCount := 0
	// as long as there are rows affected...
	for affectedRows != 0 {
		loopCount++
		now := time.Now()
		// the function in the transaction could be called multiple times
		// so we need to collect the total count first
		var txRows int64

		ctx, cancel := context.WithTimeout(context.Background(), args.Timeout)
		defer cancel()

		if args.PartitionedUpdate {
			log.Println("using partitioned update DML")
			count, err := client.PartitionedUpdate(ctx, stmt)
			if err != nil {
				return err
			}
			atomic.AddInt64(&txRows, count)
		} else {
			_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, rwt *spanner.ReadWriteTransaction) error {
				iter := rwt.QueryWithStats(ctx, stmt)
				// although we do not read the rows we need to drain the iterator properly
				defer iter.Stop()
				// drain the iterator
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
		}

		// if the query fails then abort
		if err != nil && err != iterator.Done {
			return err
		} else {
			log.Printf("loop count:%d affected rows:%v (%v)\n", loopCount, txRows, time.Since(now))
			// update the loop condition var
			affectedRows = txRows
			totalRows += txRows
		}
	}
	log.Printf("total rows affected:%v\n", totalRows)
	return nil
}
