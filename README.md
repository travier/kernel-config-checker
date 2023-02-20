# Kernel config checker

Checks that options in a given kernel config match a list of expected values.
This is usefull to check for compliance, for example with the
[ANSSI-BP-028 v2.0 profile](https://www.ssi.gouv.fr/guide/recommandations-de-securite-relatives-a-un-systeme-gnulinux/).

See also the
[OpenSCAP profiles](https://www.open-scap.org/security-policies/choosing-policy/)
for the rest of the system configuration.

## How to

1. Download a kernel config:

```
# Example from:
# CentOS Stream: https://gitlab.com/redhat/centos-stream/rpms/kernel
# Fedora: https://src.fedoraproject.org/rpms/kernel
$ curl -O https://gitlab.com/redhat/centos-stream/rpms/kernel/-/raw/c9s/kernel-x86_64-rhel.config
$ curl -O https://gitlab.com/redhat/centos-stream/rpms/kernel/-/raw/c9s/kernel-aarch64-rhel.config
$ curl -O https://src.fedoraproject.org/rpms/kernel/raw/f37/f/kernel-x86_64-fedora.config
$ curl -O https://src.fedoraproject.org/rpms/kernel/raw/f37/f/kernel-aarch64-fedora.config
```

2. Match them with a config profile:

```
$ go run main.go \
    -config examples/c9s/kernel-x86_64-rhel.config,examples/c9s/kernel-aarch64-rhel.config \
    -profile profiles/ANSSI-BP-028 \
    > results.csv
...
```

## Interpreting results

Note that it's extremely unlikely that you will have all config options set as
recommended as some options are so strict that they are impractical in most
cases (see `CONFIG_MODULES is not set` for example).

## TODO

- Validate sysctls too
- Add ideas from <https://github.com/a13xp0p0v/kconfig-hardened-check>
