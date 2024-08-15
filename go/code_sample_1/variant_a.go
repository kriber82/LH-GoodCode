package birthday

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

func SendGreetingsA(filename string, now time.Time, smtpHost string, smtpPort int) error {
	//open file
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	//create csv reader for file
	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	header := false
	for {
		rec, err := r.Read()
		//handle end of file
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		//skip header line
		if !header {
			header = true
			continue
		}

		//parse csv record
		e, err := NewEmployee(rec[0], rec[1], rec[2], rec[3]) // Lastname, Firstname, dateOfBirth, email
		if err != nil {
			return err
		}

		//check if we need to send a birthday message
		if e.IsBirthday(now) {
			//actually send birthday email
			body := strings.Replace("Happy Birthday, dear %NAME%", "%NAME%", e.Firstname, -1)
			err = sendMessageA(smtpHost, smtpPort, "sender@here.com", "Happy Birthday!", body, e.Email)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func sendMessageA(smtpHost string, smtpPort int, sender, subject, body, recipient string) error {
	smtpAddr := net.JoinHostPort(smtpHost, strconv.Itoa(smtpPort))
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n"+"%s\r\n", recipient, subject, body))
	return smtp.SendMail(smtpAddr, nil, sender, []string{recipient}, msg)
}
