package utils

import "fmt"

// JobUpdatedEmail builds the email sent to candidates who applied when a job is updated
func JobUpdatedEmail(candidateEmail, candidateName, jobTitle, company, location, description string) EmailMessage {
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 30px auto; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
    .header { background: #0284C7; padding: 30px; text-align: center; }
    .header h1 { color: #fff; margin: 0; font-size: 22px; }
    .body { padding: 30px; color: #333; line-height: 1.7; }
    .badge { display: inline-block; background: #E0F2FE; color: #0369A1; padding: 6px 16px; border-radius: 20px; font-weight: bold; font-size: 14px; margin-bottom: 20px; }
    .job-card { background: #F0F9FF; border-left: 4px solid #0284C7; padding: 16px 20px; border-radius: 4px; margin: 20px 0; }
    .job-card h2 { margin: 0 0 8px; color: #1F2937; font-size: 18px; }
    .job-card p { margin: 4px 0; color: #6B7280; font-size: 14px; }
    .footer { background: #F3F4F6; padding: 16px 30px; text-align: center; font-size: 12px; color: #9CA3AF; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>✏️ Job Updated</h1>
    </div>
    <div class="body">
      <p>Hi <strong>%s</strong>,</p>
      <span class="badge">📢 Job Details Updated</span>
      <p>A job you applied for has been updated by the recruiter. Here are the latest details:</p>
      <div class="job-card">
        <h2>%s</h2>
        <p>🏢 <strong>%s</strong></p>
        <p>📍 %s</p>
        <p style="margin-top:12px; color:#374151;">%s</p>
      </div>
      <p>Your application is still active. No action needed from your side.</p>
      <p>— The Job Portal Team</p>
    </div>
    <div class="footer">
      &copy; 2026 Job Portal. All rights reserved.
    </div>
  </div>
</body>
</html>`,
		candidateName, jobTitle, company, location, description,
	)

	return EmailMessage{
		To:      []string{candidateEmail},
		Subject: fmt.Sprintf("📢 Job Updated: %s at %s", jobTitle, company),
		Body:    body,
	}
}

// NewJobAlertEmail builds the email sent to candidates when a new job is posted
func NewJobAlertEmail(candidateEmail, candidateName, jobTitle, company, location, description string) EmailMessage {
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 30px auto; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
    .header { background: #4F46E5; padding: 30px; text-align: center; }
    .header h1 { color: #fff; margin: 0; font-size: 22px; }
    .body { padding: 30px; color: #333; }
    .job-card { background: #F9FAFB; border-left: 4px solid #4F46E5; padding: 16px 20px; border-radius: 4px; margin: 20px 0; }
    .job-card h2 { margin: 0 0 8px; color: #1F2937; font-size: 18px; }
    .job-card p { margin: 4px 0; color: #6B7280; font-size: 14px; }
    .btn { display: inline-block; margin-top: 24px; padding: 12px 28px; background: #4F46E5; color: #fff; text-decoration: none; border-radius: 6px; font-weight: bold; }
    .footer { background: #F3F4F6; padding: 16px 30px; text-align: center; font-size: 12px; color: #9CA3AF; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>🚀 New Job Opportunity!</h1>
    </div>
    <div class="body">
      <p>Hi <strong>%s</strong>,</p>
      <p>A new job has just been posted that might interest you:</p>
      <div class="job-card">
        <h2>%s</h2>
        <p>🏢 <strong>%s</strong></p>
        <p>📍 %s</p>
        <p style="margin-top:12px; color:#374151;">%s</p>
      </div>
      <p>Log in to your account to apply before the opportunity closes.</p>
    </div>
    <div class="footer">
      You're receiving this because you're registered as a candidate on Job Portal.<br>
      &copy; 2026 Job Portal. All rights reserved.
    </div>
  </div>
</body>
</html>`,
		candidateName, jobTitle, company, location, description,
	)

	return EmailMessage{
		To:      []string{candidateEmail},
		Subject: fmt.Sprintf("New Job Alert: %s at %s", jobTitle, company),
		Body:    body,
	}
}

// ApplicationConfirmationEmail builds the confirmation email sent to the candidate after applying
func ApplicationConfirmationEmail(candidateEmail, candidateName, jobTitle, company, location string) EmailMessage {
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 30px auto; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
    .header { background: #4F46E5; padding: 30px; text-align: center; }
    .header h1 { color: #fff; margin: 0; font-size: 22px; }
    .body { padding: 30px; color: #333; line-height: 1.7; }
    .success-badge { display: inline-block; background: #D1FAE5; color: #065F46; padding: 6px 16px; border-radius: 20px; font-weight: bold; font-size: 14px; margin-bottom: 20px; }
    .job-card { background: #F9FAFB; border-left: 4px solid #4F46E5; padding: 16px 20px; border-radius: 4px; margin: 20px 0; }
    .job-card h2 { margin: 0 0 8px; color: #1F2937; font-size: 18px; }
    .job-card p { margin: 4px 0; color: #6B7280; font-size: 14px; }
    .steps { margin: 24px 0; padding: 0; list-style: none; }
    .steps li { padding: 8px 0; border-bottom: 1px solid #F3F4F6; color: #374151; font-size: 14px; }
    .steps li:last-child { border-bottom: none; }
    .footer { background: #F3F4F6; padding: 16px 30px; text-align: center; font-size: 12px; color: #9CA3AF; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>✅ Application Submitted!</h1>
    </div>
    <div class="body">
      <p>Hi <strong>%s</strong>,</p>
      <span class="success-badge">✓ Successfully Applied</span>
      <p>Your application has been successfully submitted. Here's a summary:</p>
      <div class="job-card">
        <h2>%s</h2>
        <p>🏢 <strong>%s</strong></p>
        <p>📍 %s</p>
        <p>📋 Status: <strong>APPLIED</strong></p>
      </div>
      <p><strong>What happens next?</strong></p>
      <ul class="steps">
        <li>📬 The recruiter has been notified of your application</li>
        <li>👀 Your application will be reviewed shortly</li>
        <li>📧 You'll receive an email when your status is updated</li>
      </ul>
      <p>Good luck! We're rooting for you. 🤞</p>
      <p>— The Job Portal Team</p>
    </div>
    <div class="footer">
      &copy; 2026 Job Portal. All rights reserved.
    </div>
  </div>
</body>
</html>`,
		candidateName, jobTitle, company, location,
	)

	return EmailMessage{
		To:      []string{candidateEmail},
		Subject: fmt.Sprintf("✅ Application Submitted — %s at %s", jobTitle, company),
		Body:    body,
	}
}

// ApplicationReceivedEmail builds the email sent to a recruiter when someone applies
func ApplicationReceivedEmail(recruiterEmail, recruiterName, candidateName, candidateEmail, jobTitle string) EmailMessage {
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 30px auto; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
    .header { background: #059669; padding: 30px; text-align: center; }
    .header h1 { color: #fff; margin: 0; font-size: 22px; }
    .body { padding: 30px; color: #333; }
    .info-card { background: #F0FDF4; border-left: 4px solid #059669; padding: 16px 20px; border-radius: 4px; margin: 20px 0; }
    .info-card p { margin: 6px 0; color: #374151; font-size: 15px; }
    .footer { background: #F3F4F6; padding: 16px 30px; text-align: center; font-size: 12px; color: #9CA3AF; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>📩 New Application Received</h1>
    </div>
    <div class="body">
      <p>Hi <strong>%s</strong>,</p>
      <p>Someone just applied for your job posting <strong>"%s"</strong>.</p>
      <div class="info-card">
        <p>👤 <strong>Candidate:</strong> %s</p>
        <p>📧 <strong>Email:</strong> %s</p>
        <p>💼 <strong>Applied for:</strong> %s</p>
      </div>
      <p>Log in to your recruiter dashboard to review the application and update its status.</p>
    </div>
    <div class="footer">
      &copy; 2026 Job Portal. All rights reserved.
    </div>
  </div>
</body>
</html>`,
		recruiterName, jobTitle, candidateName, candidateEmail, jobTitle,
	)

	return EmailMessage{
		To:      []string{recruiterEmail},
		Subject: fmt.Sprintf("New Application for \"%s\" from %s", jobTitle, candidateName),
		Body:    body,
	}
}

// ApplicationStatusEmail builds the email sent to a candidate when their status changes
func ApplicationStatusEmail(candidateEmail, candidateName, jobTitle, company string, status string) EmailMessage {
	var (
		headerColor string
		emoji       string
		headline    string
		message     string
	)

	switch status {
	case "ACCEPTED":
		headerColor = "#059669"
		emoji = "🎉"
		headline = "Congratulations! You've been accepted!"
		message = fmt.Sprintf(
			"We're thrilled to let you know that your application for <strong>%s</strong> at <strong>%s</strong> has been <strong style='color:#059669;'>ACCEPTED</strong>.<br><br>The recruiter will be in touch with you shortly regarding next steps.",
			jobTitle, company,
		)
	case "REJECTED":
		headerColor = "#DC2626"
		emoji = "😔"
		headline = "Application Update"
		message = fmt.Sprintf(
			"Thank you for your interest in <strong>%s</strong> at <strong>%s</strong>. After careful consideration, the recruiter has decided not to move forward with your application at this time.<br><br>Don't be discouraged — keep applying and the right opportunity will come!",
			jobTitle, company,
		)
	case "REVIEWED":
		headerColor = "#D97706"
		emoji = "👀"
		headline = "Your Application is Being Reviewed"
		message = fmt.Sprintf(
			"Good news! Your application for <strong>%s</strong> at <strong>%s</strong> is currently being <strong style='color:#D97706;'>REVIEWED</strong> by the recruiter.<br><br>We'll notify you as soon as there's an update.",
			jobTitle, company,
		)
	default:
		headerColor = "#4F46E5"
		emoji = "📋"
		headline = "Application Status Update"
		message = fmt.Sprintf(
			"Your application for <strong>%s</strong> at <strong>%s</strong> has been updated to <strong>%s</strong>.",
			jobTitle, company, status,
		)
	}

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 30px auto; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
    .header { background: %s; padding: 30px; text-align: center; }
    .header h1 { color: #fff; margin: 0; font-size: 22px; }
    .body { padding: 30px; color: #333; line-height: 1.6; }
    .footer { background: #F3F4F6; padding: 16px 30px; text-align: center; font-size: 12px; color: #9CA3AF; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>%s %s</h1>
    </div>
    <div class="body">
      <p>Hi <strong>%s</strong>,</p>
      <p>%s</p>
      <p>Best of luck on your job search journey!</p>
      <p>— The Job Portal Team</p>
    </div>
    <div class="footer">
      &copy; 2026 Job Portal. All rights reserved.
    </div>
  </div>
</body>
</html>`,
		headerColor, emoji, headline,
		candidateName, message,
	)

	subjectMap := map[string]string{
		"ACCEPTED": fmt.Sprintf("🎉 You got the job! — %s at %s", jobTitle, company),
		"REJECTED": fmt.Sprintf("Application Update — %s at %s", jobTitle, company),
		"REVIEWED": fmt.Sprintf("Your application is being reviewed — %s", jobTitle),
	}

	subject, ok := subjectMap[status]
	if !ok {
		subject = fmt.Sprintf("Application Status Update — %s", jobTitle)
	}

	return EmailMessage{
		To:      []string{candidateEmail},
		Subject: subject,
		Body:    body,
	}
}
