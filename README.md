# lets-go-chat

## How to run
1. Clone the project.
2. Start Postgresql Docker container:
<pre>
docker run -it --rm --name go-postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=pass -e PGDATA=/var/lib/postgresql/data/pgdata -v ~/local-go-postgres:/var/lib/postgresql/data postgres:14.0
</pre>
3. Start Redis Docker Container
<pre>
docker run -d --name redis -p 6379:6379 redis
</pre>

Open terminal from the application directory and :
<pre>
    $ go run cmd/main.go
</pre>

## API
Use Postman or your favourite browser

1. Registration

<pre>
POST http://localhost:8080/v1/user -H 'Content-Type: application/json' -d '{"userName":"name","password":"mpassword"}'
</pre>

2. Login

<pre>
POST http://localhost:8080/v1/user/login -H 'Content-Type: application/json' -d '{"userName":"name","password":"password"}'
</pre>

3. Start chat
Use token returned by Login request

<pre>
ws://localhost:8080/v1/chat/ws.rtm.start?token=your_token
</pre>

4. Get active users

<pre>
GET http://localhost:8080/v1/user/active -H 'Content-Type: application/json'
</pre>
