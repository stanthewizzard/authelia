---
title: "authelia crypto certificate ed25519 request"
description: "Reference for the authelia crypto certificate ed25519 request command."
lead: ""
date: 2022-06-27T18:27:57+10:00
draft: false
images: []
menu:
  reference:
    parent: "cli-authelia"
weight: 330
toc: true
---

## authelia crypto certificate ed25519 request

Generate an Ed25519 private key and certificate signing request

### Synopsis

Generate an Ed25519 private key and certificate signing request.

This subcommand allows generating an Ed25519 private key and certificate signing request.

```
authelia crypto certificate ed25519 request [flags]
```

### Examples

```
authelia crypto certificate ed25519 request --help
```

### Options

```
  -c, --common-name string            certificate common name
      --country strings               certificate country
  -d, --directory string              directory where the generated keys, certificates, etc will be stored
      --duration duration             duration of time the certificate is valid for (default 8760h0m0s)
      --file.csr string               name of the file to export the certificate request data to (default "request.csr")
      --file.private-key string       name of the file to export the private key data to (default "private.pem")
  -h, --help                          help for request
  -l, --locality strings              certificate locality
      --not-before string             earliest date and time the certificate is considered valid formatted as Jan 2 15:04:05 2006 (default is now)
  -o, --organization strings          certificate organization (default [Authelia])
      --organizational-unit strings   certificate organizational unit
  -p, --postcode strings              certificate postcode
      --province strings              certificate province
      --sans strings                  subject alternative names
      --signature string              signature algorithm for the certificate (default "SHA256")
  -s, --street-address strings        certificate street address
```

### SEE ALSO

* [authelia crypto certificate ed25519](authelia_crypto_certificate_ed25519.md)	 - Perform Ed25519 certificate cryptographic operations

