package main

import (
    "github.com/c-bata/go-prompt"
    "fmt"
    "github.com/logrusorgru/aurora"
    "strings"
    "os"
    "mail"
    "conf"
)

const (
    Logo = `
 ____                                    ___
/\  _ \           /'\_/ \            __ /\_ \
\ \ \L\_\    ___ /\      \     __   /\_\\//\ \
 \ \ \L_L   / __ \ \ \__\ \  /'__ \ \/\ \ \ \ \
  \ \ \/, \/\ \L\ \ \ \_/\ \/\ \L\.\_\ \ \ \_\ \_
   \ \____/\ \____/\ \_\\ \_\ \__/.\_\\ \_\/\____\
    \/___/  \/___/  \/_/ \/_/\/__/\/_/ \/_/\/____/
                                             v1.0`
)

func completer(d prompt.Document) []prompt.Suggest {
    s := []prompt.Suggest{
        {Text: "load", Description: "load /usr/local/mailconf.json"},
        {Text: "from", Description: "from xx"},
        {Text: "to", Description: "to xx xx"},
        {Text: "cc", Description: "cc xx xx"},
        {Text: "subject", Description: "subject this is a demo"},
        {Text: "text", Description: "text it is a demo"},
        {Text: "attach", Description: "attach /usr/local/demo.txt"},
        {Text: "quit", Description: "quit the shell"},
        {Text: "", Description: "[Tips] more than one keyword use ; to split"},
    }
    return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
    var mailInfo = make(map[string]string)

    fmt.Println(aurora.Blue(Logo))

    cmd:
    for {
        t := prompt.Input("[GoMail] >>> ", completer, prompt.OptionPrefixTextColor(prompt.Brown), prompt.OptionMaxSuggestion(9))

        if "" != t {

            segments := strings.Split(t, ";")
            for _, commandline := range segments {
                command := strings.Fields(strings.Trim(commandline, " "))
                verb := command[0]
                args := command[1:]

                args_length := len(args)
                switch verb {
                case "from", "load":
                    if 1 != args_length {
                        fmt.Println(aurora.Red("[GoMail] command must have only one argument"))
                        goto cmd
                    }
                    if "load" == verb {
                        conf.LoadConf(args[0])
                        goto cmd
                    } else {
                        mailInfo[verb] = args[0]
                    }
                case "to", "cc", "subject", "text", "attach":
                    if 0 == args_length {
                        fmt.Println(aurora.Red("[GoMail] command must at least have one argument"))
                        goto cmd
                    }
                    mailInfo[verb] = strings.Join(args, " ")
                case "quit":
                    fmt.Println(aurora.Red("[GoMail] quiting now ..."))
                    os.Exit(1)
                default:
                    fmt.Println(aurora.Red("[GoMail] not support command"))
                }

            }
            if conf.MailConfInstance == nil {
                fmt.Println(aurora.Red("[GoMail] not yet load conf success"))
                continue
            }
            mail.SendTextMail(mailInfo)

        } else {
            fmt.Println(aurora.Red("[GoMail] input is none"))
        }
    }
}
