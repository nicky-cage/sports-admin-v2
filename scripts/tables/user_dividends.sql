ALTER TABLE `user_dividends`
ADD COLUMN `money_type` tinyint(3) NULL DEFAULT 1 COMMENT '钱包类型 1.中心钱包 2.场馆钱包' AFTER `updated`,
ADD COLUMN `operation_type` tinyint(3) NULL DEFAULT 1 COMMENT '操作类型  1.批量发放 2.单会员发放' AFTER `money_type`,
ADD COLUMN `flow_limit` tinyint(3) NULL DEFAULT 1 COMMENT '1 无流水限制 2需要流水限制' AFTER `operation_type`,
ADD COLUMN `flow_multiple` tinyint(3) UNSIGNED NULL DEFAULT 0 COMMENT '流水倍数' AFTER `flow_limit`,
ADD COLUMN `applicant_remark` varchar(300) NULL COMMENT '申请备注' AFTER `flow_multiple`,
ADD COLUMN `applicant` varchar(32) NULL COMMENT '申请人' AFTER `applicant_remark`,
ADD COLUMN `reviewer` varchar(32) NULL COMMENT '审核人' AFTER `applicant`,
ADD COLUMN `reviewer_remark` varchar(300) NULL COMMENT '审核备注' AFTER `reviewer`,
ADD COLUMN `state` tinyint(3) NULL DEFAULT 1 COMMENT '状态 1申请中 2审核通过 3审核拒绝' AFTER `reviewer_remark`;