package main

import (
    "fmt"
    "os"
    "path/filepath"
    "bufio"
    "strings"
    "net/smtp"
    "mime/multipart"
    "net/textproto"
    "bytes"
   
    "encoding/base64"
)

// ... (keep the main, browseFiles, and getNextAction functions as they were)
var senderEmail string
var senderPassword string

func main() {
	    // Check for environment variables
	        senderEmail = os.Getenv("EMAIL_SENDER")
		    senderPassword = os.Getenv("EMAIL_PASSWORD")

		        // If environment variables are not set, prompt the user
			    if senderEmail == "" || senderPassword == "" {
				            fmt.Println("Email sender or password not set in environment variables.")
					            fmt.Println("Please set EMAIL_SENDER and EMAIL_PASSWORD environment variables.")
						            fmt.Println("For example:")
							            fmt.Println("export EMAIL_SENDER=your.email@example.com")
								            fmt.Println("export EMAIL_PASSWORD=your_password")
									            os.Exit(1)
										        }

											    currentDir, _ := os.Getwd()
											        for {
													        browseFiles(currentDir)
														        currentDir = getNextAction(currentDir)
															    }
														    }


func browseFiles(dir string) {
    files, _ := os.ReadDir(dir)
    fmt.Printf("Current directory: %s\n", dir)
    for _, file := range files {
        fmt.Println(file.Name())
    }
}

func getNextAction(currentDir string) string {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("\nEnter command (cd/send/quit): ")
    command, _ := reader.ReadString('\n')
    command = strings.TrimSpace(command)

    switch command {
    case "cd":
        fmt.Print("Enter directory name: ")
        dirName, _ := reader.ReadString('\n')
        dirName = strings.TrimSpace(dirName)
        return filepath.Join(currentDir, dirName)
    case "send":
        sendFile(currentDir)
        return currentDir
    case "quit":
        os.Exit(0)
    }
    return currentDir
}














func sendFile(dir string) {
    reader := bufio.NewReader(os.Stdin)
    
    fmt.Print("Enter file name: ")
    fileName, _ := reader.ReadString('\n')
    fileName = strings.TrimSpace(fileName)
    
    filePath := filepath.Join(dir, fileName)
    
    // Check if file exists
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        fmt.Println("Error: File does not exist")
        return
    }
    
    fmt.Print("Enter recipient email: ")
    recipient, _ := reader.ReadString('\n')
    recipient = strings.TrimSpace(recipient)
    
    fmt.Print("Enter email subject: ")
    subject, _ := reader.ReadString('\n')
    subject = strings.TrimSpace(subject)
    
    fmt.Print("Enter email body (type 'EOF' on a new line when finished):\n")
    var bodyLines []string
    for {
        line, _ := reader.ReadString('\n')
        line = strings.TrimSpace(line)
        if line == "EOF" {
            break
        }
        bodyLines = append(bodyLines, line)
    }
    emailBody := strings.Join(bodyLines, "\n")
    
    // Read file content
    fileContent, err := os.ReadFile(filePath)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    
    // Prepare email
    from := senderEmail
    password := senderPassword
    
    to := []string{recipient}
    
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"
    
    // Create buffer for message
    var body bytes.Buffer
    writer := multipart.NewWriter(&body)
    
    // Add email body
    part, _ := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/plain"}})
    part.Write([]byte(emailBody))
    
    // Add attachment
    part, _ = writer.CreatePart(textproto.MIMEHeader{
        "Content-Type":              {"application/octet-stream"},
        "Content-Transfer-Encoding": {"base64"},
        "Content-Disposition":       {"attachment; filename=\"" + fileName + "\""},
    })
    encoder := base64.NewEncoder(base64.StdEncoding, part)
    encoder.Write(fileContent)
    encoder.Close()
    
    writer.Close()
    
    // Compose message
    message := []byte(fmt.Sprintf("To: %s\r\n"+
        "Subject: %s\r\n"+
        "MIME-Version: 1.0\r\n"+
        "Content-Type: multipart/mixed; boundary=%s\r\n"+
        "\r\n%s", recipient, subject, writer.Boundary(), body.String()))
    
    // Send email
    auth := smtp.PlainAuth("", from, password, smtpHost)
    err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
    if err != nil {
        fmt.Println("Error sending email:", err)
        return
    }
    
    fmt.Println("Email sent successfully with attachment!")
}
