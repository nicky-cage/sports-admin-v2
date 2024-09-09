ALTER TABLE `user_withdraws`
ADD COLUMN `business_type` tinyint(3) NULL COMMENT '商户类型 1天下 2风云 3.易付宝' AFTER `updated`,
ADD COLUMN `card_number` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '出款卡号' AFTER `business_type`;
ADD COLUMN `cause_failure` tinyint(255) NULL COMMENT '失败原因 例如 1 打码量不足 2 违规操作' AFTER `card_number`,
ADD COLUMN `failure_reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '失败原因内容' AFTER `cause_failure`;
ADD COLUMN `agent_admin` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '代理审核人' AFTER `failure_reason`;