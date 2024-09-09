ALTER TABLE game_electrons
ADD COLUMN platform varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '平台 web h5 app' AFTER en_name;
ADD COLUMN games_id int(10) NULL COMMENT '场馆ID' AFTER game_type;
ADD COLUMN `display_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '展示类型 正常 热门 最新' AFTER `updated`,
ADD COLUMN `img_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '图片url' AFTER `display_type`;