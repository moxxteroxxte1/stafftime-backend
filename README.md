## Nas Company Duty Roster and Time Sheet API

### Get all Users

#### Request URL
```http
  GET /api/users
```

#### Request Parameters
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `api_key` | `string` | **Required**. Your API key |

#### Respone
````JSON
    [
      {
        "id"
      }
    ]
````    

#### Create User

```http
 POST /api/users
```

| Parameter   | Type     | Description                          |
|:------------| :------- |:-------------------------------------|
| `firstname` | `string` | **Required**.                        |
| `lastname`  | `string` | **Required**.                        |
| `username`  | `string` | **Required** Username must be unique |
| `email`     | `string` | **Required**.                        |
| `password`  | `string` | **Required**.                        |

#### Update all Users

````http requeuest
    PUT /api/users
````

| Parameter   | Type     | Description                          |
|:------------| :------- |:-------------------------------------|
| `firstname` | `string` | **Optional**.                        |
| `lastname`  | `string` | **Optional**.                        |
| `username`  | `string` | **Optional** Username must be unique |
| `email`     | `string` | **Optional**.                        |
| `password`  | `string` | **Optional**.                        |

#### Get User

```http
  GET /api/users/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |

#### add(num1, num2)

Takes two numbers and returns the sum.
