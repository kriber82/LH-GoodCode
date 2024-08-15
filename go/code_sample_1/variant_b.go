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

func SendGreetingsB(employeesCsvFilename string, now time.Time, smtpHost string, smtpPort int) error {
	employeesCsvFile, err := os.Open(employeesCsvFilename)
	if err != nil {
		return err
	}
	defer employeesCsvFile.Close()

	employeesCsvReader := csv.NewReader(employeesCsvFile)
	employeesCsvReader.TrimLeadingSpace = true
	headerHasBeenSkipped := false
	for {
		employeeCsvRecord, err := employeesCsvReader.Read()
		//handle end of file
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if !headerHasBeenSkipped {
			headerHasBeenSkipped = true
			continue
		}

		//parse csv record
		lastname := employeeCsvRecord[0]
		firstname := employeeCsvRecord[1]
		dateOfBirth := employeeCsvRecord[2]
		email := employeeCsvRecord[3]
		employee, err := NewEmployee(lastname, firstname, dateOfBirth, email)
		if err != nil {
			return err
		}

		//check if we need to send a birthday message
		if employee.IsBirthday(now) {
			//actually send birthday email
			greetingMessageBody := strings.Replace("Happy Birthday, dear %NAME%", "%NAME%", employee.Firstname, -1)
			err = sendMessageB(smtpHost, smtpPort, "sender@here.com", "Happy Birthday!", greetingMessageBody, employee.Email)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func sendMessageB(smtpHost string, smtpPort int, sender, subject, body, recipient string) error {
	smtpAddr := net.JoinHostPort(smtpHost, strconv.Itoa(smtpPort))
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n"+"%s\r\n", recipient, subject, body))
	return smtp.SendMail(smtpAddr, nil, sender, []string{recipient}, msg)
}
