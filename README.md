# Astra
Simple Web UI for boltdb.
Astra is built using Go's [embed](https://pkg.go.dev/embed) package.
This enables Astra to run as a single binary.

Follow the instructions to install the app and enjoy exploring the boltdb.


## Installation

#### To install the latest version of **Astra** run:

`go install -v github.com/alwindoss/astra/cmd/astra@latest`



#### To run the installed version of **Astra** assuming your $GOBIN is added to the $PATH environment variable

Run `astra` from anywhere in the terminal

#### If you want to view and manage your boltdb using **Astra** then using the flags as shown below

Run `astra -loc=/tmp/astra -name=dummy.db`

Run `astra -h` to see the options
