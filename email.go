package main

import (
	"bufio"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

func (c *Config) verifyEmail() {

	if c.Interval[0] == '+' || c.Interval[0] == '-' {
		c.Interval = strings.Replace(c.Interval, string(c.Interval[0]), "", -1)
	}

	if i, _ := os.Stat(c.Logs); !(i.Size() > 0) {
		log.Printf("\n%s\n\n", "Log file is empty.")
		return
	}

}

func (c *Config) emailLogs() {

	reset, err := strconv.ParseBool(c.Reset)

	if err != nil {
		log.Fatalln(err)
	}

	c.generateHeaderMap()

	c.insertMessage()

	err = c.sendEmail()
	check(err)

	if reset {
		if err := os.Remove(c.Logs); err != nil {
			log.Println(err)
		}
		if _, err := os.Create(c.Logs); err != nil {
			log.Println(err)
		}
	}

}

func (c *Config) generateHeaderMap() {

	c.Headers["From"] = fmt.Sprintf(`"%s" <%s>`, c.From.Name, c.From.Email)
	c.Headers["To"] = fmt.Sprintf(`"%s" <%s>`, c.To.Name, c.To.Email)
	c.Headers["Subject"] = c.Subject
	c.Headers["MIME-version"] = "1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	for title, data := range c.Headers {
		c.Message += fmt.Sprintf("%s: %s\r\n", title, data)
	}
	c.Message += "\r\n"

}

func (c *Config) insertMessage() {

	file, err := os.Open(c.Logs)
	check(err)

	defer file.Close()

	c.Message += "<div style=\"font-family: monospace;background: #ecf0f1;padding: 20px;border-radius: 9px;font-size: 150%;margin: 30px;\">"
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c.Message += scanner.Text() + "<br>"
	}
	c.Message += "<br><br>Generated by <a href=\"https://github.com/muhammadmuzzammil1998/Log-mailer\">Log Mailer</a> on " + time.Now().Format(time.RFC1123Z) + "</div>"

}

func (c *Config) sendEmail() error {

	err := smtp.SendMail(
		c.Server+":"+c.Port,
		smtp.PlainAuth("", c.Credentials.Username, c.Credentials.Password, c.Server),
		c.Headers["From"],
		[]string{c.Headers["To"]},
		[]byte(c.Message),
	)
	if err != nil {
		return err
	}

	log.Printf("\n%s\n\n", c.Message)
	return nil
}