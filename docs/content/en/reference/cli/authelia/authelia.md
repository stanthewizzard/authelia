---
title: "authelia"
description: "Reference for the authelia command."
lead: ""
date: 2022-06-15T17:51:47+10:00
draft: false
images: []
menu:
  reference:
    parent: "cli-authelia"
weight: 320
toc: true
---

## authelia

authelia untagged-unknown-dirty (master, unknown)

### Synopsis

authelia untagged-unknown-dirty (master, unknown)

An open-source authentication and authorization server providing
two-factor authentication and single sign-on (SSO) for your
applications via a web portal.

Documentation is available at: https://www.authelia.com/

```
authelia [flags]
```

### Examples

```
authelia --config /etc/authelia/config.yml --config /etc/authelia/access-control.yml
authelia --config /etc/authelia/config.yml,/etc/authelia/access-control.yml
authelia --config /etc/authelia/config/
```

### Options

```
  -c, --config strings   configuration files to load
  -h, --help             help for authelia
```

### SEE ALSO

* [authelia access-control](authelia_access-control.md)	 - Helpers for the access control system
* [authelia build-info](authelia_build-info.md)	 - Show the build information of Authelia
* [authelia crypto](authelia_crypto.md)	 - Perform cryptographic operations
* [authelia hash-password](authelia_hash-password.md)	 - Hash a password to be used in file-based users database
* [authelia storage](authelia_storage.md)	 - Manage the Authelia storage
* [authelia validate-config](authelia_validate-config.md)	 - Check a configuration against the internal configuration validation mechanisms

