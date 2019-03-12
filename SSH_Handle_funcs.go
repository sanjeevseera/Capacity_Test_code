package main

import(
	"io/ioutil"
	"log"
	"strings"
//	"fmt"
//        "github.com/gliderlabs/ssh"
//	"io"
//	confdssh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"path/filepath"
	"github.com/Unknwon/com"
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/setting"
	"runtime/debug"
)

// Listen starts a SSH server listens on given port.
func Listen(addr string) {
	ciphers := []string{"aes128-ctr","aes192-ctr","aes256-ctr","arcfour256","arcfour128","aes128-cbc","3des-cbc"}
//	keyExchanges := []string{"diffie-hellman-group-exchange-sha1","diffie-hellman-group14-sha1","diffie-hellman-group1-sha1"}
	keyExchanges := []string{"diffie-hellman-group-exchange-sha256","diffie-hellman-group14-sha1","diffie-hellman-group1-sha1"}
	macs := []string{"hmac-md5","hmac-sha1","umac-64@openssh.com","hmac-ripemd160"}

	config := &ssh.ServerConfig{
		Config: ssh.Config{
			Ciphers:      ciphers,
			KeyExchanges: keyExchanges,
			MACs:         macs,
		},
		NoClientAuth: true,
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			pkey, err := models.SearchPublicKeyByContent(strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key))))
			if err != nil {
				log.Printf(" *ERR* %v SearchPublicKeyByContent: %v",strings.Split(addr,":")[0],err)
				return nil, err
			}
			return &ssh.Permissions{Extensions: map[string]string{"key-id": com.ToStr(pkey.ID)}}, nil
		},
	}

	keyPath := filepath.Join(setting.AppDataPath, "ssh/gogs.rsa")
	if !com.IsExist(keyPath) {
		os.MkdirAll(filepath.Dir(keyPath), os.ModePerm)
		_, stderr, err := com.ExecCmd("ssh-keygen", "-f", keyPath, "-t", "rsa", "-N", "")
		if err != nil {
			log.Printf(" *ERR* %v Fail to generate private key: %v - %s", strings.Split(addr,":")[0],err, stderr)
		}
		log.Printf(" *INFO* %v New private key is generateed: %v", strings.Split(addr,":")[0],keyPath)
	}

	privateBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Printf(" *ERR* %v Fail to load privateByte key",strings.Split(addr,":")[0])
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Printf(" *ERR* %v Fail to parse private key",strings.Split(addr,":")[0])
	}
	config.AddHostKey(private)
	//wg_listen.Add(1)
	go listen(config, addr)
}

func listen(config *ssh.ServerConfig, addr string) {
	log.Printf(" *INFO* %v listen function",strings.Split(addr,":")[0])
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf(" *ERR* %v err:[%v]",strings.Split(addr,":")[0],err)
		index := indexOf(strings.Split(addr,":")[0],vnfips)
		VNFStates[index].NewEvent = "Exit"
		VNFStates[index].spawnStat = "Spawn Failed"
	}else{
		// Once a ServerConfig has been configured, connections can be accepted.
		for true{
			conn, err := listener.Accept()
			if err != nil {
				log.Printf(" *ERR* %v listener.Accept err:[%v]",strings.Split(addr,":")[0],err)
			}
			// Before use, a handshake must be performed on the incoming net.Conn.
			sConn, chans, reqs, err := ssh.NewServerConn(conn, config)
			if err != nil {
				log.Printf(" *ERR* %v ssh.NewServerConn sConn:[%v] err:[%v] ",strings.Split(addr,":")[0],sConn,err)
			}else{
				// The incoming Request channel must be serviced.
				log.Printf(" *INFO* %v ssh.NewServerConn sConn:[%v] reqs:[%v] ",strings.Split(addr,":")[0],sConn,reqs)
				//go ssh.DiscardRequests(reqs)
				//log.Printf(" *ERR* %v sConn.Permissions.Extensions:[%v]",strings.Split(addr,":")[0],sConn.Permissions.Extensions["key-id"])
				//go handleServerConn(sConn.Permissions.Extensions["key-id"], chans)
				handleServerConn(chans,addr)
                                log.Printf(" *INFO* %v Handled ssh.NewServerConn sConn:[%v] reqs:[%v] ",strings.Split(addr,":")[0],sConn,reqs)
				//wg_Listen.Done()
				//wg_listen.Done()
				debug.FreeOSMemory()
				break  //changed now sanjeev
			}
		}
	}
}
//func handleServerConn(keyID string, chans <-chan ssh.NewChannel) {
func handleServerConn(chans <-chan ssh.NewChannel,addr string) {
	//log.Printf(" *INFO* %v HandleServerConn",strings.Split(addr,":")[0])
	index := indexOf(strings.Split(addr,":")[0],vnfips)
	if VNFStates[index].GetState() == "InitSuccess"{
		VNFStates[index].NewEvent = "Send"
		VNFStates[index].pos += 1
	}
	for newChan := range chans {
		log.Printf(" *INFO* %v newChan loop",strings.Split(addr,":")[0])
		if newChan.ChannelType() != "session" {
			log.Printf(" *ERR* %v UnknownChannelType",strings.Split(addr,":")[0])
			newChan.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		ch, reqs, err := newChan.Accept()
		if err != nil {
			log.Printf(" *ERR* %v newChan.Accept err:%v",strings.Split(addr,":")[0],err)
			continue
		}
		log.Printf(" *ERR* %v newChan.Accept ch:[%v] reqs:[%v]",strings.Split(addr,":")[0],ch, reqs)

//		go func(in <-chan *ssh.Request) {
		func(in <-chan *ssh.Request) {
			defer ch.Close()
			for req := range in {
				payload := cleanCommand(string(req.Payload))
				log.Printf(" *INFO* %v payload: ",strings.Split(addr,":")[0],payload)
			}
		}(reqs)
	}
}
func cleanCommand(cmd string) string {
	i := strings.Index(cmd, "git")
	if i == -1 {
		return cmd
	}
	return cmd[i:]
}
