# cURL package installer - Local Web server
[![Build Status](https://travis-ci.org/electrocucaracha/pkg-mgr.png)](https://travis-ci.org/electrocucaracha/pkg-mgr)
[![Go Report Card](https://goreportcard.com/badge/github.com/electrocucaracha/pkg-mgr)](https://goreportcard.com/report/github.com/electrocucaracha/pkg-mgr)
[![GoDoc](https://godoc.org/github.com/electrocucaracha/pkg-mgr?status.svg)](https://godoc.org/github.com/electrocucaracha/pkg-mgr)

This project provides a Web Server for collecting information about
package installation metrics and also it can be used to centralize the
installation of Linux packages for different distributions. This
initiative is under development but it can be tested using the
following instructions.

```bash
curl -fsSL http://bit.ly/install_pkg | PKG="docker docker-compose make git" bash
newgrp docker
git clone --depth 1 https://github.com/electrocucaracha/pkg-mgr
cd pkg-mgr/
make install
```

Once it's deployed locally, it's possible to consume the scripts
with the following instruction:

    $ curl -fsSL http://localhost:3000/install_pkg?pkg=terraform | bash

## License

Apache-2.0
