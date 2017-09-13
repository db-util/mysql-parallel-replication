DROP DATABASE IF EXISTS adt_test;

CREATE DATABASE adt_test;

SET GLOBAL binlog_format = 'ROW';
-- SET GLOBAL log_bin = 'YES';

SHOW variables WHERE variable_name LIKE "%log%";

RESET MASTER;

CREATE TABLE adt_test.test_1 (
  no int(11) NOT NULL,
  seq int(11) NOT NULL,
  uk int(11) NOT NULL,
  v text NOT NULL,
  c int(11) NOT NULL,
  modtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  regtime datetime NOT NULL,
  PRIMARY KEY (no),
  UNIQUE KEY ux_uk_no (uk),
  KEY ix_modtime (modtime),
  KEY ix_regtime (regtime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SHOW databases \G;
SHOW tables FROM adt_test \G;
SHOW CREATE TABLE adt_test.test_1 \G;