source .env

echo "${DB_HOST}:${DB_PORT}"

# wait-for "${DB_HOST}:${DB_PORT}" -- "$@"

# Watch your .go files and invoke go build if the files changed.
CompileDaemon --build="go build -o main main.go"  --command=./main