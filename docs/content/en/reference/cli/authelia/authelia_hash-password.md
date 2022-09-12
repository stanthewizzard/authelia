---
title: "authelia hash-password"
description: "Reference for the authelia hash-password command."
lead: ""
date: 2022-06-15T17:51:47+10:00
draft: false
images: []
menu:
  reference:
    parent: "cli-authelia"
weight: 330
toc: true
---

## authelia hash-password

Hash a password to be used in file-based users database

### Synopsis

Hash a password to be used in file-based users database.

```
authelia hash-password [flags] -- <password>
```

### Examples

```
authelia hash-password -- 'mypass'
authelia hash-password --sha512 -- 'mypass'
authelia hash-password --iterations=4 -- 'mypass'
authelia hash-password --memory=128 -- 'mypass'
authelia hash-password --parallelism=1 -- 'mypass'
authelia hash-password --key-length=64 -- 'mypass'
```

### Options

```
  -c, --config strings    Configuration files
  -h, --help              help for hash-password
  -i, --iterations int    set the number of hashing iterations (default 3)
  -k, --key-length int    [argon2id] set the key length param (default 32)
  -m, --memory int        [argon2id] set the amount of memory param (in MB) (default 64)
  -p, --parallelism int   [argon2id] set the parallelism param (default 4)
  -s, --salt string       set the salt string
  -l, --salt-length int   set the auto-generated salt length (default 16)
  -z, --sha512            use sha512 as the algorithm (changes iterations to 50000, change with -i)
```

### SEE ALSO

* [authelia](authelia.md)	 - authelia untagged-unknown-dirty (master, unknown)

