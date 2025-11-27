#!/bin/bash

set -e

mysql -u root -p"${MYSQL_ROOT_PASSWORD}" << EOSQL
USE ${MYSQL_DATABASE};
CREATE TABLE employees_tbl (
    id VARCHAR(20) NOT NULL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    is_admin BOOLEAN NOT NULL,
    joining_date DATETIME NOT NULL,
    registered_at DATETIME NOT NULL,
    deleted_at DATETIME NULL
);
CREATE TABLE paid_leave_employee_tbl (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    employee_id VARCHAR(20) NOT NULL,
    vacation_date DATETIME NOT NULL,
    start_at_hour INT NOT NULL,
    duration INT NOT NULL,
    given_at DATETIME NOT NULL,
    registered_at DATETIME NOT NULL,
    deleted_at DATETIME NULL,
    FOREIGN KEY (employee_id) REFERENCES employees_tbl(id)
);
CREATE TABLE paid_leave_company_tbl (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    vacation_date DATETIME NOT NULL,
    registered_at DATETIME NOT NULL,
    deleted_at DATETIME NULL
);
CREATE DATABASE ${MYSQL_AUTH_DATABASE};
USE ${MYSQL_AUTH_DATABASE};
CREATE TABLE auth_tbl (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    employee_id VARCHAR(20) NOT NULL,
    hash VARCHAR(64) NOT NULL,
    UNIQUE(employee_id),
    FOREIGN KEY (employee_id) REFERENCES ${MYSQL_DATABASE}.employees_tbl(id)
);
EOSQL
