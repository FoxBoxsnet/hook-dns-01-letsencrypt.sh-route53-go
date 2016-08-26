# Amazon Route 53 hook for `letsencrypt.sh`

This a hook for [letsencrypt.sh](https://github.com/lukas2511/letsencrypt.sh) (a [Let's Encrypt](https://letsencrypt.org/) ACME client) that allows you to use [Amazon Route 53](https://aws.amazon.com/jp/route53/) DNS records to respond to `dns-01` challenges. Requires Go and your Amazon Route 53 account `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` being in the environment.

## Installation

```
$ git clone https://github.com/lukas2511/letsencrypt.sh
$ cd letsencrypt.sh
go get github.com/FoxBoxsnet/letsencrypt.sh-dns-route53
```

## Configuration

Your account's `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are expected to be in the environment, so make sure to:

```
$ export AWS_ACCESS_KEY_ID=ACCESS_KEYXXXXXXXXXX
$ export AWS_SECRET_ACCESS_KEY=SECRET_KEYXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

## Usage

```
$ ./usr/local/bin/letsencrypt.sh \
    --cron \
    --force \
    --domain www.example.com \
    --challenge dns-01 \
    --hook hook/letsencrypt.sh-dns-route53

# !! WARNING !! No main config file found, using default config!
+ Generating account key...
+ Registering account key with letsencrypt...
Processing www.example.com
 + Signing domains...
 + Creating new directory /app/lets/certs/www.example.com ...
 + Generating private key...
 + Generating signing request...
 + Requesting challenge for www.example.com...
[  info ] 2016/08/24 13:50:30 INFO: deploy_challenge common_name: [www.example.com]
[  info ] 2016/08/24 13:50:31 INFO: Found ZONE : /hostedzone/Z32DGRNCTFM483 Domain : example.com.
[  info ] 2016/08/24 13:50:31 INFO: apply wait please wait for a while...
[  info ] 2016/08/24 13:51:04 INFO; Record status : INSYNC
 + Responding to challenge for www.example.com...
[  info ] 2016/08/24 13:51:06 INFO: clean_challenge common_name: [www.example.com]
[  info ] 2016/08/24 13:51:06 INFO: Found ZONE : /hostedzone/Z32DGRNCTFM483 Domain : example.com.
[  info ] 2016/08/24 13:51:07 INFO: apply wait please wait for a while...
[  info ] 2016/08/24 13:51:39 INFO; Record status : INSYNC
 + Challenge is valid!
 + Requesting certificate...
 + Checking certificate...
 + Done!
 + Creating fullchain.pem...
[  info ] 2016/08/24 13:51:40 INFO: deploy_cert common_name: [www.example.com]
[  info ] 2016/08/24 13:51:40 INFO: Private key       : /app/lets/certs/www.example.com/privkey.pem
[  info ] 2016/08/24 13:51:40 INFO: Private cert      : /app/lets/certs/www.example.com/cert.pem
[  info ] 2016/08/24 13:51:40 INFO: Private fullchain : /app/lets/certs/www.example.com/fullchain.pem
[  info ] 2016/08/24 13:51:40 INFO: Private chain     : /app/lets/certs/www.example.com/chain.pem
```
# Use SDK & package
| Package | License |
|:--------|:--------|
|[aws/aws-sdk-go](https://github.com/aws/aws-sdk-go) | [![Hex.pm](https://img.shields.io/hexpm/l/plug.svg?maxAge=2592000)](https://github.com/aws/aws-sdk-go/blob/master/LICENSE.txt) |
| [comail/colog](https://github.com/comail/colog) | [![Packagist](https://img.shields.io/packagist/l/doctrine/orm.svg?maxAge=2592000)](https://github.com/comail/colog/blob/master/LICENSE) |

# Special thanks
The [kappataumu/letsencrypt-cloudflare-hook](https://github.com/kappataumu/letsencrypt-cloudflare-hook) was written as a reference.
