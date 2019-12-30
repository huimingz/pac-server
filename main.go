package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
)

const PAC_URL = "https://raw.githubusercontent.com/petronny/gfwlist2pac/master/gfwlist.pac"
var pacFilePath string

func main() {

    var flags struct {
        Host      string
        Port      int
        Mode      string
        UpdatePAC bool
        Verbose   bool
        ProxyHost string
        ProxyPort int
        pacFile string
    }

    flag.StringVar(&flags.Host, "host", "127.0.0.1", "pac server host")
    flag.IntVar(&flags.Port, "port", 10010, "pac server port")
    flag.StringVar(&flags.ProxyHost, "proxy_host", "127.0.0.1", "proxy server host")
    flag.IntVar(&flags.ProxyPort, "proxy_port", 1080, "proxy server port")
    flag.BoolVar(&flags.UpdatePAC, "update", true, "update pac file from remote")
    flag.BoolVar(&flags.Verbose, "verbose", false, "")
    flag.StringVar(&flags.pacFile, "pac_file", "gfwlist.pac", "pac file")
    flag.Parse()

    pacFilePath = flags.pacFile

    if flags.Verbose {
        log.SetLevel(log.DebugLevel)
        log.Info("设置日志等级为DEBUG级别")
    }

    if flags.UpdatePAC {
        log.Info("更新pac文件...")
        err := updatePacFile()
        if err != nil {
            log.Errorf("更新pac文件失败，错误信息: %s", err.Error())
            os.Exit(1)
        } else {
            log.Infof("更新pac文件成功，存储位置: %s", pacFilePath)
        }
    }

    log.Info("读取pac文件信息...")
    data, err := getPACData()
    if err != nil {
        log.Error("读取pac文件内容失败，错误信息: %s", err.Error())
        os.Exit(1)
    } else {
        log.Info("读取pac文件信息成功")
    }


    proxyAddr := fmt.Sprintf("%s:%d", flags.ProxyHost, flags.ProxyPort)
    data = strings.Replace(data, "127.0.0.1:1080", proxyAddr, 1)

    log.Infof("PAC URL: http://%s:%d/pac", host, port)
    service(flags.Host, flags.Port, data)
}

func service(host string, port int, data string) {
    router := gin.Default()

    router.GET("/pac", func(ctx *gin.Context) {
        ctx.Header("Content-Type", "application/x-ns-proxy-autoconfig")
        ctx.String(http.StatusOK, data)
    })

    address := fmt.Sprintf("%s:%d", host, port)


    log.Fatal(router.Run(address))
}

// 从远程更新pac文件至本地
func updatePacFile() error {
    resp, err := http.Get(PAC_URL)
    if err != nil {
        return err
    }

    defer resp.Body.Close()

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    return savePac2File(data)
}

// 从本地文件读取pac文件
func getPACData() (data string, err error) {
    file, err := os.Open(pacFilePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    byteData, err := ioutil.ReadAll(file)
    if err !=nil {
        return "", err
    }

    return string(byteData), err
}

// 保存pac信息至本地文件
func savePac2File(data []byte) error {
    file, err := os.OpenFile(pacFilePath, os.O_TRUNC | os.O_CREATE | os.O_WRONLY, 0766)
    if err != nil {
        return err
    }

    defer file.Close()

    _, err = file.Write(data)
    if err != nil {
        return err
    }

    return nil
}

func init() {
    log.SetFormatter(&log.TextFormatter{})

    log.SetOutput(os.Stdout)

    log.SetLevel(log.InfoLevel)
}
