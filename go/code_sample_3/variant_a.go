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

type EmployeeCsvFile struct {
	name string
}

func SendGreetingsA(employeeCsv EmployeeCsvFile, now time.Time, smtpServer SmtpServer) error {
	employees, err := employeeCsv.ReadEmployees()
	if err != nil {
		return err
	}

	for _, employee := range employees {
		if employee.IsBirthday(now) {
			email := createBirthdayMail(employee)

			err = smtpServer.send(email)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createBirthdayMail(employee *Employee) Email {
	emailBody := strings.Replace("Happy Birthday, dear %NAME%", "%NAME%", employee.Firstname, -1)
	return Email{sender: "sender@here.com", subject: "Happy Birthday!", body: emailBody, recipient: employee.Email}
}

// --- employee_csv_file.go ---

func (employeesCsvFile *EmployeeCsvFile) ReadEmployees() ([]*Employee, error) {
	csvFile := CsvFile{name: employeesCsvFile.name}
	csvRecords, err := csvFile.ReadRecords()
	if err != nil {
		return nil, err
	}

	employees, err := employeesCsvFile.parseEmployees(csvRecords)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (employeesCsvFile *EmployeeCsvFile) parseEmployees(employeeCsvRecords [][]string) ([]*Employee, error) {
	var employees []*Employee = make([]*Employee, 0)
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

// --- csv_file.go ---

type CsvFile struct {
	name string
}

func (csvFile *CsvFile) ReadRecords() ([][]string, error) {
	fileHandle, err := os.Open(csvFile.name)
	if err != nil {
		return nil, err
	}
	defer fileHandle.Close()

	csvRecordsIncludingHeader, err := csvFile.readRecordsFromOpenFile(fileHandle)
	if err != nil {
		return nil, err
	}

	csvRecordsWithoutHeader := csvRecordsIncludingHeader[1:]
	return csvRecordsWithoutHeader, nil
}

func (csvFile *CsvFile) readRecordsFromOpenFile(fileHandle *os.File) ([][]string, error) {
	csvReader := csv.NewReader(fileHandle)
	csvReader.TrimLeadingSpace = true

	var csvRecords [][]string
	endOfCsvReached := false
	for !endOfCsvReached {
		csvRecord, err := csvReader.Read()
		if err == io.EOF {
			endOfCsvReached = true
		} else if err != nil {
			return nil, err
		} else {
			csvRecords = append(csvRecords, csvRecord)
		}
	}
	return csvRecords, nil
}

// --- smtp_server.go ---

type SmtpServer struct {
	host string
	port int
}

type Email struct {
	sender, subject, body, recipient string
}

func (server *SmtpServer) send(email Email) error {
	smtpAddr := net.JoinHostPort(server.host, strconv.Itoa(server.port))
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n"+"%s\r\n", email.recipient, email.subject, email.body))
	return smtp.SendMail(smtpAddr, nil, email.sender, []string{email.recipient}, msg)
}
