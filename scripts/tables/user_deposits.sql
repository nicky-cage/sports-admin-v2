ALTER TABLE `user_deposits`
ADD COLUMN `arrive_money` decimal(16, 2) NULL COMMENT '到账金额' AFTER `after_money`,
ADD COLUMN `confirm_money` decimal(16, 2) NULL COMMENT '确认金额' AFTER `arrive_money`;
ADD COLUMN `top_money` decimal(16, 2) NULL COMMENT '上分金额' AFTER `arrive_money`,
ADD COLUMN `discount` decimal(16, 2) NULL COMMENT '存款优惠' AFTER `top_money`,
ADD COLUMN `business_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '商户名称' AFTER `remark`,
ADD COLUMN `card_number` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '收款卡号' AFTER `business_name`;