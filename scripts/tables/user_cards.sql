alter table user_cards drop column bank_name;
alter table user_cards drop column bank_code;
alter table user_cards add column bank_id int unsigned not null default 0 comment '银行编号';
alter table user_cards change username user_name varchar(30) not null default '' comment '用户名称';
alter table user_cards change bank_realname real_name varchar(30) not null default '' comment '真实姓名';
alter table user_cards change bank_card card_number varchar(50) not null default '' comment '银行卡号';
alter table user_cards change bank_branch_name branch_name varchar(50) not null default '' comment '支行名称';
alter table user_cards change bank_address address varchar(100) not null default '' comment '地址';

alter table user_cards add column province_id int unsigned not null default 0 comment '省份编号';
alter table user_cards add column city_id int unsigned not null default 0 comment '城市编号';
alter table user_cards add column district_id int unsigned not null default 0 comment '县区编号';


create or replace view user_cards_v as 
    select 
        c.id, c.user_id, c.user_name, c.real_name, c.card_number, c.address, c.branch_name, c.created, c.updated, c.bank_id, 
        b.name as bank_name, b.code as bank_code,
        c.province_id, c.city_id, c.district_id,
        pr.name as province_name, ci.name as city_name, di.name as district_name
    from 
        user_cards as c 
        inner join banks as b on c.bank_id = b.id 
        inner join provinces as pr on pr.id = c.province_id
        inner join cities as ci on ci.id = c.city_id
        inner join districts as di on di.id = c.district_id;
