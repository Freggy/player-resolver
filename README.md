Luxor Player Resolver
=====================

This is a microservice which resolves Minecraft player UUIDs to their corresponding player name and vice versa.
The received data will then be stored in Cassandra to cache results.

Restful API
-----------

The player resolver exposes two routes:

* `/uuid/:name` - Resolves names to `UUID`s
* `/name/:uuid` - Resolves `UUID`s to names

### GET `/uuid/:name`

When performing a GET request with the given valid name, you will then receive the UUID and the name of the player.

```JSON
{"id": "{UUID}", "name": "{name}"}
```

If the name is not valid an error response is sent back.
```JSON
{"status": 400, "message": "Provided name is not valid", "type": "InvalidNameException"}
```

### GET `/name/:uuid`

When performing a GET request with the given `UUID`, you will then receive the players current name and their `UUID`.

```JSON
{"id": "{UUID}", "name": "{name}"}
```

If the provided `UUID` is not valid an error response is sent back.

```JSON
{"status": 400, "message": "Provided UUID is not vaild", "type": "InvalidUUIDException"}
```