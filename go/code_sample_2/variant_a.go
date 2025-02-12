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

func SendGreetingsA(employeesCsvFilename string, now time.Time, smtpHost string, smtpPort int) error {
	employeesCsvFile, err := os.Open(employeesCsvFilename)
	if err != nil {
		return err
	}
	defer employeesCsvFile.Close()

	employeeCsvRecords, err := readEmployeeCsvRecords(employeesCsvFile)
	if err != nil {
		return err
	}

	employees, err := parseEmployees(employeeCsvRecords)
	if err != nil {
		return err
	}

	err = sendBirthdayGreetings(employees, now, smtpHost, smtpPort)
	if err != nil {
		return err
	}

	return nil
}

func readEmployeeCsvRecords(employeesCsvFile *os.File) ([][]string, error) {
	employeesCsvReader := csv.NewReader(employeesCsvFile)
	employeesCsvReader.TrimLeadingSpace = true

	var employeeCsvRecordsIncludingHeader [][]string
	endOfCsvReached := false
	for !endOfCsvReached {
		employeeCsvRecord, err := employeesCsvReader.Read()
		if err == io.EOF {
			endOfCsvReached = true
		} else if err != nil {
			return nil, err
		} else {
			employeeCsvRecordsIncludingHeader = append(employeeCsvRecordsIncludingHeader, employeeCsvRecord)
		}
	}

	employeeRecordsWithoutHeader := employeeCsvRecordsIncludingHeader[1:]
	return employeeRecordsWithoutHeader, nil
}

func parseEmployees(employeeCsvRecords [][]string) ([]*Employee, error) {
	var employees []*Employee
	for _, employeeCsvRecord := range employeeCsvRecords {
		lastname := employeeCsvRecord[0]
		firstname := employeeCsvRecord[1]
		dateOfBirth := employeeCsvRecord[2]
		email := employeeCsvRecord[3]
		employee, err := NewEmployee(lastname, firstname, dateOfBirth, email)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func sendMessageA(smtpHost string, smtpPort int, sender, subject, body, recipient string) error {
	smtpAddr := net.JoinHostPort(smtpHost, strconv.Itoa(smtpPort))
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n"+"%s\r\n", recipient, subject, body))
	return smtp.SendMail(smtpAddr, nil, sender, []string{recipient}, msg)
}

func sendBirthdayGreetings(employees []*Employee, now time.Time, smtpHost string, smtpPort int) error {
	for _, employee := range employees {
		if employee.IsBirthday(now) {
			err := sendBirthdayGreeting(employee, smtpHost, smtpPort)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func sendBirthdayGreeting(employee *Employee, smtpHost string, smtpPort int) error {
	greetingMessageBody := strings.Replace("Happy Birthday, dear %NAME%", "%NAME%", employee.Firstname, -1)
	err := sendMessageA(smtpHost, smtpPort, "sender@here.com", "Happy Birthday!", greetingMessageBody, employee.Email)
	if err != nil {
		return err
	}

	return nil
}
