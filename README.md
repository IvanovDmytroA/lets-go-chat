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
