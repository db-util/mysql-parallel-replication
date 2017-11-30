SHOW GLOBAL VARIABLES WHERE variable_name LIKE "%log%";

DROP DATABASE IF EXISTS test_shard_src;
DROP DATABASE IF EXISTS test_shard_dst1;
DROP DATABASE IF EXISTS test_shard_dst2;

CREATE DATABASE test_shard_src;
CREATE DATABASE test_shard_dst1;
CREATE DATABASE test_shard_dst2;

SHOW DATABASES;

CREATE TABLE test_shard_src.t1 (
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

CREATE TABLE test_shard_dst1.t1 LIKE test_shard_src.t1;
CREATE TABLE test_shard_dst2.t1 LIKE test_shard_src.t1;

SHOW TABLES FROM test_shard_src \G;
SHOW CREATE TABLE test_shard_src.t1 \G;

SHOW TABLES FROM test_shard_dst1 \G;
SHOW CREATE TABLE test_shard_dst1.t1 \G;

SHOW TABLES FROM test_shard_dst2 \G;
SHOW CREATE TABLE test_shard_dst2.t1 \G;

RESET MASTER;