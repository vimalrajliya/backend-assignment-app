# Backend-Assignment-App

This Backend Application  contain Authentication APIs build with 
- **Go**
- **Gin (Web Framework)**
- **GORM (ORM for Go)**
- **SQLite** (Database)
- **Redis** (for token blacklisting)

Docker has been implemented to run the project.

## Getting started
To run the project using Docker, ensure Docker is installed and running on your system. Then execute:

`docker compose up`


APIs postman documentation can be found on the following link:

https://tinyurl.com/yxm3j7ae


Following is the discription of APIs:

### 1. Sign up (creation of user)

This API is for creating the sign up using name, email and password.
user gets created in sqlite db and has uniqness validation on email.
Also the password is stored with encryption.

`curl --location 'localhost:8080/auth' \
--header 'Content-Type: application/json' \
--data-raw '{
        "name": "vimal test",
        "email": "vimaltest@gmail.com",
        "password": "sdfasddsfnkasdsnfl424fds34"
    }'
`

### 2. Sign In

User can get the authentication token using this API provided correct email and password.

`
curl --location 'localhost:8080/auth/sign-in' \
--header 'Content-Type: application/json' \
--data-raw '{
        "email": "vimaltest@gmail.com",
        "password": "sdfasddsfnkasdsnfl424fds34"
}'
`

### 3. Authorization of token (Get User Details)

In this user details can be accessed by authentication token in bearer.

`
curl --location 'localhost:8080/api/user/' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InZpbWFsdGVzdEBnbWFpbC5jb20iLCJleHAiOjE3NDc3NTUyODQsImlhdCI6MTc0NzQ5NjA4NH0.kHUbnHEHpQ11KAI9nGkrAXwmVjThAqUEIHGHxd4yx04'
`

### 4. Mechanism to refresh a token

This API generates fresh token which can be used to extend the session
`
curl --location 'localhost:8080/api/token/refresh' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InZpbWFsdGVzdEBnbWFpbC5jb20iLCJleHAiOjE3NDc3NTUyODQsImlhdCI6MTc0NzQ5NjA4NH0.kHUbnHEHpQ11KAI9nGkrAXwmVjThAqUEIHGHxd4yx04'
`

### 5. Revocation of token (Logout)

This API puts token in blacklist for remaining expired times. Redis is used here for fast retrieval.

`
curl --location --request POST 'localhost:8080/api/user/log-out' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InZpbWFsdGVzdEBnbWFpbC5jb20iLCJleHAiOjE3NDc3NTUyODQsImlhdCI6MTc0NzQ5NjA4NH0.kHUbnHEHpQ11KAI9nGkrAXwmVjThAqUEIHGHxd4yx04'
`

