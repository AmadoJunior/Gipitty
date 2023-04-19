package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/AmadoJunior/Gipitty/config"
	"github.com/AmadoJunior/Gipitty/models"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// 👇 Email template parser

func ParseTemplateDir(dir string, templateName string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fmt.Println(path)
			if path == "templates/base.html" || path == "templates/styles.html" || path == "templates/"+templateName {
				paths = append(paths, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	fmt.Println(paths)
	return template.ParseFiles(paths...)
}

func SendEmail(user *models.DBResponse, data *EmailData, templateName string) error {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates", templateName)
	if err != nil {
		log.Fatal("Could not parse template", err)
	}
	fmt.Println(template.DefinedTemplates())
	newTemplate := template.Lookup(templateName)
	fmt.Println(newTemplate.Name())
	tempErr := newTemplate.Execute(&body, &data)
	fmt.Println(body.String())
	if tempErr != nil {
		log.Fatal("Template execution failed:", tempErr)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
