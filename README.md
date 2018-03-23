# aws-rekey

Rotates your AWS access keys in seconds

CLI app that will generate a new IAM access key for an AWS account, add it
to the `~/.aws/config` file (overwriting the previous entry), then deleting
the old access key.  This will happen only if the key is older than 30
days.  The application also modifies your `.aws/config` file with your
new credentials so you have literally nothing to do but run this binary.

```
$ rotate-aws-key [--key N] [--profile NAME]
```

It should take less than 10 seconds to complete, typically 2-4 seconds.

