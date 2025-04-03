# MySQL process exporter

## exposes metrics

### `mysql_process_seconds`

labels:

- `db_host`: MySQL host
- `id`: MySQL process id
- `user`: MySQL user
- `host`: MySQL host
- `db`: MySQL database
- `command`: MySQL command
- `state`: MySQL process state
- `info`: MySQL process info

## Configuration

- Environment variables
  - `MYSQL_USER`: MySQL user
  - `MYSQL_PASSWORD`: MySQL password
  - `MYSQL_HOST`: MySQL host
  - `MYSQL_PORT`: MySQL port
  - `MYSQL_PROCESS_EXPORTER_PORT`: MySQL process exporter port (default: 8080)

## How to use in production

use <https://github.com/Drumato/helm-charts>

## How to run in local

```bash
# terminal 1
docker compose up -d
MYSQL_USER=user MYSQL_PASSWORD=userpassword MYSQL_HOST=127.0.0.1 MYSQL_PORT=3306 ./mysql-process-exporter.exe

# terminal 2
docker exec -it <mysql_container_id> mysql -u user -p -e "SELECT SLEEP(1000);" exampledb

# terminal 3
curl http://localhost:8080/metrics
```
