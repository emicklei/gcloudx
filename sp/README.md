## spanner utils

## usage

Create a file `mutation.sql` that hold the DML query that requires iterations to
complete.

    gcloudx -v spanner dml -f mutation.sql -d <fully-qualified-spanner-database-id>

See `gcloudx spanner dml -help` for all options.
