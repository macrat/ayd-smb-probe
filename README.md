ayd-smb-probe
=============

SMB protocol plugin for [Ayd?](https://github.com/macrat/ayd) status check service.


## Usage

``` shell
$ ayd smb://username:password@hostname.example.com
```

This example is check if can login to `\\hostname.example.com` with `username` and `password`.

If omit username, it use `guest` as username.
