
backups是一个文件备份工具，支持将本地目录或文件备份至七牛云等云存储。  
backups is a file backup tool, support backup local directory or file to cloud storage.

# 安装/INSTALLATION
在[Releases](https://github.com/qf0129/backups/releases)中下载二进制文件  

# 使用/USAGE

```
backups -c conf.json
```

### 选项/OPTIONS

```bash
-c, 指定json格式的配置文件路径
-v, 查看版本
```

### 配置文件/CONFIGURATION
默认会读取conf.json文件，你也可以通过-c选项指定配置文件路径

```
{
    // 需要备份的目录或文件
    "Paths": [ 
        "/tmp",
        "/var/lib/mysql/",
        "/etc/nginx/conf.d/test.conf"
    ],
    
    // 忽略的目录或文件, 支持通配符
    "IgnorePaths": [
        "**/ib_logfile*", 
        ".git/**",
        "node_modules/**",
        ".DS_Store",
        "src/**/*.log",
        "**/*.pyc"
    ],

    // 七牛云配置
    "Qiniu": {
        "AccessKey": "xxxxx", // 七牛云AccessKey
        "SecretKey": "xxxxx", // 七牛云SecretKey
        "Bucket": "xxx",      // 七牛云Bucket
        "BucketDir": "xxx"    // 七牛云存储路径前缀
    },

    // 是否按天轮转, 默认false，为true时存储路径会增加/YYYYMMDD/目录
    "RotateByDay": true,
}
```

### 定时执行/CRON
可以搭配linux的crontab来定时执行
```
# 每天1点执行
0 1 * * * /opt/backups/backups -c /opt/backups/conf.json
```

# 功能/FEATURE
- [x] 支持忽略指定通配符路径
- [x] 支持按日期存储
- [x] 支持zip打包后再上传
- [ ] 支持腾讯云COS
- [ ] 支持阿里云OSS
- [ ] 支持cloudflare R2
- [ ] 支持AWS S3
- [ ] 支持windows\mac

