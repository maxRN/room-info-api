# Better Room API for TU Dresden

This projects aims to provide an easy-to-use API to get information about the room plans of TU Dresden.

## Getting Started

1. Have Go installed on your system.
1. For local development you should setup a local database.

   1. Go to your terminal and type

   ```shell
   docker volume create room-info
   ```

   1. Then type

   ```shell
   docker run --name room-info-mysql \
   -v room-info:/var/lib/mysql \
   -e MYSQL_ROOT_PASSWORD=booster123 \
   --restart on-failure \
   -p 31444:3306 \
   -d mysql:latest
   ```

   1. Then run `./migrations.sh` to setup the database and run all migrations.

1. Create a copy of `.env.example` called `.env.local` and add your TU username and password in there.
1. In your terminal type `go run .` to start the local server.

## Limitations

Only a handful rooms of the APB are currently supported.
Adding new rooms takes a bit of manual work, PRs are appreciated :)

## Outlook

The following improvements/features are planned:

- An endpoint to find free rooms
- Refactoring the SQL database schema
- Adding some tests, especially for the parsing of the HTML code
- Migrate to GORM
