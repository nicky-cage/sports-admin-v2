/** 银行信息表 **/
drop table if exists banks;
create table if not exists banks (
    id int unsigned not null auto_increment,
    name varchar(50) not null default '' comment '银行名称',
    code varchar(20) not null default '' comment '编码',
    remark varchar(50) not null default '' comment '备注',
    status tinyint not null default 1 comment '状态',
    sort int not null default 0 comment '排序',
    primary key(id),
    unique key(name),
    unique key(code)
);

INSERT INTO banks (name, code, remark, status, sort) VALUES ('国家开发银行','CDB','国家开发银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('中国工商银行','ICBC','中国工商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('中国农业银行','ABC','中国农业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('中国银行','BOC','中国银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('中国建设银行','CCB','中国建设银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('中国邮政储蓄银行','PSBC','中国邮政储蓄银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('交通银行','COMM','交通银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('招商银行','CMB','招商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('上海浦东发展银行','SPDB','上海浦东发展银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('兴业银行','CIB','兴业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('华夏银行','HXBANK','华夏银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('广东发展银行','GDB','广东发展银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('中国民生银行','CMBC','中国民生银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('中信银行','CITIC','中信银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('中国光大银行','CEB','中国光大银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('恒丰银行','EGBANK','恒丰银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('浙商银行','CZBANK','浙商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('渤海银行','BOHAIB','渤海银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('平安银行','SPABANK','平安银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('上海农村商业银行','SHRCB','上海农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('玉溪市商业银行','YXCCB','玉溪市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('尧都农商行','YDRCB','尧都农商行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('北京银行','BJBANK','北京银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('上海银行','SHBANK','上海银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('江苏银行','JSBANK','江苏银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('杭州银行','HZCB','杭州银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('南京银行','NJCB','南京银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('宁波银行','NBBANK','宁波银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('徽商银行','HSBANK','徽商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('长沙银行','CSCB','长沙银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('成都银行','CDCB','成都银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('重庆银行','CQBANK','重庆银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('大连银行','DLB','大连银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('南昌银行','NCB','南昌银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('福建海峡银行','FJHXBC','福建海峡银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('汉口银行','HKB','汉口银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('温州银行','WZCB','温州银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('青岛银行','QDCCB','青岛银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('台州银行','TZCB','台州银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('嘉兴银行','JXBANK','嘉兴银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('常熟农村商业银行','CSRCB','常熟农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('南海农村信用联社','NHB','南海农村信用联社',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('常州农村信用联社','CZRCB','常州农村信用联社',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('内蒙古银行','H3CB','内蒙古银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('绍兴银行','SXCB','绍兴银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('顺德农商银行','SDEB','顺德农商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('吴江农商银行','WJRCB','吴江农商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('齐商银行','ZBCB','齐商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('贵阳市商业银行','GYCB','贵阳市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('遵义市商业银行','ZYCBANK','遵义市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('湖州市商业银行','HZCCB','湖州市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('龙江银行','DAQINGB','龙江银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('晋城银行JCBANK','JINCHB','晋城银行JCBANK',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('浙江泰隆商业银行','ZJTLCB','浙江泰隆商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('广东省农村信用社联合社','GDRCC','广东省农村信用社联合社',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('东莞农村商业银行','DRCBCL','东莞农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('浙江民泰商业银行','MTBANK','浙江民泰商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('广州银行','GCB','广州银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('辽阳市商业银行','LYCB','辽阳市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('江苏省农村信用联合社','JSRCU','江苏省农村信用联合社',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('廊坊银行','LANGFB','廊坊银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('浙江稠州商业银行','CZCB','浙江稠州商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('德阳商业银行','DYCB','德阳商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('晋中市商业银行','JZBANK','晋中市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('苏州银行','BOSZ','苏州银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('桂林银行','GLBANK','桂林银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('乌鲁木齐市商业银行','URMQCCB','乌鲁木齐市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('成都农商银行','CDRCB','成都农商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('张家港农村商业银行','ZRCBANK','张家港农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('东莞银行','BOD','东莞银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('莱商银行','LSBANK','莱商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('北京农村商业银行','BJRCB','北京农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('天津农商银行','TRCB','天津农商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('上饶银行','SRBANK','上饶银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('富滇银行','FDB','富滇银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('重庆农村商业银行','CRCBANK','重庆农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('鞍山银行','ASCB','鞍山银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('宁夏银行','NXBANK','宁夏银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('河北银行','BHB','河北银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('华融湘江银行','HRXJB','华融湘江银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('自贡市商业银行','ZGCCB','自贡市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('云南省农村信用社','YNRCC','云南省农村信用社',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('吉林银行','JLBANK','吉林银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('东营市商业银行','DYCCB','东营市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('昆仑银行','KLB','昆仑银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('鄂尔多斯银行','ORBANK','鄂尔多斯银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('邢台银行','XTB','邢台银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('晋商银行','JSB','晋商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('天津银行','TCCB','天津银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('营口银行','BOYK','营口银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('吉林农信','JLRCU','吉林农信',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('山东农信','SDRCU','山东农信',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('西安银行','XABANK','西安银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('河北省农村信用社','HBRCU','河北省农村信用社',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('宁夏黄河农村商业银行','NXRCU','宁夏黄河农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('贵州省农村信用社','GZRCU','贵州省农村信用社',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('阜新银行','FXCB','阜新银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('湖北银行黄石分行','HBHSBANK','湖北银行黄石分行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('浙江省农村信用社联合社','ZJNX','浙江省农村信用社联合社',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('新乡银行','XXBANK','新乡银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('湖北银行宜昌分行','HBYCBANK','湖北银行宜昌分行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('乐山市商业银行','LSCCB','乐山市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('江苏太仓农村商业银行','TCRCB','江苏太仓农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('驻马店银行','BZMD','驻马店银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('赣州银行','GZB','赣州银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('无锡农村商业银行','WRCB','无锡农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('广西北部湾银行','BGB','广西北部湾银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('广州农商银行','GRCB','广州农商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('江苏江阴农村商业银行','JRCB','江苏江阴农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('平顶山银行','BOP','平顶山银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('泰安市商业银行','TACCB','泰安市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('南充市商业银行','CGNB','南充市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('重庆三峡银行','CCQTGB','重庆三峡银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('中山小榄村镇银行','XLBANK','中山小榄村镇银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('邯郸银行','HDBANK','邯郸银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('库尔勒市商业银行','KORLABANK','库尔勒市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('锦州银行','BOJZ','锦州银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('齐鲁银行','QLBANK','齐鲁银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('青海银行','BOQH','青海银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('阳泉银行','YQCCB','阳泉银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('盛京银行','SJBANK','盛京银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('抚顺银行','FSCB','抚顺银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('郑州银行','ZZBANK','郑州银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('深圳农村商业银行','SRCB','深圳农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('潍坊银行','BANKWF','潍坊银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('九江银行','JJBANK','九江银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('江西省农村信用','JXRCU','江西省农村信用',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('河南省农村信用','HNRCU','河南省农村信用',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('甘肃省农村信用','GSRCU','甘肃省农村信用',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('四川省农村信用','SCRCU','四川省农村信用',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('广西省农村信用','GXRCU','广西省农村信用',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('陕西信合','SXRCCU','陕西信合',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('武汉农村商业银行','WHRCB','武汉农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('宜宾市商业银行','YBCCB','宜宾市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('昆山农村商业银行','KSRB','昆山农村商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('石嘴山银行','SZSBK','石嘴山银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('衡水银行','HSBK','衡水银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('信阳银行','XYBANK','信阳银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('鄞州银行','NBYZ','鄞州银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('张家口市商业银行','ZJKCCB','张家口市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('许昌银行','XCYH','许昌银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('济宁银行','JNBANK','济宁银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('开封市商业银行','CBKF','开封市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('威海市商业银行','WHCCB','威海市商业银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('湖北银行','HBC','湖北银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('承德银行','BOCD','承德银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('丹东银行','BODD','丹东银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('金华银行','JHBANK','金华银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('朝阳银行','BOCY','朝阳银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('临商银行','LSBC','临商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('包商银行','BSB','包商银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('兰州银行','LZYH','兰州银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('周口银行','BOZK','周口银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('德州银行','DZBANK','德州银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('三门峡银行','SCCB','三门峡银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('安阳银行','AYCB','安阳银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('安徽省农村信用社','ARCU','安徽省农村信用社',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('湖北省农村信用社','HURCB','湖北省农村信用社',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('湖南省农村信用社','HNRCC','湖南省农村信用社',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('广东南粤银行','NYNB','广东南粤银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('洛阳银行','LYBANK','洛阳银行',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('农信银清算中心','NHQS','农信银清算中心',0,0 );
INSERT INTO banks (name, code, remark, status, sort) VALUES ('城市商业银行资金清算中心','CBBQS','城市商业银行资金清算中心',0,0 );

alter table banks change status state tinyint unsigned not null default 0 comment '状态';

alter table banks add icon varchar(200) not null default '' comment '图标';
alter table banks add image varchar(200) not null default '' comment '图片';
