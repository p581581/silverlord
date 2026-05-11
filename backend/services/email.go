package services

import (
	"bytes"
	"fmt"
	"html/template"

	"silverlord/backend/config"
	"silverlord/backend/models"

	gomail "gopkg.in/gomail.v2"
)

type OrderEmailData struct {
	OrderNo          string
	TotalLockedPrice int
	ShipFee          int
	Items            []models.OrderItem
	PaymentURL       string
	BankName         string
	BankAccount      string
	BankAccountName  string
}

var orderConfirmTpl = template.Must(template.New("order").Funcs(template.FuncMap{
	"mult": func(a, b int) int { return a * b },
	"sub":  func(a, b int) int { return a - b },
}).Parse(`
<!DOCTYPE html>
<html lang="zh-TW">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    background: #f0ede8;
    font-family: 'Georgia', 'Times New Roman', serif;
    color: #2c2c2c;
    padding: 40px 16px;
  }
  .wrapper { max-width: 580px; margin: 0 auto; }

  /* Header */
  .header {
    background: #1a1a2e;
    border-radius: 12px 12px 0 0;
    padding: 36px 40px 28px;
    text-align: center;
  }
  .header-brand {
    font-size: 11px;
    letter-spacing: 0.25em;
    text-transform: uppercase;
    color: #a89060;
    margin-bottom: 10px;
  }
  .header-title {
    font-size: 24px;
    font-weight: normal;
    color: #f5f0e8;
    letter-spacing: 0.06em;
  }
  .header-divider {
    width: 40px;
    height: 1px;
    background: #a89060;
    margin: 14px auto 0;
  }

  /* Body */
  .body {
    background: #ffffff;
    padding: 36px 40px;
  }
  .greeting {
    font-size: 15px;
    color: #4a4a4a;
    line-height: 1.7;
    margin-bottom: 28px;
  }

  /* Order no */
  .order-block {
    border: 1px solid #e8e3dc;
    border-radius: 8px;
    padding: 16px 20px;
    margin-bottom: 28px;
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .order-label {
    font-size: 11px;
    letter-spacing: 0.15em;
    text-transform: uppercase;
    color: #a89060;
    white-space: nowrap;
  }
  .order-no {
    font-family: 'Courier New', monospace;
    font-size: 16px;
    font-weight: bold;
    color: #1a1a2e;
    letter-spacing: 0.05em;
  }

  /* Table */
  .section-label {
    font-size: 11px;
    letter-spacing: 0.15em;
    text-transform: uppercase;
    color: #a89060;
    margin-bottom: 10px;
  }
  table {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 4px;
    font-size: 14px;
  }
  thead tr {
    border-bottom: 1.5px solid #1a1a2e;
  }
  thead th {
    padding: 8px 6px;
    text-align: left;
    font-size: 11px;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: #888;
    font-weight: normal;
  }
  thead th:last-child { text-align: right; }
  tbody tr {
    border-bottom: 1px solid #f0ede8;
  }
  tbody td {
    padding: 10px 6px;
    color: #2c2c2c;
    vertical-align: middle;
  }
  tbody td:last-child { text-align: right; }

  /* Subtotal rows */
  .subtotal-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 9px 6px;
    border-top: 1px solid #f0ede8;
    gap: 8px;
  }
  .subtotal-label {
    font-size: 12px;
    color: #888;
  }
  .subtotal-amount {
    font-size: 14px;
    color: #4a4a4a;
  }

  /* Total */
  .total-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 8px;
    border-top: 1.5px solid #1a1a2e;
    padding: 16px 6px 0;
    margin-top: 6px;
    margin-bottom: 36px;
  }
  .total-label {
    font-size: 12px;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: #888;
  }
  .total-amount {
    font-size: 22px;
    font-weight: bold;
    color: #1a1a2e;
    letter-spacing: 0.03em;
  }

  /* Bank info */
  .bank-block {
    background: #f9f7f4;
    border: 1px solid #e8e3dc;
    border-radius: 8px;
    padding: 18px 20px;
    margin-bottom: 28px;
  }
  .bank-title {
    font-size: 11px;
    letter-spacing: 0.15em;
    text-transform: uppercase;
    color: #a89060;
    margin-bottom: 12px;
  }
  .bank-row {
    display: flex;
    justify-content: space-between;
    padding: 5px 0;
    font-size: 14px;
    border-bottom: 1px solid #ede9e2;
  }
  .bank-row:last-child { border-bottom: none; }
  .bank-key { color: #888; }
  .bank-val { font-weight: 600; color: #1a1a2e; letter-spacing: 0.04em; }

  /* CTA Button */
  .cta-wrap { text-align: center; margin-bottom: 32px; }
  .cta-btn {
    display: inline-block;
    background: #1a1a2e;
    color: #f5f0e8 !important;
    text-decoration: none;
    font-family: 'Georgia', serif;
    font-size: 14px;
    letter-spacing: 0.12em;
    padding: 14px 36px;
    border-radius: 4px;
  }

  /* Hint */
  .hint {
    background: #f9f7f4;
    border-left: 3px solid #a89060;
    border-radius: 0 6px 6px 0;
    padding: 14px 18px;
    font-size: 13px;
    color: #666;
    line-height: 1.7;
  }

  /* Footer */
  .footer {
    background: #1a1a2e;
    border-radius: 0 0 12px 12px;
    padding: 20px 40px;
    text-align: center;
    font-size: 11px;
    letter-spacing: 0.1em;
    color: #6a6a7e;
  }
</style>
</head>
<body>
<div class="wrapper">

  <div class="header">
    <div class="header-brand">Silver Lord</div>
    <div class="header-title">訂單確認通知</div>
    <div class="header-divider"></div>
  </div>

  <div class="body">
    <p class="greeting">
      感謝您的訂購。<br>
      您的訂單已成立，請於 <strong>48 小時內</strong>完成匯款，以利我們盡速為您處理出貨事宜。
    </p>

    <div class="order-block">
      <span class="order-label">訂單編號</span>
      <span class="order-no">{{.OrderNo}}</span>
    </div>

    <p class="section-label">訂購明細</p>
    <table>
      <thead>
        <tr>
          <th>商品</th>
          <th style="text-align:center">數量</th>
          <th style="text-align:right">單價</th>
          <th>小計</th>
        </tr>
      </thead>
      <tbody>
        {{range .Items}}
        <tr>
          <td>{{.ProductName}}</td>
          <td style="text-align:center">{{.Quantity}}</td>
          <td style="text-align:right">NT$ {{.LockedUnitPrice}}</td>
          <td>NT$ {{mult .LockedUnitPrice .Quantity}}</td>
        </tr>
        {{end}}
      </tbody>
    </table>

    <div class="subtotal-row">
      <span class="subtotal-label">商品小計</span>
      <span class="subtotal-amount">NT$ {{sub .TotalLockedPrice .ShipFee}}</span>
    </div>
    <div class="subtotal-row">
      <span class="subtotal-label">運費</span>
      <span class="subtotal-amount">{{if eq .ShipFee 0}}免運{{else}}NT$ {{.ShipFee}}{{end}}</span>
    </div>
    <div class="total-row">
      <span class="total-label">訂單總計</span>
      <span class="total-amount">NT$ {{.TotalLockedPrice}}</span>
    </div>

    {{if or .BankName .BankAccount .BankAccountName}}
    <div class="bank-block">
      <div class="bank-title">匯款資訊</div>
      {{if .BankName}}
      <div class="bank-row">
        <span class="bank-key">銀行</span>
        <span class="bank-val">{{.BankName}}</span>
      </div>
      {{end}}
      {{if .BankAccountName}}
      <div class="bank-row">
        <span class="bank-key">戶名</span>
        <span class="bank-val">{{.BankAccountName}}</span>
      </div>
      {{end}}
      {{if .BankAccount}}
      <div class="bank-row">
        <span class="bank-key">帳號</span>
        <span class="bank-val">{{.BankAccount}}</span>
      </div>
      {{end}}
    </div>
    {{end}}

    {{if .PaymentURL}}
    <div class="cta-wrap">
      <a href="{{.PaymentURL}}" class="cta-btn">完成匯款後，點此回報</a>
    </div>
    {{end}}

    <div class="hint">
      匯款完成後，請點擊上方按鈕（或前往網站「匯款回報」頁面），輸入匯款末五碼通知我們。<br>
      確認入帳後將盡速安排出貨，如有任何問題歡迎與我們聯繫。
    </div>
  </div>

  <div class="footer">
    © Silver Lord &nbsp;·&nbsp; 此為系統自動發送，請勿直接回覆
  </div>

</div>
</body>
</html>
`))

// SendOrderConfirmEmail 寄送訂單確認信，SMTP 設定來自 config（env）
func SendOrderConfirmEmail(cfg *config.Config, to string, data OrderEmailData) error {
	if cfg.SmtpHost == "" || to == "" {
		return nil
	}

	var buf bytes.Buffer
	if err := orderConfirmTpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("email template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", cfg.SmtpFrom, cfg.SmtpFromName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("【Silver Lord】訂單確認 %s", data.OrderNo))
	m.SetBody("text/html", buf.String())

	d := gomail.NewDialer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUser, cfg.SmtpPassword)
	return d.DialAndSend(m)
}
