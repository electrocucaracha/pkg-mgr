# cURL package installer - Local Web server
[![Build Status](https://travis-ci.org/electrocucaracha/pkg-mgr.png)](https://travis-ci.org/electrocucaracha/pkg-mgr)
[![Go Report Card](https://goreportcard.com/badge/github.com/electrocucaracha/pkg-mgr)](https://goreportcard.com/report/github.com/electrocucaracha/pkg-mgr)
[![GoDoc](https://godoc.org/github.com/electrocucaracha/pkg-mgr?status.svg)](https://godoc.org/github.com/electrocucaracha/pkg-mgr)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This project provides a Web Server for collecting information about
package installation metrics and also it can be used to centralize the
installation of Linux packages for different distributions. This
initiative is under development but it can be tested using the
[All-in-One deployment script](aio.sh).

    $ curl -fsSL https://raw.githubusercontent.com/electrocucaracha/pkg-mgr/master/aio.sh | bash

Once it's dependencies are installed and services are running locally,
it's possible to consume the scripts. The following example shows how
install and configure docker properly:

    $ curl -fsSL http://localhost:3000/install_pkg?pkg=docker | bash
