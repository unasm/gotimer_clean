

CREATE TABLE `black` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `ip` varchar(24) NOT NULL DEFAULT '',
    `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0表示添加成功，3同步成功 7同步失败，10标记删除，14删除成功',
    `desc` varchar(256) NOT NULL DEFAULT '',
    `creator` varchar(64) NOT NULL DEFAULT '',
    `create_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `ip` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8  ;


CREATE TABLE `dim` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '维度名称',
    `desc` varchar(512) NOT NULL DEFAULT '',
    `user` varchar(32) NOT NULL DEFAULT '',
    `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0未授权1已生效2未生效',
    `create_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    `update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;  


CREATE TABLE `rule` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `uri` varchar(256) NOT NULL DEFAULT '' COMMENT 'URI',
    `dim_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '限制元素，0memberId, 1表示memberid+IP 2表示设备标识udid',
    `user` varchar(64) NOT NULL DEFAULT '',
    `times` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '次数',
    `expire` int(11) NOT NULL DEFAULT '0' COMMENT '超时周期, 分、小时、天',
    `way` int(11) NOT NULL DEFAULT '0' COMMENT '阻断方式,1.图形验证码、2.短信验证码、3.直接返回',
    `ext` varchar(512) NOT NULL DEFAULT '扩展字段、建议存放json',
    `reason` varchar(512) NOT NULL DEFAULT '' COMMENT '创建、修改策略原因',
    `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0表示添加成功，3同步成功 7同步失败，10标记删除，14删除成功',
    `create_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8  ;
