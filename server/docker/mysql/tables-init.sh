#!/bin/bash
set -e

mysql --protocol=socket -uroot -p$MYSQL_ROOT_PASSWORD <<EOSQL

USE ${MYSQL_DATABASE};

CREATE TABLE station (
    station_id INT PRIMARY KEY AUTO_INCREMENT,
    station_owner VARCHAR(255) NOT NULL,
    latitude DECIMAL(9, 6) NOT NULL CHECK (latitude >= -90 AND latitude <= 90),
    longitude DECIMAL(10, 6) NOT NULL CHECK (longitude >= -180 AND longitude <= 180)
);

CREATE TABLE measurement (
    measurement_id INT PRIMARY KEY AUTO_INCREMENT,
    random_num INT NOT NULL
);

CREATE TABLE entry (
    entry_id INT PRIMARY KEY AUTO_INCREMENT,

    station_id INT NOT NULL,
    latitude DECIMAL(9, 6) NOT NULL CHECK (latitude >= -90 AND latitude <= 90),
    longitude DECIMAL(10, 6) NOT NULL CHECK (longitude >= -180 AND longitude <= 180),

    measurement_id INT NOT NULL,

    entry_time DATETIME NOT NULL,

    FOREIGN KEY (station_id) REFERENCES station (station_id),
    FOREIGN KEY (measurement_id) REFERENCES measurement (measurement_id)
);

EOSQL