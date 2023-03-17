# sensus-mwater-adapter

![Gitlab code coverage](https://img.shields.io/gitlab/pipeline-coverage/fffinkel/sensus-mwater-adapter?branch=main)

![badge](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/fffinkel/bb5d76c3d157a2497d578e1a30564c4a/raw/coverage.json)

[![Go Report Card](https://goreportcard.com/badge/github.com/fffinkel/sensus-mwater-adapter)](https://goreportcard.com/report/github.com/fffinkel/sensus-mwater-adapter)

## Name

Sensus-mWater Adapter

## Description

The Sensus-mWater Adapter is a program that is able to translate Sensus EMR
meter reading CSVs to mWater meter reading accounting transactions.

It receives CSVs through HTTP POST requests, translates each line in the CSV to
an mWater transaction, and assembles a final POST request to mWater's API.

Sending clients must use Basic authentication.

## Installation

## Usage

### Logging

### Diagnosis

### Error Handling

### Replay CSVs

### Configuration

## Testing

## Contributing

Pull requests are always welcome. Please make sure to update tests as
appropriate.

If you have any questions or comments, please feel free to open a GitHub issue
to discuss.

### Useful Documentation

[mWater API](https://api.mwater.co/)

## License

See [LICENSE](LICENSE)
