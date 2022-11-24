# Kernel config checker

Checks that options in a given kernel config match a list of expected values
(Y, N, is not set). This is usefull to check for compliance, for example for
the ANSSI guide: TODO.

## How to

1. Download a kernel config:

```
$ curl ...
```

2. Match it with a config profile:

```
$ go run main.go -c config -p profile -r results.csv
```

## TODO

- Validate sysctls too
- Integrate ideas from <https://github.com/a13xp0p0v/kconfig-hardened-check>
