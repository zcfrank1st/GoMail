package conf

import (
    "io/ioutil"
    "encoding/json"
    "fmt"
    "github.com/logrusorgru/aurora"
)

type MailConf struct {
    Host string `json:"host"`
    Port string    `json:"port"`
    AuthUsername string `json:"username"`
    AuthPassword string `json:"password"`
}

var MailConfInstance *MailConf

func LoadConf(filepath string) {
    if cb, err := ioutil.ReadFile(filepath); err == nil {
        fmt.Println(aurora.Blue("[GoMail] load mail conf success"))
        json.Unmarshal(cb, &MailConfInstance)
        fmt.Println(aurora.Green(MailConfInstance))
    } else {
        fmt.Println(aurora.Red("[GoMail] load mail conf failed"))
    }
}