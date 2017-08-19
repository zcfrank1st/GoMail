package mail

import (
    "strings"
    "fmt"
    "github.com/logrusorgru/aurora"
    "gopkg.in/gomail.v2"
    "crypto/tls"
    "conf"
    "strconv"
)

func SendTextMail(mailInfo map[string]string) {
    m := gomail.NewMessage()

    if from, ok := mailInfo["from"]; ok {
        m.SetHeader("From", from)
    }
    if to, ok := mailInfo["to"]; ok {
        m.SetHeader("To", strings.Fields(to)...)
    }
    if cc, ok := mailInfo["cc"]; ok {
        for _, ele := range strings.Fields(cc) {
            m.SetAddressHeader("Cc", ele, "")
        }
    }
    if subject, ok := mailInfo["subject"]; ok {
        m.SetHeader("Subject", subject)
    }
    if text, ok := mailInfo["text"]; ok {
        m.SetBody("text/plain", text)
    }
    if attach, ok := mailInfo["attachment"]; ok {
        m.Attach(attach)
    }

    fmt.Println(aurora.Blue("[GoMail] mail sending ..."))

    port, _ := strconv.Atoi(conf.MailConfInstance.Port)
    d := gomail.NewPlainDialer(conf.MailConfInstance.Host, port, conf.MailConfInstance.AuthUsername, conf.MailConfInstance.AuthPassword)
    d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    if err := d.DialAndSend(m); err != nil {
        fmt.Println(aurora.Red("[GoMail] send mail failed: "), err)
    } else {
        fmt.Println(aurora.Green("[GoMail] send mail success"))
    }
}