setup: setup-mysql-servers

setup-mysql-servers: pull-tatum-mysql build-mysql-docker

pull-tatum-mysql:
	git clone https://github.com/tutumcloud/tutum-docker-mysql.git

build-mysql-docker:
	cd tutum-docker-mysql && docker build -t tutum/mysql .

run-mysql:
	docker run -d --name mysql1 -p 3306:3306 tutum/mysql
	docker run -d --name mysql2 -p 3307:3306 tutum/mysql

clean:
	docker rm -f mysql1 mysql2

