# methods for working users

## Create User

creates a new user in mocsi with a unique nickname and email address. Also returns 2 tokens

### **Shemas**

```
type NewUser {
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
      "here.will.be.accessToken",
      "here.will.be.refreshToken"
    }
  }
}
```


## User authorization

checks if the user is in mocsi, if the entered data is correct, then returns two tokens

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
      "here.will.be.accessToken",
      "here.will.be.refreshToken"
    }
  }
}
```

## Error response

returned in case of error

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