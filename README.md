# Bolt UI [![CI](https://github.com/boreq/bolt-ui/workflows/CI/badge.svg)][actions]

A web interface which lets you browse your [Bolt database](https://github.com/etcd-io/bbolt).

The program is designed mainly with debugging in mind and provides a simple way
of navigating your database structure. As you may want to temporarily expose
this program on a server and access it from another machine the web interface
is secure by default by using TLS as well as a secure token.

![Bolt UI][screenshot]

## Installation

Currently the easiest way of installing Bolt UI is by using the Go programming
language toolchain:

    $ go install github.com/boreq/bolt-ui/cmd/bolt-ui@latest

## Usage

To view `bolt.database` using Bolt UI execute the following command:

    $ bolt-ui bolt.database

The security features can be disabled by using command line flags if you are
using the program locally.

[actions]: https://github.com/boreq/bolt-ui/actions
[screenshot]: https://user-images.githubusercontent.com/1935975/128639070-6c335b7a-26d9-4575-ae94-2250e31149c1.png
