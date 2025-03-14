# Interfaces

This document will go through the following interfaces:
1. Event Message (EN) - interface
2. Notification Message (NM) - interface

### event-producer microservice

Interfaces implemented:

| interface name | full name                           | producer           | consumer       |
|----------------|-------------------------------------|--------------------|----------------|
| EP.EM          | event-producer.event-message        | application - http | event-producer |
| EP.NM          | event-producer.notification-message | event-producer     | event-consumer |

### event-consumer microservice

Interfaces implemented:

| interface name | full name                           | producer           | consumer       |
|----------------|-------------------------------------|--------------------|----------------|
| EC.NM          | event-consumer.notification-message | event-producer     | event-consumer |

## 1. Event Message (EM) - interface

The purpose of this interface is to create an event message and then use this event message to create a notification message.
An event message is created using the /event endpoint in the event-producer microservice.

```text
curl -i -X POST http://localhost:8080/event -H "Accept: application/json" --data '{
"description":"hello world"
}'
```

Event Message json **payload**:

```json
{
  "description":"hello world"
}
```

Event Message **schema**:

```json
{
  "title": "Event Message",
  "description": "Event Message schema",
  "type": "object",
  "properties": {
    "description": {
      "description": "some information",
      "type": "string"
    }
  }
}
```

## 2. Notification Message (NM) - interface

The purpose of this interface is to **produce events** (EP.NM) to the redis stream and **consume events** (EC.NM) from the redis stream.
Both the event-producer and event-consumer microservices will be using this interface.

The event-producer microservice creates a Notification Message that will be sent to the redis stream to be consumed by the event-consumer.
A Notification Message is only sent to the redis stream after an Event Message is created in the event-producer which then builds the Notification Message
using fields from the Event Message.

Notification Message **json**:

```json
{
  "id": 123,
  "description": "hello world",
  "timestamp": "2025-02-25T22:06:25.478310425Z"
}
```

Notification Message **schema**:

```json
{
  "title": "Notification Message",
  "description": "Notification Message schema",
  "type": "object",
  "properties": {
    "id": {
      "description": "random generated id",
      "type": "integer"
    },
    "description": {
      "description": "some information from the event message",
      "type": "string"
    },
    "timestamp": {
      "description": "creation time of notification message",
      "type": "date-time"
    }
  }
}
```
