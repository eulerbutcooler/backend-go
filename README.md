# backend-go

A basic CRUD backend in go.

*To Run:*
`docker compose up --build`

*APIs:*

1. POST /users - To add new user

Sample JSON - 
```
{
"name":"am",
"email":"am@gmail.com"
}
```
---
2. GET /users - To get all the users
---
3. GET /users/id - To get a specific user
---
4. PUT /users/id - To update a specific user

Sample JSON -
```
{
"name":"p"
"email":"p@gmail.com"
}
```
---
5. DELETE /users/id - To delete a specific user
