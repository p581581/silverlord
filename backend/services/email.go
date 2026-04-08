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
	Items            []models.OrderItem
}

var orderConfirmTpl = template.Must(template.New("order").Funcs(template.FuncMap{
	"mult": func(a, b int) int { return a * b },
}).Parse(`
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><style>
  body { font-family: sans-serif; background:#f5f5f5; padding:20px; }
  .card { background:white; border-radius:12px; padding:24px; max-width:560px; margin:0; }
  h2 { color:#1a1a2e; }
  .order-no { background:#f0f0f0; padding:8px 16px; border-radius:6px; font-family:monospace; font-size:1.1em; }
  table { width:100%; border-collapse:collapse; margin:16px 0; }
  th { background:#1a1a2e; color:white; padding:8px; text-align:left; }
  td { padding:8px; border-bottom:1px solid #eee; }
  .total { font-size:1.2em; font-weight:bold; color:#c0392b; text-align:right; margin-top:8px; }
  .hint { color:#888; font-size:0.85em; margin-top:16px; }
</style></head>
<body>
<div class="card">
  <h2>🥈 白銀平台 — 訂單確認</h2>
  <p>感謝您的訂購！您的訂單已建立，請於 <strong>30 分鐘內</strong>完成匯款。</p>
  <p>訂單編號：<span class="order-no">{{.OrderNo}}</span></p>
  <table>
    <tr><th>商品</th><th>數量</th><th>單價</th><th>小計</th></tr>
    {{range .Items}}
    <tr>
      <td>{{.ProductName}}</td>
      <td>{{.Quantity}}</td>
      <td>NT$ {{.LockedUnitPrice}}</td>
      <td>NT$ {{mult .LockedUnitPrice .Quantity}}</td>
    </tr>
    {{end}}
  </table>
  <div class="total">合計：NT$ {{.TotalLockedPrice}}</div>
  <p class="hint">匯款完成後，請至網站「匯款回報」頁面填寫末五碼，我們將盡快確認。</p>
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
	m.SetHeader("From", cfg.SmtpFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("【白銀平台】訂單確認 %s", data.OrderNo))
	m.SetBody("text/html", buf.String())

	d := gomail.NewDialer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUser, cfg.SmtpPassword)
	return d.DialAndSend(m)
}
