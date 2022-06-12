实习期间写的一个考核项目，整体目标是实现一个网站的后端服务，实现用户、顾问注册，创建、回复订单、订单评论等接口。

## 相关技术

- [`gin`](https://github.com/gin-gonic/gin) 
- [`validator`](https://github.com/go-playground/validator)数据校验
- 关系型数据库`MySQL`并使用[`gendry`](https://github.com/didi/gendry)动态构建sql语句
- `Jwt`鉴权，并构建了token和role两个中间件
- `redis`缓存订单数据和评论数据
- [`cron`](https://github.com/robfig/cron)定时处理用户的订单状态（加急订单1h恢复普通状态，普通订单24h过期）
- [`zap`](https://github.com/uber-go/zap) 日志管理

## 参考链接

- [Uber Go 语言编码规范](https://github.com/xxjwxc/uber_go_guide_cn)
- [在 Golang 中用名字调用函数](https://studygolang.com/articles/1035)
- [Cron定时任务](https://juejin.cn/post/7066479548115714084)
- [gin框架参数零值json绑定的问题](https://blog.csdn.net/weixin_42279809/article/details/107800081)

