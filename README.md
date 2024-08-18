# 由于网页端wechat客户端功能不完全，现已放弃该项目，转而使用通过[hook客户端方式的onebot v12实现方案](https://github.com/barryblueice/Vanilla-Client)。

# 新坑地址：[Vanilla-Client：又一个微信端的onebot v12实现](https://github.com/barryblueice/Vanilla-Client)

## 简介

基于openwechat的onebot v11客户端实现

## 许可证

采用 [GPLv3](https://github.com/barryblueice/gocq-wechat/blob/main/LICENSE) 协议开源，不鼓励、不支持一切商业使用。

## 声明

该项目仅供学习交流，严禁用于商业用途，请于24小时内删除。</br>
本人目前仍是go新手小白，还请各位大佬轻喷~ >ᯅ<

## 上游依赖

- [openwechat](https://github.com/eatmoreapple/openwechat)：golang版个人微信号API, 突破登录限制，类似开发公众号一样，开发个人微信号

## Onebot11支持

- [ ] HTTP
- [ ] HTTP Webhook
- [ ] 正向 Websocket
- [x] 反向 Websocket

## 目前已实现功能：

- [x] 群聊/私聊信息发送文本
- [x] 群聊/私聊信息发送图片
- [x] 群聊检测被at
