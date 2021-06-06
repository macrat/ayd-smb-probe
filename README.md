ayd-smb-probe
=============

SMB protocol plugin for [Ayd?](https://github.com/macrat/ayd) status check service.


## Install

1. Download binary from [release page](https://github.com/macrat/ayd-smb-probe/releases).

2. Save downloaded binary as `ayd-smb-probe` to somewhere directory that registered to PATH.


## Usage

``` shell
$ ayd smb://username:password@hostname.example.com
```

This example is check if can login to `\\hostname.example.com` with `username` and `password`.

If omit username, it use `guest` as username.
