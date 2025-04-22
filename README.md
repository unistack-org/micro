# Micro
![Coverage](https://img.shields.io/badge/Coverage-44.1%25-yellow)
[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Doc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/go.unistack.org/micro/v3?tab=overview)
[![Status](https://git.unistack.org/unistack-org/micro/actions/workflows/job_tests.yml/badge.svg?branch=v3)](https://git.unistack.org/unistack-org/micro/actions?query=workflow%3Abuild+branch%3Av3+event%3Apush)
[![Lint](https://goreportcard.com/badge/go.unistack.org/micro/v3)](https://goreportcard.com/report/go.unistack.org/micro/v3)

Micro is a standard library for microservices.

## Overview

Micro provides the core requirements for distributed systems development including SYNC and ASYNC communication. 

## Features

Micro abstracts away the details of distributed systems. Main features:

- **Dynamic Config** - Load and hot reload dynamic config from anywhere. The config interface provides a way to load application 
level config from any source such as env vars, cmdline, file, consul, vault, etc... You can merge the sources and even define fallbacks.

- **Data Storage** - A simple data store interface to read, write and delete records. It includes support for memory, file and 
s3. State and persistence becomes a core requirement beyond prototyping and Micro looks to build that into the framework.

- **Service Discovery** - Automatic service registration and name resolution. Service discovery is at the core of micro service 
development.

- **Message Encoding** - Dynamic message encoding based on content-type. The client and server will use codecs along with content-type 
to seamlessly encode and decode Go types for you. Any variety of messages could be encoded and sent from different clients. The client 
and server handle this by default.

- **Async Messaging** - Pub/Sub is built in as a first class citizen for asynchronous communication and event driven architectures.
Event notifications are a core pattern in micro service development.

- **Synchronization** - Distributed systems are often built in an eventually consistent manner. Support for distributed locking and 
leadership are built in as a Sync interface. When using an eventually consistent database or scheduling use the Sync interface.

- **Pluggable Interfaces** - Micro makes use of Go interfaces for each system abstraction. Because of this these interfaces 
are pluggable and allows Micro to be runtime agnostic.

## License

Micro is Apache 2.0 licensed.

