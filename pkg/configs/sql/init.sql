drop table if exists `follow`;
drop table if exists `videos`;
drop table if exists `user`;
drop table if exists `message`;
-- 用户表 users --
create table if not exists `user`
(
    `id`         BIGINT       NOT NULL AUTO_INCREMENT COMMENT 'PK',
    `name`       VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'UserName',
    `password`   VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'UserPassword',
    `created_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'User account create time',
    `updated_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'User account update time',
    `deleted_at` TIMESTAMP    NULL     DEFAULT NULL COMMENT 'User account delete time',
    primary key (`id`),
    key `name_password_idx` (`name`, `password`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='User table';

create table if not exists `follow`
(
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id`         BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `follow_id`       BIGINT UNSIGNED NOT NULL COMMENT '关注的用户ID',
     primary key (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='关注表 user_id 关注了 follow_id';

create table if not exists `message`
(
    `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '消息id',
    `to_user_id `   BIGINT UNSIGNED NOT NULL COMMENT '接收者的id',
    `from_user_id ` BIGINT UNSIGNED NOT NULL COMMENT '该发送者的id',
    `content`       LONGTEXT NOT NULL COMMENT '消息内容',
    `created_at`    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    primary key (`id`)
)ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_general_ci COMMENT ='消息表';


create table if not exists `videos`
(
    `id`              bigint primary key auto_increment comment '视频ID',
    `created_at`      timestamp not null default current_timestamp comment 'User account create time',
    `updated_at`      timestamp not null default current_timestamp on update current_timestamp comment 'User account update time',
    `deleted_at`      timestamp null     default null comment 'User account delete time',
    `title`           varchar(255) comment '视频标题',
    `label`           varchar(255) comment '视频标签',
    `owner_id`        bigint    not null comment '视频所有者',
    `likes_count`     bigint             default 0 comment '点赞人数',
    `comment_archive` longtext comment '评论信息的最新归档',
    foreign key (owner_id) references `user` (id)
)