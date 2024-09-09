CREATE TABLE `user_account_sets` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `bill_no` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '订单号',
  `user_id` int(10) unsigned DEFAULT NULL COMMENT '会员id',
  `username` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '用户名',
  `abnormal_bill_no` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '异常订单号',
  `type` tinyint(3) DEFAULT '1' COMMENT '1 上分 2 下分',
  `money` decimal(16,2) DEFAULT NULL COMMENT '钱',
  `reason` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '原因',
  `applicant_remark` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '申请备注',
  `applicant` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '申请人',
  `audit_remark` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '审核备注',
  `audit` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '审核人',
  `status` tinyint(255) DEFAULT '1' COMMENT '1未处理 2成功 3失败',
  `created` int(10) unsigned DEFAULT NULL COMMENT '申请时间',
  `updated` int(10) unsigned DEFAULT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;