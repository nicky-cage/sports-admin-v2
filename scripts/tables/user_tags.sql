-- 标签分类
drop table if exists user_tag_categories;
create table if not exists user_tag_categories (
    id int unsigned not null auto_increment primary key,
    name varchar(50) not null default '' comment '名称',
    remark varchar(200) not null default '' comment '备注'
);

-- 标签
drop table if exists user_tags;
create table if not exists user_tags (
    id int unsigned not null auto_increment primary key,
    category_id int unsigned not null default 0 comment '分类编号',
    name varchar(50) not null default '' comment '名称',
    remark varchar(200) not null default '' comment '备注',
    index(category_id)
);

-- 插入默认数据
INSERT INTO user_tag_categories (name) VALUES
('风控管理');
SELECT last_insert_id() INTO @last_id;
INSERT INTO user_tags (category_id, name) VALUES
(@last_id, '风控管理'),
(@last_id, '白名单'),
(@last_id, '限制取款'),
(@last_id, '财务抽查'),
(@last_id, '人工出款');

INSERT INTO user_tag_categories (name) VALUES
('财务管理');
SELECT last_insert_id() INTO @last_id;
INSERT INTO user_tags (category_id, name) VALUES
(@last_id, '风控管理'),
(@last_id, '白名单'),
(@last_id, '限制取款'),
(@last_id, '财务抽查'),
(@last_id, '人工出款');

INSERT INTO user_tag_categories (name) VALUES
('权限控制');
SELECT last_insert_id() INTO @last_id;
INSERT INTO user_tags (category_id, name) VALUES
(@last_id, '风控管理'),
(@last_id, '白名单'),
(@last_id, '限制取款'),
(@last_id, '财务抽查'),
(@last_id, '人工出款');

-- 会员标签信息
CREATE OR REPLACE VIEW user_tags_v AS
    SELECT c.id AS category_id, c.name AS category_name, t.id AS id, t.name AS name, t.remark AS remark
    FROM user_tag_categories AS c INNER JOIN user_tags AS t ON c.id = t.category_id;