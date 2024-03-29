# Sensus-mWater Adapter

![Go Tests](https://github.com/fffinkel/sensus-mwater-adapter/actions/workflows/test.yaml/badge.svg) ![Go Coverage Badge](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/fffinkel/bb5d76c3d157a2497d578e1a30564c4a/raw/coverage.json) [![Go Report Card](https://goreportcard.com/badge/github.com/fffinkel/sensus-mwater-adapter)](https://goreportcard.com/report/github.com/fffinkel/sensus-mwater-adapter)

The Sensus-mWater Adapter is a program that is able to translate Sensus EMR
meter reading CSVs to mWater meter reading accounting transactions.

It receives CSVs through HTTP POST requests, translates each line in the CSV to
an mWater transaction, and assembles a final POST request to mWater's API.

## Authentication

Sending clients must use Basic authentication.

The adapter authenticates with the mWater API using a username and password
that are POSTed in exchange for a Client ID.

## Usage

Process a Sensus CSV:

```
$ sensus-mwater-adapter input.csv
```

Process a Sensus CSV, but do not send the mWater HTTP requests.

```
$ sensus-mwater-adapter \
    --to-account abcd0123efgh4567ijkl8910 \
    --from-account abcd0123efgh4567ijkl8910 \
    --username mwateruser \
    --password hunter2 \
    --dry-run \
    input.csv
```

### Configuration

* username: mWater API username
* password: mWater API password
* to-account: Accounts receivable ID
* from-account: Water sales ID
* dry-run: Run the program but do not POST to mWater

### Error Handling

For each CSV the adapter receives, the adapter attempts to parse as much of the
CSV as possible, and constructs a request with all rows it was able to
successfully parse. It logs all rows that it was unable to parse.

### Logging

### Diagnosis

### Replay CSVs

## Testing

From the root of the sensus-mwater-adapter checkout:

```
go test -v ./...
```

Go tests currently run in a GitHub Action on all pushes to main. Pass/fail
status and coverage percentage are shown in badges on this README.

## Contributing

Pull requests are always welcome. Please make sure to update tests as
appropriate.

If you have any questions or comments, please feel free to open a GitHub issue
to discuss.

### Useful Documentation

[mWater API](https://api.mwater.co/)

## License

See [LICENSE](LICENSE)
