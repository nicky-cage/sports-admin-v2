CREATE TABLE `finance_logs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `bill_no` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '订单号',
  `type` tinyint(3) NOT NULL DEFAULT '0' COMMENT '0 存款 1 提款 2 代理提款',
  `operating` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '操作',
  `result` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '结果',
  `consuming` varchar(12) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '耗时,单位秒',
  `operator` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '操作人',
  `remark` varchar(600) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '备注',
  `created` int(10) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;