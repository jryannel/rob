# Remote Object Labs

This is a concept for a remote object system based on NATS. It is a work in progress.

## Concept

A object is identified by an ID.
The ID is a string that is unique within the scope of the NATS server.
The ID is used to publish and subscribe to messages about the object.

Value messages are published to the object ID + '.' + value name
Method messages are published to the object ID + '.' + method name
Signal messages are published to the object ID + '.' + signal name

Question: shall values, methods and signals be published to the subjects like "<id>.<kind>.<name>"

## Client Features

- Values can be read on the client side and be notified about changes.
- Methods can be invoked on the client side.
- Signal can be notified on the client side.

A client can not register methods, provide values or publish signals.

## Service Features

- Values can be published by the service. The service can be notified about changes.
- There is an initial value for each value using a value provider.
- Methods can be registered by the service. The service can be notified about invocations.
- Signals can be published by the service.

A service is notified about changes to values, invocations of methods and signals.
