---
title: "authelia validate-config"
description: "Reference for the authelia validate-config command."
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

## authelia validate-config

Check a configuration against the internal configuration validation mechanisms

### Synopsis

Check a configuration against the internal configuration validation mechanisms.

This subcommand allows validation of the YAML and Environment configurations so that a configuration can be checked
prior to deploying it.

```
authelia validate-config [flags]
```

### Examples

```
authelia validate-config
authelia validate-config --config config.yml
```

### Options

```
  -c, --config strings   configuration files to load (default [configuration.yml])
  -h, --help             help for validate-config
```

### SEE ALSO

* [authelia](authelia.md)	 - authelia untagged-unknown-dirty (master, unknown)

