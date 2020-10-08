# awscli-console-plugin
AWSCLI plugin to access AWS Console using your IAM or STS credentials

# Demo

# Motivation
The following library is distributed as `awscli` plugin, but could be used a standalone tool to access AWS Console
using IAM access & secret keys or STS temporary credentials.
The code is based on the following article https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_enable-console-custom-url.html

Often when working with multiple aws accounts there are several profiles defined in the `~/.aws/config` file. 
I find it tedious to follow the separate process login into individual accounts for Console access (SSO certainly helps). 
So I came up with this extension to get instant access to a given AWS Console while working in the terminal.

# Installation & Usage
To install it as a plugin for `awslic` follow these steps

Install `awscli-console-plugin` using pip
```bash
$ pip install .
```

Modify the `plugin` sections in the `~/.aws/config` file
```bash
[plugins]
console = console
[profile auth]
region=ap-southeast-2
aws_access_key_id=AKIAXXXXXXXXXXX
aws_secret_access_key=XYXYXYYXYXYXYXYXY/ZZZZZZ
[profile dev]
role_arn=arn:aws:iam::1234567890123:role/OrganizationAccountAccessRole
region=ap-southeast-2
source_profile=auth
```

Verify that the plugin is operational
```bash
$ aws console help
```
Usage
```bash
aws console --profile dev
```