# Amazon Route 53 hook for `letsencrypt.sh`

This a hook for [letsencrypt.sh](https://github.com/lukas2511/letsencrypt.sh) (a [Let's Encrypt](https://letsencrypt.org/) ACME client) that allows you to use [Amazon Route 53](https://aws.amazon.com/jp/route53/) DNS records to respond to `dns-01` challenges. Requires Go and your Amazon Route 53 account AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY being in the environment.

## Installation

```
$ git clone https://github.com/lukas2511/letsencrypt.sh
$ cd letsencrypt.sh
$ mkdir hooks
$ go get github.com/FoxBoxsnet/letsencrypt.sh-dns-route53
$ go build -ldflags "-s" letsencrypt.sh-dns-route53
```

## Configuration

Your account's CloudFlare email and API key are expected to be in the environment, so make sure to:

```
$ export AWS_ACCESS_KEY_ID=ACCESS_KEYXXXXXXXXXX
$ export AWS_SECRET_ACCESS_KEY=SECRET_KEYXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

## Usage

```
$ ./letsencrypt.sh -c -d example.com -t dns-01 -k './route53'
# INFO: Using main config file /etc/nginx/certs/example.com/config
Processing example.com
 + Signing domains...
 + Creating new directory /etc/nginx/certs/example.com ...
 + Generating private key...
 + Generating signing request...
 + Requesting challenge for example.com...
[  info ] 2016/08/23 14:19:53 INFO: Found ZONE : /hostedzone/XXXXXXXXXXXXXX Domain : example.com.
[  info ] 2016/08/23 14:19:53 INFO: deploy_challenge
[  info ] 2016/08/23 14:19:53 INFO: apply wait please wait for a while.
[  info ] 2016/08/23 14:20:26 INFO; Record status : INSYNC
 + Responding to challenge for example.com...
[  info ] 2016/08/23 14:20:27 INFO: Found ZONE : /hostedzone/XXXXXXXXXXXXXX Domain : example.com.
[  info ] 2016/08/23 14:20:27 INFO: clean_challenge
[  info ] 2016/08/23 14:20:28 INFO: apply wait please wait for a while.
[  info ] 2016/08/23 14:21:00 INFO; Record status : INSYNC
 + Challenge is valid!
 + Requesting certificate...
 + Checking certificate...
 + Done!
 + Creating fullchain.pem...
[  info ] 2016/08/23 14:21:03 INFO: Found ZONE : /hostedzone/XXXXXXXXXXXXXX Domain : example.com.
 + Done!



#
# !! WARNING !! No main config file found, using default config!
#
Processing example.com
 + Signing domains...
 + Creating new directory /home/user/letsencrypt.sh/certs/example.com ...
 + Generating private key...
 + Generating signing request...
 + Requesting challenge for example.com...
 + CloudFlare hook executing: deploy_challenge
 + DNS not propagated, waiting 30s...
 + DNS not propagated, waiting 30s...
 + Responding to challenge for example.com...
 + CloudFlare hook executing: clean_challenge
 + Challenge is valid!
 + Requesting certificate...
 + Checking certificate...
 + Done!
 + Creating fullchain.pem...
 + CloudFlare hook executing: deploy_cert
 + ssl_certificate: /home/user/letsencrypt.sh/certs/example.com/fullchain.pem
 + ssl_certificate_key: /home/user/letsencrypt.sh/certs/example.com/privkey.pem
 + Done!
```