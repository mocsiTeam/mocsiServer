# methods for working users

## Create User

creates a new user in mocsi with a unique nickname and email address. Also returns 2 tokens.

### **Shemas**

```
input NewUser {
    email: String!
    password: String!
    nickname: String!
    firstname: String!
    lastname: String!
}

type Tokens {
    accessToken: String!
    refreshToken: String!
}
```

### **Mutations request**

```
createUser(
    input: NewUser!
): Tokens!
```

### **Example**

Request

```
mutation {
    createUser(input: {email: "p1ck0@github.com", password: "secrtet", 
    nickname: "p1ck0", firstname: "Kameha", lastname: "Meha"})
  {
    accessToken
    refreshToken
  }
}
```

Response

```
{
  "data": {
    "createUser": {
      "accessToken":"here.will.be.accessToken",
      "refreshToken":"here.will.be.refreshToken"
    }
  }
}
```


## User authorization

checks if the user is in mocsi, if the entered data is correct, then returns two tokens.

### **Shemas**

```
type Tokens {
    accessToken: String!
    refreshToken: String!
}

type Login {
    nickname: String!
    password: String!
}
```

### **Mutations request**

```
login(
    input: Login!
): Tokens!
```

### **Example**

Request

```
mutation {
    login(input: {nickname: "p1ck0", password: "sercet"}) {
      accessToken
      refreshToken
    }
}
```

Response

```
{
  "data": {
    "login": {
      "accessToken":"here.will.be.accessToken",
      "refreshToken":"here.will.be.refreshToken"
    }
  }
}
```

## Refresh token

Updates the access token.

### **Mutations request**
```
input RefreshTokenInput{
  token: String!
}
```

```
refreshToken(
  input: RefreshTokenInput!
): String!
```

### **Example**

Request

```
mutation {
    refreshToken(input: RefreshTokenInput!)
}
```

Response

```
{
  "data": {
    "login": {
      "here.will.be.accessToken"
    }
  }
}
```

## Getting an authorized user

Allows you to get information about the current user.

For this request you will need to add a token to the HTTP headers.
```
  {
    "token":"here.will.be.accessToken"
  }
```
### **Shemas**

```
type User {
  id: ID!
  nickname: String!
  firstname: String!
  lastname: String!
  email: String!
  role: String!
  groups: [String!]!
  error: String!
}
```

### **Query request**

```
getAuthUser: User!
```

### **Example**

Request

```
query {
  getAuthUser {
    id
    nickname
    firstname
    lastname
    email
    role
    groups
    error
  }
}
```

Response

```
{
  "data": {
    "getAuthUser": {
      "id": "1",
      "nickname": "p1ck0",
      "firstname": "Kameha",
      "lastname": "Meha",
      "email": "p1ck0@github.com",
      "role": "1",
      "groups": [],
      "error": ""
    }
  }
}
```

## Error response

Returned in case of error.

```
{
  "errors": [
    {
      "message": "error text",
      "path": [
        "createUser"
      ],
      "extensions": {
        "code": error code(has an integer type),
        "status": false
      }
    }
  ],
  "data": null
}
```
