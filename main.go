package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/joho/godotenv"
	"github.com/skip2/go-qrcode"
)

// 跪拜赛博神佛事件
func Pray() {
	Announcement := `
   
	   gocq-wechat client，Powered by barryblueice。
   
	   本项目遵循GPLv3协议开源，不鼓励、不支持一切商业使用。
	   该项目仅供学习交流，严禁用于商业用途，请于24小时内删除。
	   `
	_pray := `

	                        _oo0oo_
	                       o8888888o
	                       88" . "88
	                       (| -_- |)
	                       0\  =  /0
	                     ___/'---'\___
	                   .' \\|     |// '.
	                  / \\|||  :  |||// \
	                 / _||||| -:- |||||- \
	                |   | \\\  - /// |   |
	                | \_|  ''\---/''  |_/ |
	                \  .-\__  '-'  ___/-. /
	              ___'. .'  /--.--\  '. .'___
	           ."" '<  '.___\_<|>_/___.' >' "".
	          | | :  '- \'.;'\ _ /';.'/ - ' : | |
	          \  \ '_.   \_ __\ /__ _/   .-' /  /
	      ====='-.____'.___ \_____/___.-'___.-'=====
	                        '=---='
	
	
	      ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	
	            佛祖保佑     永不宕机     永无BUG
`
	log.Println(Announcement)
	time.Sleep(1 * time.Second)
	log.Println(_pray)
	time.Sleep(1 * time.Second)
}

// 控制台二维码事件
func ConsoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}

// 登录处理事件
func login_requests(bot *openwechat.Bot, reloadStorage openwechat.HotReloadStorage) error {
	var err error // 在函数内部定义err变量
	if _, err = os.Stat("login_token.json"); err == nil {
		log.Println("已检测到鉴权文件，client将根据该鉴权文件执行登录事件……")
		err = bot.HotLogin(reloadStorage)
		if err != nil {
			log.Println("登录失败！鉴权文件可能已过期。正在启动扫码登录事件……")
			time.Sleep(1 * time.Second)
			log.Println("请扫描二维码登录：")
			err = bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption())
			if err != nil {
				log.Printf("登录失败！错误：%s！", err)
				os.Exit(1)
			}
		} else {
			log.Println("登录成功")
			err = nil
		}
	} else if os.IsNotExist(err) {
		log.Println("登录失败！未检测到鉴权文件！正在启动扫码登录事件……")
		time.Sleep(1 * time.Second)
		log.Println("请扫描二维码登录：")
		err = bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption())
		if err != nil {
			log.Printf("登录失败！错误：%s！", err)
			os.Exit(1)
		}
	}
	return err
}

func isEmpty(s string) bool {
	return len(s) == 0
}

// 环境参数处理事件
func EnvHandle() (string, string, string) {
	if _, err := os.Stat(".env"); err != nil {
		// envContent := []byte("# LOGIN_MODE登录模式：确认登录模式\n# 默认可选两种登录模式：1.WEB（网站模式）；2.DESKTOP（桌面模式）\nLOGIN_MODE = \"\"\n\n# PORT：登录端口\nPORT = \n\n# 反向连接地址（暂时仅支持反向连接）\nURL = \"\"")
		envContent := []byte(`
# LOGIN_MODE登录模式：确认登录模式
# 默认可选两种登录模式：1.WEB（网站模式）；2.DESKTOP（桌面模式）
LOGIN_MODE = ""

# PORT：登录端口
PORT = 

# 反向连接地址（暂时仅支持反向连接）
URL = ""
	`)
		err := os.WriteFile(".env", envContent, 0644)
		log.Println("环境文件已成功创建并写入默认值，请修改环境文件默认值后重新启动client")
		if err != nil {
			log.Fatalf("无法创建环境文件文件: %v", err)
			os.Exit(1)
		}
		os.Exit(1)
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("无法加载环境文件文件: %v", err)
		os.Exit(1)
	}

	LOGIN_MODE := os.Getenv("LOGIN_MODE")
	PORT := os.Getenv("PORT")
	URL := os.Getenv("URL")

	if isEmpty(strings.TrimSpace(LOGIN_MODE)) {
		log.Fatalf("LOGIN_MODE不能为空，请检查参数！")
		os.Exit(1)
	} else if isEmpty(strings.TrimSpace(PORT)) {
		log.Fatalf("PORT不能为空，请检查参数！")
		os.Exit(1)
	} else if isEmpty(strings.TrimSpace(URL)) {
		log.Fatalf("URL不能为空，请检查参数！")
		os.Exit(1)
	}

	return LOGIN_MODE, PORT, URL
}

func main() {
	Pray()
	LOGIN_MODE, PORT, URL := EnvHandle()
	log.Println("开始执行登录事件")

	// 登录情况判断：网页 & 桌面端

	var bot *openwechat.Bot

	if LOGIN_MODE == "DESKTOP" {
		bot = openwechat.DefaultBot(openwechat.Desktop)
	} else if LOGIN_MODE == "WEB" {
		bot = openwechat.DefaultBot()
	} else {
		log.Fatalf("LOGIN_MODE参数不合法，请检查参数！")
		os.Exit(1)
	}

	log.Printf("登录模式：%s", LOGIN_MODE)

	// 端口判断：只是判断端口在不在可用范围内

	TargetPort, err := strconv.Atoi(PORT)

	if err != nil {
		log.Fatal("端口参数不合法，请检查参数！")
		os.Exit(1)
	}
	if TargetPort < 0 && TargetPort > 65535 {
		log.Fatal("端口参数不合法，请检查参数！")
		os.Exit(1)
	}

	log.Printf("端口号：%s", PORT)

	// 连接事件预留，这里暂且只留反向连接

	log.Printf("连接地址：%s", URL)

	// 登录事件处理

	reloadStorage := openwechat.NewFileHotReloadStorage("login_token.json")
	bot.UUIDCallback = ConsoleQrCode
	defer reloadStorage.Close()
	err = login_requests(bot, reloadStorage)
	if err != nil {
		log.Printf("登录失败！错误：%s！", err)
		os.Exit(1)
	}

	self, err := bot.GetCurrentUser()
	if err != nil {
		log.Printf("登录信息获取失败！错误：%s！", err)
		os.Exit(1)
	}

	// 登录成功后

	// SelfID := self.ID()
	SelfID := strconv.Itoa(int(self.ID()))
	SelfNickname := self.NickName
	log.Printf("欢迎回来！用户：%s，wxid：%s", SelfNickname, SelfID)
	go WebsocketReverseInit(URL, SelfID, self)
	// go FriendListInit(self)
	bot.MessageHandler = func(msg *openwechat.Message) {
		if isEmpty(msg.Content) != true {
			if msg.IsSendByFriend() {
				sender, _ := msg.Sender()
				log.Printf("用户 %s（%s）发送私聊消息：%s", sender.NickName, sender.ID(), msg.Content)
				if conn != nil {
					RecievePrivateText(self, 0, sender.ID(), msg.Content)
				}
			} else if msg.IsSendByGroup() {
				group, _ := msg.Sender()
				sender, _ := msg.SenderInGroup()
				log.Printf("用户 %s 于群聊 %s（%s）发送消息：%s", sender.NickName, group.NickName, group.ID(), msg.Content)
				if conn != nil {
					RecieveGroupText(self, 0, sender.ID(), msg.Content, group.ID())
				}
			}
		}
	}
	bot.Block()
}
