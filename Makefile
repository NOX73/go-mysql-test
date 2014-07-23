MYSQL_PASS=12345

help:
	cat Makefile

setup: mysql-build deps mysql-run

mysql-build:
	git clone https://github.com/tutumcloud/tutum-docker-mysql.git
	cd tutum-docker-mysql && docker build -t tutum/mysql .

mysql-run:
	docker run -d --name mysql1 -p 3306:3306 -e MYSQL_PASS="$(MYSQL_PASS)" tutum/mysql
	docker run -d --name mysql2 -p 3307:3306 -e MYSQL_PASS="$(MYSQL_PASS)" tutum/mysql

mysql-logs:
	@((docker logs mysql1 && docker logs mysql2) | grep "uadmin")

mysql-stop:
	docker stop mysql1
	docker stop mysql2

mysql-start:
	docker start mysql1
	docker start mysql2

clean:
	docker rm -f mysql1
	docker rm -f mysql2

run: 
	go run main.go "admin:$(MYSQL_PASS)@tcp($(IP):3306)/hello" "admin:$(MYSQL_PASS)@tcp($(IP):3307)/hello"

deps:
	go get github.com/go-sql-driver/mysql
