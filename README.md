# r53.go

This script, written in golang acts as its own pseudo-ddns provider. It uses the AWS v2 golang sdk to update a DNS record in Route 53 to point to a given IP address.

## Pre-requisites:

You must have an IAM account with api credentials and permissions to update the Route53 records in question, ie.:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowUpdateDns",
            "Effect": "Allow",
            "Action": "route53:ChangeResourceRecordSets",
            "Resource": "*"
        }
    ]
}
```

Place your IAM credentials in `~/.aws/credentials` on the service which is running this script. eg.
```
[default]
aws_access_key_id=AK****************YU
aws_secret_access_key=pQ***************************************tq
```

## Building & running the script on your own server

The compiled executable for linux [r53-ddns](https://github.com/jackphilippi/r53-ddns/blob/main/r53-ddns) is already included in this repo and is ready to use, so if you prefer to use that, ignore the first 3 steps:

1. Ensure you have golang installed @ minimum version `1.19`
2. `go mod tidy`
3. `go build`
4. Copy the compiled `r53-ddns` file to your server/router where you want it to run (in my case, it's a GL.inet Flint router with OpenWRT)
    a. You can use scp for this, eg: `scp r53-ddns root@192.168.0.1:/root/r53-ddns`
5. SSH into the server and run `./r53-ddns` (see below for args)

## Usage

To use this script, you need to provide the following arguments:

- `zone_id`: The ID of the hosted zone in Route 53 where you want to update the DNS record.
- `hostname`: The hostname of the DNS record you want to update.
- `ip_address`: The IP address you want to set for the DNS record.

Here's an example of how to use the script:

`./r53 <zone_id> <hostname> <ip_address> | tee -a /root/ddns-output.log`

To set this up to be executed regularly, you can set a cron job with `crontab -e`:

```
0 */6 * * * /root/r53 Z000000000000000000MW home.foo.bar $(/root/get-ip.sh) | tee -a /root/ddns-output.log
   ^             ^                                                 ^              
   |             |                                                 |
   |             this is the path to the r53 executable script     |
   |                                                               this is a local script which gets the current public ip address of the server (see below)
   this cron job will run every 6 hours
```