# Micro [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Doc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/go.unistack.org/micro/v3?tab=overview) [![Status](https://github.com/unistack-org/micro/workflows/build/badge.svg?branch=master)](https://github.com/unistack-org/micro/actions?query=workflow%3Abuild+branch%3Amaster+event%3Apush) [![Lint](https://goreportcard.com/badge/go.unistack.org/micro/v3)](https://goreportcard.com/report/go.unistack.org/micro/v3)

Micro is a standard library for microservices.

## Overview

Micro provides the core requirements for distributed systems development including RPC and Event driven communication. 

## Features

Micro abstracts away the details of distributed systems. Here are the main features.

- **Authentication** - Auth is built in as a first class citizen. Authentication and authorization enable secure 
zero trust networking by providing every service an identity and certificates. This additionally includes rule 
based access control.

- **Dynamic Config** - Load and hot reload dynamic config from anywhere. The config interface provides a way to load application 
level config from any source such as env vars, file, etcd. You can merge the sources and even define fallbacks.

- **Data Storage** - A simple data store interface to read, write and delete records. It includes support for memory, file and 
CockroachDB by default. State and persistence becomes a core requirement beyond prototyping and Micro looks to build that into the framework.

- **Service Discovery** - Automatic service registration and name resolution. Service discovery is at the core of micro service 
development. When service A needs to speak to service B it needs the location of that service.

- **Load Balancing** - Client side load balancing built on service discovery. Once we have the addresses of any number of instances 
of a service we now need a way to decide which node to route to. We use random hashed load balancing to provide even distribution 
across the services and retry a different node if there's a problem. 

- **Message Encoding** - Dynamic message encoding based on content-type. The client and server will use codecs along with content-type 
to seamlessly encode and decode Go types for you. Any variety of messages could be encoded and sent from different clients. The client 
and server handle this by default.

- **Transport** - gRPC or http based request/response with support for bidirectional streaming. We provide an abstraction for synchronous communication. A request made to a service will be automatically resolved, load balanced, dialled and streamed.

- **Async Messaging** - PubSub is built in as a first class citizen for asynchronous communication and event driven architectures.
Event notifications are a core pattern in micro service development.

- **Synchronization** - Distributed systems are often built in an eventually consistent manner. Support for distributed locking and 
leadership are built in as a Sync interface. When using an eventually consistent database or scheduling use the Sync interface.

- **Pluggable Interfaces** - Micro makes use of Go interfaces for each system abstraction. Because of this these interfaces 
are pluggable and allows Micro to be runtime agnostic.

## Getting Started

To be created.

## License

Micro is Apache 2.0 licensed.

