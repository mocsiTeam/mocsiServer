# Methods for working with rooms

## Create room

Creates a new Room in Mocsi with unique name, display name and password. Also returns information about the room.

### **Shemas**

```
type Room {
  id: ID!
  unique_name: String!
  name: String!
  link: String!
  owner: User!
  editors: [User!]
  users: [User!]
}

input NewRoom {
  unique_name: String!
  name: String!
  password: String!
}
```

### **Mutations request**

```
createRoom(
    input: NewRoom!
): Room!
```

### **Example**

Request

```
mutation {
    createRoom(
        input: {
            unique_name: "p1ck0_room",
            name: "p1ck0 room"
            password: ""
        }
    ) 
    {
        id
        unique_name
        name
        link
        owner {
            id
            nickname //the whole user structure can be queried
        }
        editors {
            id
            nickname //the whole user structure can be queried
        }
        users {
            id
            nickname //the whole user structure can be queried
        }
    }
}
```

Response


```
{
  "data": {
    "createRoom": {
      "id": "1",
      "unique_name": "p1ck0_room",
      "name": "p1ck0 room",
      "link": "domain/p1ck0_room",
      "owner": {
          "id": "1",
          "nickname": "p1ck0"
      },
      "editors": {
          "id": "1",
          "nickname": "p1ck0"
      },
      "users": {
          "id": "1",
          "nickname": "p1ck0"
      }
    }
  }
}
```

## Get My Rooms

You get information about the rooms to which you have access.

### **Shemas**

```
type Room {
  id: ID!
  unique_name: String!
  name: String!
  link: String!
  owner: User!
  editors: [User!]
  users: [User!]
}
```

### **Query request**

```
getMyRooms: [Room!]
```
### **Example**

Request

```
query {
    getMyRooms {
        id
        unique_name
        name
        link
        owner {
            id
            nickname //the whole user structure can be queried
        }
        editors {
            id
            nickname //the whole user structure can be queried
        }
        users {
            id
            nickname //the whole user structure can be queried
        }
    }
}
```

Response

```
{
  "data": {
    "getMyRooms": {
      "id": "1",
      "unique_name": "p1ck0_room",
      "name": "p1ck0 room",
      "link": "domain/p1ck0_room",
      "owner": {
          "id": "1",
          "nickname": "p1ck0"
      },
      "editors": {
          "id": "1",
          "nickname": "p1ck0"
      },
      "users": {
          "id": "1",
          "nickname": "p1ck0"
      }
    }
  }
}
```

## Get Rooms

You get information about the rooms to which you have access.

### **Shemas**

```
type Room {
  id: ID!
  unique_name: String!
  name: String!
  link: String!
  owner: User!
  editors: [User!]
  users: [User!]
}
```

### **Query request**

```
getRooms(id: [ID!]): [Room!]
```
### **Example**

Request

```
query {
    getRooms(id: ["1"]) {
        id
        unique_name
        name
        link
        owner {
            id
            nickname //the whole user structure can be queried
        }
        editors {
            id
            nickname //the whole user structure can be queried
        }
        users {
            id
            nickname //the whole user structure can be queried
        }
    }
}
```

Response

```
{
  "data": {
    "getRooms": {
      "id": "1",
      "unique_name": "p1ck0_room",
      "name": "p1ck0 room",
      "link": "domain/p1ck0_room",
      "owner": {
          "id": "1",
          "nickname": "p1ck0"
      },
      "editors": {
          "id": "1",
          "nickname": "p1ck0"
      },
      "users": {
          "id": "1",
          "nickname": "p1ck0"
      }
    }
  }
}
```

## Add users to room

Adds a users to the room you are the owner or editor.

### **Shemas**

```
input UsersToRoom {
  roomID: ID!
  usersID: [ID!]
}
```

### **Mutations request**

```
addUsersToRoom(
    input: UsersToRoom!
): String!
```

### **Example**

Request

```
mutation {
    addUsersToRoom(
        input: {
            roomID: "1",
            usersID: ["2", "3"]        
        }
    )
}
```

Response

```
{
  "data": {
    "addUsersToRoom": {
        "users_added"
    }
  }
}
```


## Add editors to room

Adds a editors to the room you are the owner.

### **Shemas**

```
input UsersToRoom {
  roomID: ID!
  usersID: [ID!]
}
```

### **Mutations request**

```
addEditorsToRoom(
    input: UsersToRoom!
): String!
```

### **Example**

Request

```
mutation {
    addEditorsToRoom(
        input: {
            roomID: "1",
            usersID: ["2", "3"]        
        }
    )
}
```

Response

```
{
  "data": {
    "addEditorsToRoom": {
        "users_became_editors"
    }
  }
}
```

## Kick users from room.

You deny the user access to the room. If you are the owner or editor of a room.

### **Shemas**

```
input UsersToRoom {
  roomID: ID!
  usersID: [ID!]
}
```

### **Mutations request**

```
kickUsersFromRoom(
    input: UsersToRoom!
): String!
```

### **Example**

Request

```
mutation {
    kickUsersFromRoom(
        input: {
            roomID: "1",
            usersID: ["2", "3"]        
        }
    )
}
```

Response

```
{
  "data": {
    "kickUsersFromRoom": {
        "users_kicked"
    }
  }
}
```

## Delete room

Removing a room if you are the owner


### **Mutations request**

```
deleteRoom(id: ID!): String!
```

### **Example**

Request

```
mutation {
    deleteRoom(id: "1")
}
```

Response

```
{
  "data": {
    "deleteRoom": {
        "room_deleted"
    }
  }
}
```
