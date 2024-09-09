ALTER TABLE messages MODIFY COLUMN `type` tinyint(1) NULL DEFAULT 0 COMMENT '类型 0:通知 1:活动' AFTER `id`,
ALTER TABLE messages MODIFY COLUMN `send_type` tinyint(1) NULL DEFAULT 0 COMMENT '发送类型 1:全体会员 2指定会员' AFTER `contents`,
ALTER TABLE messages ADD COLUMN `img_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '图标地址' AFTER `last_admin`,
ALTER TABLE messages MODIFY COLUMN `send_target` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '发送目标' AFTER `send_type`,
ALTER TABLE messages ADD COLUMN `status` tinyint(3) NULL DEFAULT 1 COMMENT '1 启用 2 停用' AFTER `deleted`,
ALTER TABLE messages ADD COLUMN `is_top` tinyint(3) NULL DEFAULT 1 COMMENT '1 不置顶 2置顶' AFTER `status`;
ALTER TABLE messages CHANGE COLUMN `status` `state` tinyint(4) NULL DEFAULT 1 COMMENT '1 启用 2 停用' AFTER `deleted`;