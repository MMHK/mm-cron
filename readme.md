# MM Cron

如果不是Linux 自带的 `conrtab` 不支持 秒级的调用，也不用重复造这轮子。

## 项目依赖

- [jobrunner](https://github.com/bamzi/jobrunner)， 感谢这位兄台完成了大部分的工作，基于他的工作基础，只需要略略的写下wrapper就好了。

## 编译项目

- 安装 ``Golang`` 环境， Golang >= 1.1
- 进入源码目录，运行
    ```
    go get -v
    ```
- 编译源码
    ```
    go build
    ```
    
## 使用

主要的配置项在 `conf.json` 里面，`task` 节点包含一堆待执行的命令。

`time` 这个是根据`crontab` 一样的定义：`<秒> <分钟> <小时> <日> <月份> <星期>`

`cmd` 就是具体需要执行的命令

```
{
	"task" : [
		{"time": "*/15 * * * * *", "cmd": "wget --spider http://demo2.mixmedia.com/freelink/driver_xmas/webhook/upload/insurance"},
		{"time": "*/15 * * * * *", "cmd": "wget --spider http://ins.demo2.mixmedia.com/webhook/queue"}
	]
}
```

或者直接使用Web UI

http://127.0.0.1:5566/ui/


## 生成 `swagger` 文档

- 安装 [swagger-go](https://github.com/go-swagger/go-swagger)
- 在项目目录执行
```bash
swagger generate spec -o ./web_root/swagger/swagger.json
