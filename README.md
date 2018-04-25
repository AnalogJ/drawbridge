<p align="center">
  <a href="https://github.com/AnalogJ/drawbridge">
  <img width="300" alt="drawbridge_view" src="https://rawgit.com/AnalogJ/drawbridge/master/logo.svg">
  </a>
</p>

# Drawbridge

[![Circle CI](https://img.shields.io/circleci/project/github/AnalogJ/drawbridge.svg?style=flat-square)](https://circleci.com/gh/AnalogJ/drawbridge)
[![Coverage Status](https://img.shields.io/codecov/c/github/AnalogJ/drawbridge.svg?style=flat-square)](https://codecov.io/gh/AnalogJ/drawbridge)
[![GitHub license](https://img.shields.io/github/license/AnalogJ/drawbridge.svg?style=flat-square)](https://github.com/AnalogJ/drawbridge/blob/master/LICENSE)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/analogj/drawbridge)
[![Go Report Card](https://goreportcard.com/badge/github.com/AnalogJ/drawbridge?style=flat-square)](https://goreportcard.com/report/github.com/AnalogJ/drawbridge)
[![GitHub release](http://img.shields.io/github/release/AnalogJ/drawbridge.svg?style=flat-square)](https://github.com/AnalogJ/drawbridge/releases)
[![Github All Releases](https://img.shields.io/github/downloads/analogj/drawbridge/total.svg?style=flat-square)](https://github.com/AnalogJ/drawbridge/releases)

Bastion/Jumphost tunneling made easy

# Introduction
> A Jump/Bastion host is a special purpose computer on a network specifically designed and configured to withstand attacks.
> The computer generally hosts a single application, for example a proxy server, and all other services are removed or limited to reduce the threat to the computer.
> It is hardened in this manner primarily due to its location and purpose, which is either on the outside of a firewall or in a demilitarized zone (DMZ) and usually involves access from untrusted networks or computers.
> - [Bastion Host - Wikipedia](https://en.wikipedia.org/wiki/Bastion_host)

In secure cloud architectures, jump/bastion hosts are the primary method to access the internal/protected network.
This means that all traffic can be audited, and that a single server can be shut down in the event that the network is compromised.


However as this architecture is scaled up and deployed across multiple environments (testing, staging, production), it can
be complicated to maintain a single `~/.ssh/config` file that allows you to tunnel into your various jump host protected internal networks.

Drawbridge aims to solve this problem in a flexible and scalable way.


# Features

- Single binary (available for macOS and linux), only depends on `ssh`, `ssh-agent` and `scp`
- Uses customizable templates to ensure that Drawbridge can be used by any organization, in any configuraton
- Helps organize your SSH config files and PEM files
- Generates SSH Config files for your servers spread across multiple environments and stacks.
	- multiple ssh users/keypairs
	- multiple environments
	- multiple stacks per environment
	- etc..
- Can be used to SSH directly into an internal node, routing though bastion, leveraging SSH-Agent
- Able to download files from internal hosts (through the jump/bastion host) using SCP syntax
- Supports HTTP proxy to access internal stack urls.
- Lists all managed config files in a heirarchy that makes sense to your organization
- Custom templated files can be automatically generated when a new SSH config is created.
	- eg. Chef knife.rb configs, Pac/Proxy files, etc.
- Cleanup utility is built-in
- `drawbridge update` lets you update the binary inplace.
- Pretty colors. The CLI is all colorized to make it easy to skim for errors/warnings

# Usage

```
$ drawbridge help
 ____  ____    __    _    _  ____  ____  ____  ____    ___  ____
(  _ \(  _ \  /__\  ( \/\/ )(  _ \(  _ \(_  _)(  _ \  / __)( ___)
 )(_) ))   / /(__)\  )    (  ) _ < )   / _)(_  )(_) )( (_-. )__)
(____/(_)\_)(__)(__)(__/\__)(____/(_)\_)(____)(____/  \___/(____)
github.com/AnalogJ/drawbridge                  darwin-amd64-1.0.7

NAME:
   drawbridge - Bastion/Jumphost tunneling made easy

USAGE:
   drawbridge [global options] command [command options] [arguments...]

VERSION:
   1.0.7

AUTHOR:
   Jason Kulatunga <jason@thesparktree.com>

COMMANDS:
     create         Create a drawbridge managed ssh config & associated files
     list           List all drawbridge managed ssh configs
     connect        Connect to a drawbridge managed ssh config
     download, scp  Download a file from an internal server using drawbridge managed ssh config, syntax is similar to scp command.
     delete         Delete drawbridge managed ssh config(s)
     update         Update drawbridge to the latest version
     help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

```

# Actions

## Create

Using the `questions` & `config_template` defined in the configuration file (`~/drawbridge.yaml`) Drawbridge will attempt to
generate a managed ssh config file. Drawbrige will prompt the user for any questions which it is unable to determine an
answer (no default value and no flag value specified).

Questions & Templates can be customized completely to match your organization.

```
$ drawbridge create --environment prod --shard us-west-2

Current Answers:
environment: prod
shard: us-west-2
stack_name: app
Please enter a value for `shard_type` [string] - Is this a live (green) or idle (blue) stack?:
idle
Please enter a value for `username` [string] - What username do you use to login to this stack?:
aws
WARNING: PEM file missing. Place it at the following location before attempting to connect. /Users/jason/.ssh/drawbridge/pem/prod/aws-prod.pem
Writing template to /Users/jason/.ssh/drawbridge/prod-app-idle-us-west-2

```

You can also enable `DRYRUN` mode to see exactly what files Drawbrige would generate, without actually writing any files.

```
$ drawbridge create --environment prod --dryrun
...
2018/04/22 23:56:23 Writing template to /Users/jason/.ssh/drawbridge/prod-app-idle-us-west-1
[DRYRUN] Would have written content to /Users/jason/.ssh/drawbridge/prod-app-idle-us-west-1:

# This file was automatically generated by Drawbridge
# Do not modify.
#
...
```


## Connect
```
$ drawbridge connect
Rendered Drawbridge Configs:
├── [prod]  environment
│   └── [app]  stack_name
│       ├── [us-east-1]  shard
│       │   ├── [1]  shard_type: idle, username: aws
│       │   └── [2]  shard_type: live, username: aws
│       └── [us-east-2]  shard
│           ├── [3]  shard_type: idle, username: aws
│           └── [4]  shard_type: live, username: aws
├── [stage]  environment
│   └── [app]  stack_name
│       └── [us-east-2]  shard
│           ├── [5]  shard_type: idle, username: aws
│           └── [6]  shard_type: live, username: aws
└── [test]  environment
    └── [app]  stack_name
        ├── [us-east-1]  shard
        │   ├── [7]  shard_type: idle, username: aws
        │   └── [8]  shard_type: live, username: aws
        └── [us-east-2]  shard
            ├── [9]  shard_type: idle, username: aws
            └── [10]  shard_type: live, username: aws

Enter number of drawbridge config you would like to connect to (1-10):
```

`drawbridge connect` will connect you to the bastion/jump host using a specified Drawbridge config file. It'll also add
the associated PEM key to your `ssh-agent`.

If you want to connect directly to a internal server, you can do so by specifying the hostname/short name using the `--dest` flag

`drawbridge connect --dest db-1`

## Delete

```
$ drawbridge delete
...
        └── [us-east-2]  shard
            ├── [9]  shard_type: idle, username: aws
            └── [10]  shard_type: live, username: aws

Enter number of drawbridge config you would like to delete:
10
Are you sure you would like to delete this config and associated templates? (PEM files will not be deleted)

environment: test
shard: us-east-2
shard_type: live
stack_name: app
username: aws

Please confirm [true/false]:
true
Deleting config file: /Users/jason/.ssh/drawbridge/test-app-live-us-east-2
Deleting answers file
Finished

```

You can use the `--force` flag to disable the confirm prompt. The `--all` flag can be used to delete all Drawbridge managed
configs in one command.

You can use the following command to completely wipe out all Drawbridge files and start over.

`drawbridge delete --all --force`


## Update

```
$ drawbridge update

Current: v1.0.7 [2018-04-22]. Available: v1.0.7 [2018-04-23]
Release notes are available here: https://github.com/AnalogJ/drawbridge/releases/tag/v1.0.7

```

## Download



# Configuration
We support a global YAML configuration file that must be located at `~/drawbridge.yaml`

Check the [example.drawbridge.yml](https://github.com/AnalogJ/drawbridge/blob/master/example.drawbridge.yaml) file for a fully commented version.

# Testing [![Circle CI](https://img.shields.io/circleci/project/github/AnalogJ/drawbridge.svg?style=flat-square)](https://circleci.com/gh/AnalogJ/drawbridge)
Drawbridge provides an extensive test-suite based on `go test`.
You can run all the integration & unit tests with `go test $(glide novendor)`

CircleCI is used for continuous integration testing: https://circleci.com/gh/AnalogJ/drawbridge

# Contributing
If you'd like to help improve Drawbridge, clone the project with Git and install dependencies by running:

```
$ git clone git://github.com/AnalogJ/drawbridge
$ glide install
```

Work your magic and then submit a pull request. We love pull requests!

If you find the documentation lacking, help us out and update this README.md.
If you don't have the time to work on Drawbridge, but found something we should know about, please submit an issue.


# To-Do List

We're actively looking for pull requests in the following areas:


- the ability to open the ssh tunnel, with http port binding locally (connect)
	- local ports chosen will be dynamic and depend on the hash of the config filepath (unique on the config level) https://stats.stackexchange.com/questions/26344/how-to-uniformly-project-a-hash-to-a-fixed-number-of-buckets
	- the ability to create/update a pac file, which points to a proxy server inside behind the bastion (--pac)


# Versioning
We use SemVer for versioning. For the versions available, see the tags on this repository.

# Authors
Jason Kulatunga - Initial Development - @AnalogJ

# License

- MIT
- [Logo: Castle by Jemis mali from the Noun Project](https://thenounproject.com/search/?q=castle&i=1063814)


# References

- https://github.com/moul/awesome-ssh/blob/master/README.md
- https://github.com/dbrady/ssh-config
- https://github.com/k4m4/terminals-are-sexy
- https://github.com/n1trux/awesome-sysadmin
- https://github.com/cjbarber/ToolsOfTheTrade
- https://github.com/dastergon/awesome-sre
- https://stackoverflow.com/questions/17355667/replace-current-process
- https://stats.stackexchange.com/questions/26344/how-to-uniformly-project-a-hash-to-a-fixed-number-of-buckets
- moul/advanced-ssh-config
- https://github.com/emre/storm
- https://stackoverflow.com/questions/12484398/global-template-data
- https://stackoverflow.com/questions/35612456/how-to-use-golang-template-missingkey-option
- https://medium.com/@dgryski/consistent-hashing-algorithmic-tradeoffs-ef6b8e2fcae8
- https://github.com/mitchellh/go-homedir
- https://gobyexample.com/execing-processes
- https://groob.io/posts/golang-execve/
- https://www.digitalocean.com/community/tutorials/how-to-configure-custom-connection-options-for-your-ssh-client
- https://github.com/zalando/go-keyring
- https://github.com/tmc/keyring
- https://github.com/jaraco/keyring
- https://github.com/99designs/keyring
- https://unix.stackexchange.com/questions/64795/must-i-store-a-private-key-in-a-file
