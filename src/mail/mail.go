package mail

import (
    "github.com/jordan-wright/email"
    "strings"
    "net/smtp"
    "fmt"
    "github.com/logrusorgru/aurora"
)

const (
    SmtpAddress = ""
    AuthUsername = ""
    AuthPassword = ""
    AuthHost = ""
)


func SendTextMail(mailInfo map[string]string) {
    e := email.NewEmail()

    if from, ok := mailInfo["from"]; ok {
       e.From = from
    }
    if to, ok := mailInfo["to"]; ok {
        e.To = strings.Fields(to)
    }
    if cc, ok := mailInfo["cc"]; ok {
        e.Cc = strings.Fields(cc)
    }
    if subject, ok := mailInfo["subject"]; ok {
        e.Subject = subject
    }
    if text, ok := mailInfo["text"]; ok {
        e.Text = []byte(text)
    }
    if attach, ok := mailInfo["attachment"]; ok {
        e.AttachFile(attach)
    }

    fmt.Println(aurora.Blue("[GoMail] mail sending ..."))
    if err := e.Send(SmtpAddress, smtp.PlainAuth("", AuthUsername, AuthPassword, AuthHost)); err != nil {
        fmt.Println(aurora.Red("[GoMail] send mail failed: "), err)
    } else {
        fmt.Println(aurora.Green("[GoMail] send mail success"))
    }
}