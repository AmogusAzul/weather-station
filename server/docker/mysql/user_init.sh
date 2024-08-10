#!/bin/bash
set -e

mysql --protocol=socket -uroot -p$MYSQL_ROOT_PASSWORD <<EOSQL

-- Create both users for the servers
CREATE USER '${DB_DATA_SERVER_USER}'@'%' IDENTIFIED WITH 'caching_sha2_password' BY '${DB_DATA_SERVER_PASSWORD}';
CREATE USER '${DB_USER_SERVER_USER}'@'%' IDENTIFIED WITH 'caching_sha2_password' BY '${DB_USER_SERVER_PASSWORD}';

-- Granting write privileges to the data-server user
GRANT ALL PRIVILEGES ON ${MYSQL_DATABASE}.* TO '${DB_DATA_SERVER_USER}'@'%';

-- Grant only-read permissions to the user-server user for safety reasons
GRANT SELECT ON ${MYSQL_DATABASE}.* TO '${DB_USER_SERVER_USER}'@'%';

FLUSH PRIVILEGES;

EOSQL