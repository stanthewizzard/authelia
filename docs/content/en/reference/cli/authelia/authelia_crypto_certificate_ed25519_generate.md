---
title: "authelia crypto certificate ed25519 generate"
description: "Reference for the authelia crypto certificate ed25519 generate command."
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

## authelia crypto certificate ed25519 generate

Generate an Ed25519 private key and certificate

### Synopsis

Generate an Ed25519 private key and certificate.

This subcommand allows generating an Ed25519 private key and certificate.

```
authelia crypto certificate ed25519 generate [flags]
```

### Examples

```
authelia crypto certificate ed25519 request --help
```

### Options

```
      --ca                            create the certificate as a certificate authority certificate
  -c, --common-name string            certificate common name
      --country strings               certificate country
  -d, --directory string              directory where the generated keys, certificates, etc will be stored
      --duration duration             duration of time the certificate is valid for (default 8760h0m0s)
      --extended-usage strings        specify the extended usage types of the certificate
      --file.ca-certificate string    certificate authority certificate to use when signing this certificate (default "ca.public.crt")
      --file.ca-private-key string    certificate authority private key to use to signing this certificate (default "ca.private.pem")
      --file.certificate string       name of the file to export the certificate data to (default "public.crt")
      --file.private-key string       name of the file to export the private key data to (default "private.pem")
  -h, --help                          help for generate
  -l, --locality strings              certificate locality
      --not-before string             earliest date and time the certificate is considered valid formatted as Jan 2 15:04:05 2006 (default is now)
  -o, --organization strings          certificate organization (default [Authelia])
      --organizational-unit strings   certificate organizational unit
      --path.ca string                source directory of the certificate authority files, if not provided the certificate will be self-signed
  -p, --postcode strings              certificate postcode
      --province strings              certificate province
      --sans strings                  subject alternative names
      --signature string              signature algorithm for the certificate (default "SHA256")
  -s, --street-address strings        certificate street address
```

### SEE ALSO

* [authelia crypto certificate ed25519](authelia_crypto_certificate_ed25519.md)	 - Perform Ed25519 certificate cryptographic operations

