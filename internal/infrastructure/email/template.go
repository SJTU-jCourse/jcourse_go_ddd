package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/email"
)

type VerificationCodeTemplate struct{}

func NewVerificationCodeTemplate() email.Template {
	return &VerificationCodeTemplate{}
}

func (t *VerificationCodeTemplate) Execute(ctx context.Context, input *auth.VerificationCode) email.RenderedEmail {
	const templateContent = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>验证码</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 20px;
            background-color: #f9f9f9;
        }
        .header {
            text-align: center;
            margin-bottom: 20px;
        }
        .code {
            font-size: 24px;
            font-weight: bold;
            color: #007bff;
            text-align: center;
            padding: 10px;
            background-color: #e9ecef;
            border-radius: 5px;
            margin: 20px 0;
        }
        .footer {
            margin-top: 20px;
            font-size: 12px;
            color: #666;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>验证码</h2>
        </div>
        <p>您好！</p>
        <p>您的验证码是：</p>
        <div class="code">{{.Code}}</div>
        <p>验证码有效期为 5 分钟，请尽快使用。</p>
        <p>如果您没有请求验证码，请忽略此邮件。</p>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复。</p>
        </div>
    </div>
</body>
</html>
`

	tmpl, err := template.New("verification").Parse(templateContent)
	if err != nil {
		// Fallback to simple format if template parsing fails
		return email.RenderedEmail{
			Title: "验证码",
			Body:  fmt.Sprintf("您的验证码是：%s\n验证码有效期为 5 分钟，请尽快使用。", input.Code),
		}
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, input); err != nil {
		// Fallback to simple format if template execution fails
		return email.RenderedEmail{
			Title: "验证码",
			Body:  fmt.Sprintf("您的验证码是：%s\n验证码有效期为 5 分钟，请尽快使用。", input.Code),
		}
	}

	return email.RenderedEmail{
		Title: "验证码",
		Body:  body.String(),
	}
}