CREATE TABLE `friends` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `friend_uid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `remark` varchar(255) DEFAULT NULL,
    `add_source` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `friend_requests` (
    `id` int (11) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` varchar(64)COLLATE utf8mb4_unicode_ci NOT NULL ,
    `req_uid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL ,
    `req_msg` varchar(255) DEFAULT NULL,
    `req_time` timestamp NOT NULL,
    `handle_result` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `handle_msg` varchar(255) DEFAULT NULL,
    `handled_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `groups`(
    `id`int (11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `icon`varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `status`tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `creator_uid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `group_type`int (11) unsigned NOT NULL,
    `is_verify` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `notification` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `notification_uid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    primary key (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `group_requests`(
     `id` int (11) unsigned NOT NULL AUTO_INCREMENT,
     `group_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
     `req_uid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL ,
     `req_msg` varchar(255) DEFAULT NULL,
     `req_time` timestamp NOT NULL,
     `handle_uid` int (11) unsigned NOT NULL,
     `handle_result` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
     `handle_msg` varchar(255) DEFAULT NULL,
     `handled_at` timestamp NULL DEFAULT NULL,
     `join_source` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
     `inviter_uid` varchar(64)COLLATE utf8mb4_unicode_ci NOT NULL,
    primary key (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `group_members`(
    `id` int (11) unsigned NOT NULL AUTO_INCREMENT,
    `group_id`varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `user_id`varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `role_level`tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `join_time`timestamp NULL DEFAULT NULL,
    `join_source`varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `inviter_uid`varchar(64)COLLATE utf8mb4_unicode_ci NOT NULL,
    `operator_uid`varchar(64)COLLATE utf8mb4_unicode_ci NOT NULL,
    primary key (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;