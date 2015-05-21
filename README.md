# AWS service DNS

[![Build Status](https://travis-ci.org/kreuzwerker/awssd.svg)](https://travis-ci.org/kreuzwerker/awssd)

AWS service DNS (`awssd`) is a utility that creates Route53 A records for EC2 instances following a tag convention. This can be used to support reverse proxy / load balancing strategies other then ELB for autoscaling groups or other elastic groups of EC2 instances.

## Usage

`awssd` should be run periodically to gather *running* EC2 instances and group them by the value for a given tag key. It will then proceed to create or update ("upsert") `A` record sets in a given zone with the public (or private) IP addresses of those instances, using the group-by value of the tag as name.

`awssd` will never delete records. For security reasons `awssd` should nevertheless be used in a dedicated hosted zone (e.g. `api.example.com`). If the service records should be accessible without specifying the name servers in this zone the usual delegation via `NS` can be setup.

When running `awssd` you need to specify:

* the Route53 domain without a trailing `.` (`-d` flag) to determine the fully-qualified record set values and finding the hosted zone id
* the tag key to group instances by (`-g` flag)

Optionally you can:

* perform a dry-run (`-dry` flag) that does not touch record sets
* pass [EC2 filters](see http://docs.aws.amazon.com/AWSEC2/latest/CommandLineReference/ApiReference-cmd-DescribeInstances.html) (`-f` flag) in key=value notation (or multiple, comma-separated filters in the same notation)
* prefer private IP addresses over public IP addresses (`-p` flag): when no public IP addresses are available for an instance, a value of `false` has no effect (defaults to `true`)
* specify the EC2 region (`-r` flag), defaults to `eu-west-1`
* specify the record set TTL in seconds (`-t` flag), default to 60
* enforce verbosity (`-v` flag)

Credential detection follow the aws-sdk standard (checking environment variables, profile configuration and instance profiles).

### Example

Let an account contain 6 instances:

1. `i-1`, public IP 1.2.3.4, tags `service=foo`, `environment=staging`
* `i-2`, public IP 2.2.3.4, tags `service=bar`, `environment=staging`
* `i-3`, public IP 3.2.3.4, tags `service=foo`, `environment=staging`
* `i-4`, public IP 4.2.3.4, tags `service=bar`, `environment=staging`
* `i-5`, public IP 5.2.3.4, tags `service=foo`, `environment=production`
* `i-6`, public IP 6.2.3.4, tags `service=bar`, `environment=production`

Calling `awssd` like this

* `awssd -d api.example.com -f tag:environment=staging -p=true -g service -v` (when running the binary build) or
* `docker run -it kreuzwerker/awssd -d api.example.com -f tag:environment=staging -p=true -g service -v` (when running it via Docker)

will

1. group records by the value of the `service` tag
* filter instances with the key- / value-pair `environment=staging`
* register `A` records for the *public* IP addresses of those instances
* in the Route53 zone `api.example.com`
* be verbose about it

Specifically, this will upsert the following Route 53 records (in the zone with the name `api.example.com`:

```
foo.api.example.com.		60	IN	A	1.2.3.4
foo.api.example.com.		60	IN	A	3.2.3.4
```

and

```
bar.api.example.com.		60	IN	A	2.2.3.4
bar.api.example.com.		60	IN	A	4.2.3.4
```

### IAM Policy

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:DescribeInstances"
        "route53:ListHostedZones"
      ],
      "Effect": "Allow",
      "Resource": [
        "*"
      ]
    },
    {
      "Action": [
        "route53:ChangeResourceRecordSets",
        "route53:ListResourceRecordSets"
      ],
      "Effect": "Allow",
      "Resource": [
        "arn:aws:route53:::hostedzone/hosted-zone-id-for-your-service-zone"
      ]
    }
  ]
}
```

## Future roadmap

* [ ] Support deletion of records as soon as we have metadata / tag support in R53
* [ ] Support retries / error handling as soon as the aws-sdk supports it
* [ ] Support pagination as soon as the aws-sdk supports it (**WARNING**: currently the tools **DOES NOT** paginate)
* [ ] Support the creation of individual instance names, maybe with wordlists from sources such as Dynamo + appropriate reverse DNS entries
