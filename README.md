# 品速系列工具

该工具使用go实现 目的是为了方便一些群友 不会操作 想傻瓜式一键刷入 写的

该工具功能有：

1：一键刷入dropbear（轻量型ssh）默认账号为：root 默认密码为online 同时默认开启了sftp服务

2：一键刷入forward（自己写的短信转发程序）github地址：[mixingh/sms_push: 品速R200 CPE 短信推送 (github.com)](https://github.com/mixingh/sms_push/tree/master)推荐使用c语言版本的分支 该工具内置的也是c语言版本的

3：一键刷入IPv6穿透工具（lucky，ddns-go）两个都可以 看自己需求 lucky功能更多一些

4：一键刷入IPv4穿透工具（nps）nps和frp 我选择nps！



声明：
1：dropbear是我从官网下载源码 移植编译armV7架构的 包括upx curl bash 我都有移植 有需要我再上传

2：该工具会关闭CPE的SELINUX模式，介意勿使用

3：如果发现未能成功转发 那应该是您cpe的证书过期了 前上面的forward github地址下载3.3m那个 替换到cpe的/usr/bin目录中

4：关于IPv6穿透工具 使用lucky/ddns-go的前提是您的IPv6入栈了；检查方法：[IPV6版_在线tcping_tcp延迟测试_持续ping_禁ping_tcping_端口延迟测试 (itdog.cn)](https://www.itdog.cn/tcping_ipv6)打开这个网站 输入您cpe的ipv6地址 测试一下是否通了（地图绿色）

5：IPv4穿透工具：使用该工具有个前提：您有服务器并且搭建了nps客户端
