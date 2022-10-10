#
************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.29)
# Database: fire_boom
# Generation Time: 2022-06-28 13:52:44 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


#
# Dump of table fb_data_source
# ------------------------------------------------------------

DROP TABLE IF EXISTS `fb_data_source`;

CREATE TABLE `fb_data_source`
(
    `id`          int unsigned NOT NULL,
    `name`        varchar(255) NOT NULL DEFAULT '' COMMENT '数据源名称',
    `source_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '数据源类型: mysql, openAPI, sqlLite, pg, graphql, 自定义脚本',
    `config`      text COMMENT '数据源对应的配置项：命名空间、请求配置、连接配置、文件路径、是否配置为外部数据源等',
    `create_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_del`      tinyint unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='数据源配置';


#
# Dump of table fb_file
# ------------------------------------------------------------

DROP TABLE IF EXISTS `fb_file`;

CREATE TABLE `fb_file`
(
    `id`          string       NOT NULL,
    `name`        varchar(255) NOT NULL DEFAULT '' COMMENT '文件名',
    `path`        varchar(255) NOT NULL DEFAULT '' COMMENT '相对路径',
    `create_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_del`      tinyint unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='文件上传 meta 信息表';


#
# Dump of table fb_prisma
# ------------------------------------------------------------

DROP TABLE IF EXISTS `fb_prisma`;

CREATE TABLE `fb_prisma`
(
    `id`          int unsigned NOT NULL,
    `name`        varchar(255) NOT NULL DEFAULT '0' COMMENT 'prisma',
    `file_id`     int unsigned NOT NULL DEFAULT '0' COMMENT '文件 id',
    `create_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_del`      tinyint unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='prisma对应的数据模型';


/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

DROP TABLE IF EXISTS `fb_authentication`;

CREATE TABLE `fb_authentication`
(
    `id`            int unsigned NOT NULL,
    `name`          varchar(255) NOT NULL DEFAULT '' COMMENT '身份验证名称',
    `auth_supplier` string       NOT NULL DEFAULT 'openid' COMMENT '验证供应商:openid、github、google',
    `switch_state`  tinyint unsigned NOT NULL DEFAULT '0' COMMENT '开关状态: 0-全部关闭 1-cookie 2-token 3-全部cookie、token都开启',
    `config`        text COMMENT '身份验证配置对应的配置项：供应商id、appID、appSecret、服务发现地址、重定向url等',
    `create_time`   timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`   timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_del`        tinyint unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='身份验证配置';

DROP TABLE IF EXISTS `fb_role`;

CREATE TABLE `fb_role`
(
    `id`          int unsigned NOT NULL,
    `code`        varchar(20)  NOT NULL DEFAULT '' COMMENT '角色编码',
    `remark`      varchar(100) NOT NULL DEFAULT '' COMMENT '描述',
    `create_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_del`      tinyint unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色配置';



DROP TABLE IF EXISTS `fb_storage_bucket`;

CREATE TABLE `fb_storage_bucket`
(
    `id`          int unsigned NOT NULL,
    `name`        varchar(20)  NOT NULL DEFAULT '' COMMENT '名称',
    `switch`      varchar(20)  NOT NULL DEFAULT '0' COMMENT '开关: 0-关 1-开',
    `config`      varchar(100) NOT NULL DEFAULT '' COMMENT '描述',
    `create_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_del`      tinyint unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='存储配置';


CREATE TABLE `fb_user`
(
    `id`                  int unsigned NOT NULL,
    `name`                varchar(20)  NOT NULL DEFAULT '' COMMENT '名称',
    `phone`               varchar(20)  NOT NULL DEFAULT '0' COMMENT '手机号',
    `email`               varchar(50)  NOT NULL DEFAULT '' COMMENT '邮箱',
    `account`             varchar(20)  NOT NULL DEFAULT '' COMMENT '账号',
    `encryption_password` varchar(255) NOT NULL DEFAULT '' COMMENT '加密密码',
    `create_at`           timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at`           timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_del`              tinyint unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='存储配置';

#
# Dump of table fb_env
# ------------------------------------------------------------

DROP TABLE IF EXISTS `fb_operations`;

CREATE TABLE `fb_operations`
(
    `id`          int unsigned NOT NULL,
    `method`      varchar(255) NOT NULL DEFAULT '' COMMENT '请求类型 GET、POST、PUT、DELETE',
    `status`      tinyint(1) NOT NULL DEFAULT 0 COMMENT '状态 1公有 2私有',
    `remark`      varchar(255) NOT NULL DEFAULT '' COMMENT '说明',
    `legal`       tinyint      NOT NULL DEFAULT 0 COMMENT '是否合法 1合法 2非法',
    `path`        varchar(255) NOT NULL DEFAULT '' COMMENT '路径',
    `content`     varchar(255) NOT NULL DEFAULT '' COMMENT '内容',
    `enable`      tinyint      NOT NULL DEFAULT 0 COMMENT '开关 0开 1关',
    `create_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_del`      tinyint unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='环境变量';


#
# Dump of table oauth_user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `oauth_user`;


CREATE TABLE `oauth_user`
(
    `id`                  int unsigned NOT NULL,
    `name`                varchar(50) NOT NULL DEFAULT '' COMMENT '姓名',
    `nick_name`           varchar(50) NOT NULL DEFAULT '' COMMENT '昵称',
    `user_name`             varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
    `encryption_password` varchar(50) NOT NULL DEFAULT '' COMMENT '加密后密码',
    `mobile`              varchar(11) NOT NULL default '' COMMENT '手机号',
    `email`               varchar(50) NOT NULL default '' COMMENT '邮箱',
    `mate_data`           text        NOT NULL default '' COMMENT '用户信息',
    `last_login_time`     timestamp   NOT NULL COMMENT '最后登陆时间',
    `status`              tinyint     not null default 0 comment '状态',
    `create_time`         timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`         timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_del`              tinyint unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

