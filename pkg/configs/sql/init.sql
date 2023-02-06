-- 用户表 users --
CREATE TABLE `user`
(
    `id`         bigint NOT NULL AUTO_INCREMENT COMMENT 'PK',
    `name`       varchar(128) NOT NULL DEFAULT '' COMMENT 'UserName',
    `password`   varchar(128) NOT NULL DEFAULT '' COMMENT 'UserPassword',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'User account create time',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'User account update time',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'User account delete time',
     PRIMARY KEY (`id`),
     KEY `name_password_idx` (`name`,`password`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='User table';

# CREATE TABLE `follows`
# (
#
# )