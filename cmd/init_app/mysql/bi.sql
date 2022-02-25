DROP TABLE IF EXISTS `app`;
CREATE TABLE `app`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  `descibe` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT NULL,
  `app_id` varchar(225) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT NULL,
  `app_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT NULL,
  `create_by` int(11) NULL DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `update_by` int(11) NULL DEFAULT 0,
  `app_manager` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  `is_close` tinyint(4) NULL DEFAULT 0 COMMENT '是否关闭 0为false 1 为 true',
  `save_mouth` int(11) NULL DEFAULT 1 COMMENT '保存n个月',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `app_name`(`app_name`) USING BTREE,
  UNIQUE INDEX `app_id`(`app_id`) USING BTREE,
  INDEX `app_create_by`(`create_by`, `app_name`, `is_close`) USING BTREE,
  INDEX `app_isclose`(`is_close`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 41 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_german2_ci ROW_FORMAT = Dynamic;
DROP TABLE IF EXISTS `attribute`;
CREATE TABLE `attribute`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `attribute_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '' COMMENT '属性名',
  `show_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '' COMMENT '显示名',
  `data_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '' COMMENT '数据类型',
  `attribute_type` tinyint(4) NULL DEFAULT 1 COMMENT '默认为1 （1为预置属性，2为自定义属性）',
  `attribute_source` tinyint(4) NULL DEFAULT 1 COMMENT '默认为1 （1为用户属性，2为事件属性）',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `app_id` int(11) NULL DEFAULT 0 COMMENT 'appid',
  `status` tinyint(4) NULL DEFAULT 0 COMMENT '是否显示 0为不显示 1为显示 默认不显示',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `attribute_name_attribute_source`(`attribute_name`, `attribute_source`, `app_id`) USING BTREE,
  INDEX `attribute_id_source`(`app_id`, `attribute_source`, `attribute_name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4022 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_german2_ci ROW_FORMAT = Dynamic;
DROP TABLE IF EXISTS `debug_device`;
CREATE TABLE `debug_device`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `appid` int(11) NULL DEFAULT 0,
  `device_id` varchar(225) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT NULL,
  `create_by` int(11) NULL DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `debug_device_uq`(`appid`, `device_id`) USING BTREE,
  INDEX `debug_device_appid_createby`(`appid`, `create_by`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_german2_ci ROW_FORMAT = Dynamic;
DROP TABLE IF EXISTS `gm_operater_log`;
CREATE TABLE `gm_operater_log`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `operater_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '操作者名字',
  `operater_id` int(11) NULL DEFAULT 0 COMMENT '操作者id',
  `operater_action` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '请求路由',
  `created` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `method` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '请求方法',
  `body` blob NOT NULL COMMENT '请求body',
  `operater_role_id` int(11) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `operater_action`(`operater_action`) USING BTREE,
  INDEX `operater_id`(`operater_id`) USING BTREE,
  INDEX `operater_role_id`(`operater_role_id`) USING BTREE,
  INDEX `operater_id_act_role`(`operater_action`, `operater_id`, `operater_role_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2940 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = DYNAMIC;
DROP TABLE IF EXISTS `gm_role`;
CREATE TABLE `gm_role`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `role_list` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;
INSERT INTO `gm_role` VALUES (1, 'admin', '超级管理员', '[{\"path\":\"/behavior-analysis\",\"component\":\"layout\",\"redirect\":\"/behavior-analysis/index\",\"alwaysShow\":false,\"meta\":{\"title\":\"行为分析\",\"icon\":\"el-icon-link\"},\"children\":[{\"path\":\"event/:id\",\"component\":\"views/behavior-analysis/event\",\"name\":\"event\",\"meta\":{\"title\":\"事件分析\",\"dynamic\":true,\"icon\":\"el-icon-data-line\"}},{\"path\":\"retention/:id\",\"component\":\"views/behavior-analysis/retention\",\"name\":\"retention\",\"meta\":{\"title\":\"留存分析\",\"dynamic\":true,\"icon\":\"el-icon-data-analysis\"}},{\"path\":\"funnel/:id\",\"component\":\"views/behavior-analysis/funnel\",\"name\":\"funnel\",\"meta\":{\"title\":\"漏斗分析\",\"dynamic\":true,\"icon\":\"el-icon-data-board\"}},{\"path\":\"trace/:id\",\"component\":\"views/behavior-analysis/trace\",\"name\":\"trace\",\"meta\":{\"title\":\"智能路径分析\",\"dynamic\":true,\"icon\":\"el-icon-bicycle\"}}]},{\"path\":\"/user-analysis\",\"component\":\"layout\",\"redirect\":\"/user-analysis/attr\",\"alwaysShow\":false,\"meta\":{\"title\":\"用户分析\",\"icon\":\"el-icon-pie-chart\"},\"children\":[{\"path\":\"attr/:id\",\"component\":\"views/user-analysis/index\",\"name\":\"attr\",\"meta\":{\"title\":\"用户属性分析\",\"dynamic\":true,\"icon\":\"el-icon-s-custom\"}},{\"path\":\"group\",\"component\":\"views/user-analysis/group\",\"name\":\"group\",\"meta\":{\"title\":\"用户分群\",\"icon\":\"el-icon-user\"}},{\"isInside\":true,\"path\":\"user_list\",\"component\":\"views/user-analysis/user_list\",\"name\":\"user_list\",\"meta\":{\"title\":\"用户列表\",\"icon\":\"el-icon-user-solid\"}},{\"isInside\":true,\"path\":\"user_info/:uid/:index\",\"component\":\"views/user-analysis/user_info\",\"name\":\"user_info\",\"meta\":{\"title\":\"用户事件详情\",\"dynamic\":true,\"icon\":\"el-icon-s-custom\"}}]},{\"path\":\"/manager\",\"component\":\"layout\",\"redirect\":\"/manager/event\",\"alwaysShow\":false,\"meta\":{\"title\":\"数据管理\",\"icon\":\"el-icon-edit\"},\"children\":[{\"path\":\"event\",\"component\":\"views/manager/event\",\"name\":\"event\",\"meta\":{\"title\":\"事件管理\",\"icon\":\"el-icon-s-management\"}},{\"path\":\"log\",\"component\":\"views/manager/log\",\"name\":\"log\",\"meta\":{\"title\":\"埋点管理\",\"icon\":\"el-icon-notebook-1\"}}]},{\"path\":\"/permission\",\"component\":\"layout\",\"redirect\":\"/permission/role\",\"alwaysShow\":true,\"meta\":{\"title\":\"权限\",\"icon\":\"el-icon-user-solid\"},\"children\":[{\"path\":\"role\",\"component\":\"views/permission/role\",\"name\":\"RolePermission\",\"meta\":{\"title\":\"角色管理\",\"icon\":\"el-icon-s-check\"}},{\"path\":\"user\",\"component\":\"views/permission/user\",\"name\":\"user\",\"meta\":{\"title\":\"用户管理\",\"icon\":\"el-icon-user\"}},{\"path\":\"operater_log\",\"component\":\"views/permission/operater_log\",\"name\":\"operater_log\",\"meta\":{\"title\":\"操作日志列表\",\"icon\":\"el-icon-s-order\"}}]},{\"path\":\"/app\",\"component\":\"layout\",\"children\":[{\"path\":\"/app/app\",\"component\":\"views/app/index\",\"name\":\"index\",\"meta\":{\"title\":\"应用管理\",\"icon\":\"el-icon-s-goods\"}}]}]', '2022-02-24 21:03:07', '2022-01-07 14:56:23');
DROP TABLE IF EXISTS `gm_user`;
CREATE TABLE `gm_user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `role_id` int(11) NULL DEFAULT NULL COMMENT '角色id',
  `realname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '真实姓名',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `last_login_time` varchar(225) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
  `is_del` tinyint(4) NULL DEFAULT 0 COMMENT '是否禁止该账号',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `gm_user_username`(`username`) USING BTREE COMMENT '角色名唯一索引',
  INDEX `gm_user_username_pwd`(`username`, `password`, `is_del`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;
INSERT INTO `gm_user` VALUES (1, 'admin', '21232f297a57a5a743894a0e4a801fc3', 1, '肖文龙', '2021-10-21 10:48:08', '2022-01-07 14:49:28', '2022-01-07 14:49:29', 0);
DROP TABLE IF EXISTS `meta_attr_relation`;
CREATE TABLE `meta_attr_relation`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) NULL DEFAULT 0,
  `event_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  `event_attr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `event_name_event_attr`(`app_id`, `event_name`, `event_attr`) USING BTREE,
  INDEX `event_name_event_attr1`(`app_id`, `event_name`) USING BTREE,
  INDEX `event_name_event_attr2`(`app_id`, `event_attr`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1444027 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_german2_ci ROW_FORMAT = Dynamic;
DROP TABLE IF EXISTS `meta_event`;
CREATE TABLE `meta_event`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `appid` int(11) NULL DEFAULT NULL,
  `event_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  `show_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  `yesterday_count` int(11) NULL DEFAULT 0,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `meta_event_appid_event_name`(`appid`, `event_name`) USING BTREE,
  INDEX `meta_event_appid`(`appid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 223627 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_german2_ci ROW_FORMAT = Dynamic;
DROP TABLE IF EXISTS `pannel`;
CREATE TABLE `pannel`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `folder_id` int(11) NULL DEFAULT 0,
  `pannel_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  `managers` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  `create_by` int(11) NULL DEFAULT 0,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `report_tables` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `pannel_unique`(`folder_id`, `pannel_name`, `create_by`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 19 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_german2_ci ROW_FORMAT = Dynamic;
DROP TABLE IF EXISTS `pannel_folder`;
CREATE TABLE `pannel_folder`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `folder_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  `create_by` int(11) NULL DEFAULT 0,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `appid` int(11) NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `pannel_folder_unique`(`folder_name`, `create_by`, `appid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_german2_ci ROW_FORMAT = Dynamic;
DROP TABLE IF EXISTS `report_table`;
CREATE TABLE `report_table`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `appid` int(11) NULL DEFAULT NULL,
  `user_id` int(11) NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  `rt_type` tinyint(8) NULL DEFAULT 0,
  `data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `report_table_appid_user_id_name_type`(`appid`, `user_id`, `name`, `rt_type`) USING BTREE,
  INDEX `report_table_appid_user_id`(`appid`, `user_id`, `rt_type`) USING BTREE,
  INDEX `report_table_id_user_id`(`id`, `user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 56 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_german2_ci ROW_FORMAT = Dynamic;
DROP TABLE IF EXISTS `user_group`;
CREATE TABLE `user_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_name` varchar(255) COLLATE utf8mb4_german2_ci NOT NULL DEFAULT '',
  `group_remark` varchar(255) COLLATE utf8mb4_german2_ci NOT NULL DEFAULT '',
  `create_by` int(11) NOT NULL DEFAULT '0',
  `user_count` int(11) NOT NULL DEFAULT '0',
  `user_list` blob NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `appid` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_group_name` (`group_name`,`appid`) USING BTREE,
  KEY `user_group_appid` (`id`,`appid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_german2_ci

