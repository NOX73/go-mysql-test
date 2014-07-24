go-mysql-test
=============

Application can concurrently download data from mysql servers, dump it to cvs and zip files to arhive.

Commands:
  - `make deps` install dependencies

Run application:
==============

```sh
go run main.go "<login1>:<pass1>@tcp(<ip_addres1>:<port1>)/<db_name1>" "<login2>:<pass2>@tcp(<ip_addres2>:<port2>)/<db_name2>"
```

Helpful
=============

To running mysql servers you can use make files command:
  - `mysql-build` - downloads docker image and install it.
  - `mysql-run` - runs two mysql servers inside docker containers on 3306 and 3307 ports.
  - `mysql-dump` - upload sql dump with random data to both mysql servers.

If you use docker mysql servers, you can run `IP=127.0.0.1 make run` to run application with default configuration.
  
Run `make` command for more details.
