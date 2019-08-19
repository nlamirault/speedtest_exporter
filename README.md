# speedtest_exporter

[![License Apache 2][badge-license]](LICENSE)
[![GitHub version](https://badge.fury.io/gh/nlamirault%2Fspeedtest_exporter.svg)](https://badge.fury.io/gh/nlamirault%2Fspeedtest_exporter)

* Master : [![Circle CI](https://circleci.com/gh/nlamirault/speedtest_exporter/tree/master.svg?style=svg)](https://circleci.com/gh/nlamirault/speedtest_exporter/tree/master)
* Develop : [![Circle CI](https://circleci.com/gh/nlamirault/speedtest_exporter/tree/develop.svg?style=svg)](https://circleci.com/gh/nlamirault/speedtest_exporter/tree/develop)

This Prometheus exporter check your network connection. Metrics are :

* Latency
* Download bandwidth
* Upload bandwidth


## Installation

You can download the binaries :

* Architecture amd64 [ [linux](https://bintray.com/artifact/download/nlamirault/oss/speedtest_exporter-0.3.0_linux_amd64) / [darwin](https://bintray.com/artifact/download/nlamirault/oss/speedtest_exporter-0.3.0_darwin_amd64) / [freebsd](https://bintray.com/artifact/download/nlamirault/oss/speedtest_exporter-0.3.0_freebsd_amd64) / [netbsd](https://bintray.com/artifact/download/nlamirault/oss/speedtest_exporter-0.3.0_netbsd_amd64) / [openbsd](https://bintray.com/artifact/download/nlamirault/oss/speedtest_exporter-0.3.0_openbsd_amd64) / [windows](https://bintray.com/artifact/download/nlamirault/oss/speedtest_exporter-0.3.0_windows_amd64.exe) ]
* Architecture arm [ [linux](https://bintray.com/artifact/download/nlamirault/oss/speedtest_exporter-0.3.0_linux_arm) ]
* Architecture arm64 [ [linux](https://bintray.com/artifact/download/nlamirault/oss/speedtest_exporter-0.3.0_linux_arm64) ]


## Usage

Launch the Prometheus exporter :

```bash
$ speedtest_exporter -log.level=debug
```

## Development

* Initialize environment

```bash
$ make init
```

* Build tool :

```bash
$ make build
```

* Launch unit tests :

```bash
$ make test
```

## Local Deployment

* Launch Prometheus using the configuration file in this repository:

```bash
$ prometheus -config.file=prometheus.yml
```

* Launch exporter:

```bash
$ speedtest_exporter -log.level=debug
```

* Check that Prometheus find the exporter on `http://localhost:9090/targets`


## Contributing

See [CONTRIBUTING](CONTRIBUTING.md).


## License

See [LICENSE](LICENSE) for the complete license.


## Changelog

A [changelog](ChangeLog.md) is available


## Contact

Nicolas Lamirault <nicolas.lamirault@gmail.com>

[badge-license]: https://img.shields.io/badge/license-Apache2-green.svg?style=flat
