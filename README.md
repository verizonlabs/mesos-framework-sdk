## Mesos Framework SDK ##
This library aims to be a general purpose Golang library for writing
Mesos frameworks.

### Getting started ###
Mesos frameworks, at a minimum, require a scheduler to tell the framework
how to run tasks on the cluster when offers come in.

If you are unfamiliar with the Mesos architecture you can find more information [here](http://mesos.apache.org/documentation/latest/architecture/).

This SDK specifically targets the newer V1 streaming API of Mesos.
V0 support is not provided. The API can be broken down into a lower-level API as well as a higher-level API which provides more abstractions at the cost of being more opinionated.

### Low-Level API ###
The Low-level API gives maximum flexibility to the developer.

In this case, the low-level API is simply the protobufs that are generated
by protoc, along with the defined interfaces per component.

We use the standard protoc implementation and pull directly from the official Mesos codebase to keep versioning simple.
This allows us to pull protobufs for a specific version of Mesos and run protoc to generate the appropriate bindings.

### High-Level API ###
The high-level API takes an opinionated stance on the architecture of a mesos framework.

- Scheduler: This component implements all calls a scheduler would make towards the master.
- Client: A generic HTTP client for handling HTTP connections and calls.
- Task Manager: Handles how tasks are stored and delegated to other components.
- Resource Manager: Handles how offers and corresponding resources are delegated to other components.
- Storage: The storage tier handles how we interface with any storage backend.
- Logging: A simple logging interface to follow the MLOG format.
- HA (High Availability): An etcd-based leader election module.
- Server: A module for creating servers on components. This for example, can be used for an API, or to serve up a custom executor binary.
- Task: This defines how we define a "Task" to run on a cluster via JSON.
- Structures: Data structures that are useful for Mesos Frameworks.
- Utils: Various utility functions that are useful when writing a framework.

### Creating a Basic Framework ###
A basic framework will handle, at minimum:
- Taking offers from the Mesos master and declining or accepting them as necessary.
- Handing out Task's to be run on the resources offered by the Mesos master, if any.
- Ensuring state is synced with the master via reconciliation.

A basic framework can then be minimally made with a scheduler, event controller, event handler, task manager, and resource manager.

These four components will be enough to handle the previous mentioned responsibilities of a basic Mesos framework.

Using the default implementation of the high level SDK will make this easy.

The only custom logic one needs to write is the event controller.  This is a sort of "router" that routes events coming from the Mesos master and passes them off to your event handlers.

### Building ###

The SDK is link-only and not built on its own.

#### Testing ####

Tests can be run with `make test`. Similarly, you can run benchmarks with `make bench`.

### [License](LICENSE) ###
