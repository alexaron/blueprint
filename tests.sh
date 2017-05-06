# Runs all tests (Linux) excluding vendor directory
go test $(go list ./... | grep -v /vendor/)