# integrated-admin

Shipu 后台管理系统

# 开发说明
## 一、安装依赖
```bash
# 安装gin
go get -u "github.com/gin-gonic/gin"
# 安装pongo2
go get -u "github.com/flosch/pongo2"
go get -u "gitlab.com/go-box/pongo2gin"
# 导入session包
go get -u "github.com/gin-contrib/sessions"
# 导入session存储引擎
go get -u "github.com/gin-contrib/sessions/cookie"
```

## 二、配置网站
```bash
# 将nginx配置文件复制到本地环境
cp nginx.conf.default nginx.conf
# 修改配置文件,将网站目录指定为 {你的开发目录}/pulbic
# 如源代码目录为: /home/fly/integrated-admin
# 则可以配置网站: root /home/fly/integrated-admin/public
vim nginx.conf
```

## 三、运行程序
```bash
# 可依据你本地环境自行修改此启动脚本
# 一般情况下不需要修改可直接启动运行
./start.sh
```

## 四、参考模板
本系统前端使用layuiadmin，相应的模板文件位于: 
```bash
ls -lha ./layuiAdmin.std-v1.4.0
# 已经压缩目录, 一般使用浏览器打开以下文件即可: 
ls -lha ./layuiAdmin.std-v1.4.0/dist/views/index.html
# 未经压缩目录, 一般不需要管此目录
ls -lha ./layuiAdmin.std-v1.4.0/src
```
layui官方网站: https://www.layui.com/doc/

## 五、实时编译、执行
```
# 下载 realize 插件: 
go get github.com/oxequa/realize
# 进入项目目录
cd 项目目录
# 运行启动脚本
./start.sh
```

