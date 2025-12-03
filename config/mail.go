package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/wneessen/go-mail"
)

func SendResetEmail(toEmail, resetURL string, expiryMinutes int) error {
	// prepare env variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))

	if err != nil {
		return fmt.Errorf("invalid SMTP_PORT value: %w", err)
	}

	smtpSenderName := os.Getenv("SMTP_SENDER_USER")
	smtpEmail := os.Getenv("SMTP_EMAIL_FROM")
	smtpPassword := os.Getenv("SMTP_EMAIL_PASS")

	// fmt.Println(smtpHost, smtpPort, smtpSenderName, smtpEmail, smtpPassword)

	message := mail.NewMsg()

	if err := message.From(smtpSenderName); err != nil {
		return fmt.Errorf("failed to set From address: %w", err)
	}

	if err := message.To(toEmail); err != nil {
		return fmt.Errorf("failed to set To address: %w", err)
	}

	message.Subject("Reset your password")

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<style>
			*, ::after, ::before, ::backdrop, ::file-selector-button {
				margin: 0;
				padding: 0;
				box-sizing: border-box;
				border: 0 solid;

			}
			
			h1, h2, h3, h4, h5, h6 {
				font-size: inherit;
				font-weight: inherit;
			}
			
			ol, ul, menu {
				list-style: none;
			}
		</style>
		<body style="font-family: 'Segoe UI', Arial, sans-serif;">
			<div
				style="height: auto; min-height: 100%%; width: 100%%; position: absolute;
				background-image: linear-gradient(135deg, #0EA5E9, #003366);"
			>
				<div
					style="padding-left: 40px; padding-right: 40px; margin-left: auto; margin-right: auto;
					background-color: #FFFFFF;"
				>
					<h1
						style="padding-top: 20px; padding-bottom: 20px; font-size: 24px; line-height: 1.333;
						color: #003366; font-weight: 700;"
					>
						TEC UKDC
					</h1>
				</div>
				<div
					style="width: calc(10/12 * 100%%); margin-left: auto; margin-right: auto; margin-top: 60px;
					background-color: #FFFFFF; padding-top: 40px; padding-bottom: 40px; border-radius: 12px;
					text-align: center;"
				>
					<p
						style="font-size: 18px; line-height: 1.555; margin-top: 12px; margin-bottom: 24px;
						padding-left: 16px; padding-right: 16px;"
					>
						<b>Hello, </b><br>
						You requested a password reset. Click the link below to reset your password.
					</p>
					<p>
						<a href="%s"
							style="width:100%%; padding: 12px; background-color: #003366; color: #FFFFFF;
							font-size: 18px; line-height: 1.555; font-weight: 600; border-radius: 8px; cursor: pointer;
							margin-top: 12px; margin-bottom: 12px; text-decoration-line: none;"
						>
							Reset Password
						</a>
					</p>
					<p
						style="font-size: 18px; line-height: 1.555; margin-top: 24px; margin-bottom: 12px;
						padding-left: 16px; padding-right: 16px;"
					>
						This link will expire in %d minutes. If you didn't request this, you can ignore this email.
					</p>
					<p
						style="font-size: 14px; line-height: 1.428; margin-top: 12px; margin-bottom: 12px;
						padding-left: 16px; padding-right: 16px; text-align: center; font-weight: 600;"
					>
						This is an automated e-mail.
					</p>
				</div>
				<div style="width: 100%%; height: 60px;">
				</div>
			</div>
		</body>
		</html>
	`, resetURL, expiryMinutes)

	message.SetBodyString(mail.TypeTextHTML, htmlBody)

	client, err := mail.NewClient(
		smtpHost,
		mail.WithPort(smtpPort),
		mail.WithSSL(),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(smtpEmail),
		mail.WithPassword(smtpPassword),
		mail.WithTimeout(10*time.Second),
	)

	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}

	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
