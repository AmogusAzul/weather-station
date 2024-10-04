FROM mysql:9.0

# Copy the init scripts to the MySQL init directory inside the container
COPY ./tables-init.sh /docker-entrypoint-initdb.d/tables-init.sh
COPY ./user-init.sh /docker-entrypoint-initdb.d/user-init.sh

# Copy the MySQL configuration file into the MySQL configuration directory
COPY ./my.cnf /etc/mysql/conf.d/my.cnf
