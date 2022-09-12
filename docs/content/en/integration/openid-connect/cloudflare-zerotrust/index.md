---
title: "Cloudflare Zero Trust"
description: "Integrating Cloudflare Zero Trust with the Authelia OpenID Connect Provider."
lead: ""
date: 2022-06-15T17:51:47+10:00
draft: false
images: []
menu:
  integration:
    parent: "openid-connect"
weight: 620
toc: true
community: true
---

## Tested Versions

* [Authelia]
  * [v4.35.6](https://github.com/authelia/authelia/releases/tag/v4.35.6)

## Before You Begin

You are required to utilize a unique client id and a unique and random client secret for all [OpenID Connect] relying
parties. You should not use the client secret in this example, you should randomly generate one yourself. You may also
choose to utilize a different client id, it's completely up to you.

This example makes the following assumptions:

* __Cloudflare Team Name:__ `example-team`
* __Authelia Root URL:__ `https://auth.example.com`
* __Client ID:__ `cloudflare`
* __Client Secret:__ `cloudflare_client_secret`

*__Important Note:__ Cloudflare does not properly URL encode the secret. This means you'll either have to use
only alphanumeric characters for the secret or URL encode it yourself.*

## Configuration

### Application

*__Important Note:__ It is a requirement that the Authelia URL's can be requested by Cloudflare's servers. This usually
means that the URL's are accessible to foreign clients on the internet. There may be a way to configure this without
accessibility to foreign clients on the internet on Cloudflare's end but this is beyond the scope of this document.*

To configure [Cloudflare Zero Trust] to utilize Authelia as an [OpenID Connect] Provider:

1. Visit the [Cloudflare Zero Trust Dashboard](https://dash.teams.cloudflare.com)
2. Visit `Settings`
3. Visit `Authentication`
4. Under `Login nethods` select `Add new`
5. Select `OpenID Connect`
6. Set the following values:
   1. Name: `Authelia`
   2. App ID: `cloudflare`
   3. Client Secret: `cloudflare_client_secret`
   4. Auth URL: `https://auth.example.com/api/oidc/authorization`
   5. Token URL: `https://auth.example.com/api/oidc/token`
   6. Certificate URL: `https://auth.example.com/jwks.json`
   7. Enable `Proof Key for Code Exchange (PKCE)`
   8. Add the following OIDC Claims: `preferred_username`, `mail`
7. Click Save

### Authelia

The following YAML configuration is an example __Authelia__
[client configuration](../../../configuration/identity-providers/open-id-connect.md#clients) for use with [Cloudflare]
which will operate with the above example:

```yaml
- id: cloudflare
  description: Cloudflare ZeroTrust
  secret: cloudflare_client_secret
  public: false
  authorization_policy: two_factor
  redirect_uris:
    - https://example-team.cloudflareaccess.com/cdn-cgi/access/callback
  scopes:
    - openid
    - profile
    - email
  userinfo_signing_algorithm: none
```

## See Also

* [Cloudflare Zero Trust Generic OIDC Documentation](https://developers.cloudflare.com/cloudflare-one/identity/idp-integration/generic-oidc/)

[Authelia]: https://www.authelia.com
[Cloudflare]: https://www.cloudflare.com/
[Cloudflare Zero Trust]: https://www.cloudflare.com/products/zero-trust/
[OpenID Connect]: ../../openid-connect/introduction.md
