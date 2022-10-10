#!/bin/bash

fb_authentication = `CREATE TABLE "fb_authentication"
(
    "id"            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name"          TEXT    default '',
    "auth_supplier" TEXT    default '',
    "config"        TEXT    default '',
    "create_time"   TEXT    DEFAULT "2022-08-09 00:00:00",
    "update_time"   TEXT    DEFAULT "2022-08-09 00:00:00",
    "is_del"        INTEGER DEFAULT 0,
    "switch_state"  TEXT    default ''
);`


fb_data_source = `CREATE TABLE "fb_data_source"
(
    "id"          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "source_type" integer DEFAULT 0,
    "switch"      integer DEFAULT 0,
    "is_del"      integer DEFAULT 0,
    "name"        TEXT    DEFAULT '0',
    "config"      TEXT    DEFAULT '0',
    "create_time" TEXT    DEFAULT "2022-08-09 00:00:00",
    "update_time" TEXT    DEFAULT "2022-08-09 00:00:00"
);`


fb_env = `CREATE TABLE "fb_env"
(
    "id"          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "env_type"    integer DEFAULT 0,
    "is_del"      integer DEFAULT 0,
    "key"         TEXT    DEFAULT '',
    "dev_env"     TEXT    DEFAULT '',
    "pro_env"     TEXT    DEFAULT '',
    "create_time" TEXT    DEFAULT "2022-08-09 00:00:00",
    "update_time" TEXT    DEFAULT "2022-08-09 00:00:00"
);`


fb_file = `CREATE TABLE "fb_file"
(
    "id"          text NOT NULL PRIMARY KEY,
    "name"        TEXT    DEFAULT '',
    "path"        text    DEFAULT '',
    "create_time" TEXT    DEFAULT "2022-08-09 00:00:00",
    "update_time" TEXT    DEFAULT "2022-08-09 00:00:00",
    "is_del"      integer DEFAULT 0

);`

fb_operations = `CREATE TABLE "fb_operations"
(
    "id"             INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "method"         TEXT    DEFAULT '',
    "is_public"      integer DEFAULT 0,
    "remark"         TEXT    DEFAULT '',
    "legal"          integer DEFAULT 0,
    "path"           TEXT    DEFAULT '',
    "content"        TEXT    DEFAULT '',
    "enable"         integer DEFAULT 0,
    "create_time"    TEXT    DEFAULT "2022-08-09 00:00:00",
    "update_time"    TEXT    DEFAULT "2022-08-09 00:00:00",
    "is_del"         integer DEFAULT 0,
    "operation_type" TEXT
);`

fb_prisma = `CREATE TABLE "fb_prisma"
(
    "id"          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "file_id"     integer DEFAULT 0,
    "name"        TEXT    DEFAULT '',
    "create_time" TEXT    DEFAULT "2022-08-09 00:00:00",
    "update_time" TEXT    DEFAULT "2022-08-09 00:00:00",
    "is_del"      integer DEFAULT 0
);`

fb_role = `CREATE TABLE "fb_role"
(
    "id"     INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "code"   TEXT    DEFAULT '',
    "remark" TEXT    DEFAULT '',
    "create_time"    DEFAULT "2022-08-09 00:00:00",
    "update_time"    DEFAULT "2022-08-09 00:00:00",
    "is_del" integer DEFAULT 0,
);`

fb_storage_bucket = `CREATE TABLE "fb_storage_bucket"
(
    "id"          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name"        TEXT    default '',
    "switch"      TEXT    default '',
    "config"      TEXT    default '',
    "create_time" TEXT    DEFAULT '2022-07-07 11:00:00',
    "update_time" TEXT    DEFAULT '2022-07-07 11:00:00',
    "is_del"      integer DEFAULT 0
);`

fb_user = `CREATE TABLE "fb_user"
(
    "id"                  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name"                TEXT    DEFAULT '',
    "phone"               TEXT    DEFAULT '',
    "email"               TEXT    DEFAULT '',
    "account"             TEXT    DEFAULT '',
    "encryption_password" TEXT    DEFAULT '',
    "create_at"           TEXT    DEFAULT '',
    "update_at"           TEXT    DEFAULT '',
    "is_del"              integer default 0,
);`

oauth_user=`CREATE TABLE `oauth_user`
(
    "id"                  INTEGER unsigned NOT NULL,
    "name"                TEXT NOT NULL DEFAULT '' ,
    "nick_name"           TEXT NOT NULL DEFAULT '' ,
    "user_name"             TEXT NOT NULL DEFAULT '' ,
    "encryption_password" TEXT NOT NULL DEFAULT '' ,
    "mobile"              TEXT NOT NULL default '' ,
    "email"               TEXT NOT NULL default '' ,
    "mateData"               TEXT NOT NULL default '' ,
    "last_login_time"     TEXT   NOT NULL ,
    "status"              INTEGER     not null default 0 ,
    "create_time"         TEXT   NOT NULL DEFAULT '' ,
    "update_time"         TEXT NULL DEFAULT '',
    "is_del"              INTEGER unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY ("id")
);`

sqlite3 /home/fire_boom.db