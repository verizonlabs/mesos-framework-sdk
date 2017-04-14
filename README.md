## Mesos Framework SDK ##
This library aims to be a general purpose Golang library for writing
mesos frameworks.

### Getting started ###
Mesos frameworks at a minimum require a scheduler to tell the framework
how to run tasks on the cluster when offers come in.

If you are unfamilar with the Mesos architecture, please read here first:
http://mesos.apache.org/documentation/latest/architecture/

The SDK has two levels of API's, a low-level API and a higher-level API.

### Low-Level API ###
The Low-level API gives maximum flexibility to the developer.

In this case, the low-level API is simply the protobufs that are generated
by protoc, along with the defined interfaces per component.

We use the standard protoc implementation and pull directly from the mesos git repos to keep versioning simple.

In this way, we can pull a mesos tag from the apache git repository and run protoc to create the appropriate protobufs for a particular version of mesos.

### High-Level API ###
The high-level api takes an opinionated stance on the architecture of a mesos framework.

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
- Taking offers from the mesos master and declining or accepting them as necessary.
- Handing out Task's to be run on the resources offered by the mesos-master, if any.
- Ensuring state is synced with the master via Reconciliation.

A basic framework can then be minimally made with a scheduler, a scheduler event controller, a task manager and a resource manager.

These four components will be enough to handle the previous mentioned responsiblities of a basic mesos framework.

Using the default implementations from the high level SDK will make this easy.

The only custom logic one needs to write is the event controller.  This is a sort of "router" that routes events coming from the mesos master and cluster.

When a certain event  occurs, you define how you react to those events.  Please see the examples section to see an implementation of this.

The interface is within the scheduler/events package.
