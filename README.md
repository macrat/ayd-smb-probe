ayd-smb-probe
=============

[![CI](https://github.com/macrat/ayd-smb-probe/actions/workflows/ci.yml/badge.svg)](https://github.com/macrat/ayd-smb-probe/actions/workflows/ci.yml)

SMB protocol plugin for [Ayd?](https://github.com/macrat/ayd) status check service.


## Install

1. Download binary from [release page](https://github.com/macrat/ayd-smb-probe/releases).

2. Save downloaded binary as `ayd-smb-probe` to somewhere directory that registered to PATH.


## Usage

``` shell
$ ayd smb://username:password@hostname.example.com/share/path/to/file
```

This example is check if can access to `\\hostname.example.com\share\path\to\file` with `username` and `password`.

The path to file or directory is optional.

The username and password is also optional. It uses `guest` as username if omitted.
