<template>
  <footer class="footer">
    <div class="footer-inner">
      <div class="footer-grid">

        <!-- 聯絡我們 -->
        <div class="footer-col">
          <h3 class="footer-heading">
            <i class="pi pi-phone" /> 聯絡我們
          </h3>
          <ul class="footer-list">
            <li>
              <i class="pi pi-phone footer-icon" />
              <span>客服專線：0927-066-022<br /><span class="muted">（週一至週五 09:00–18:00）</span></span>
            </li>
            <li>
              <i class="pi pi-envelope footer-icon" />
              <span>service@silverlord.com.tw</span>
            </li>
            <li>
              <i class="pi pi-map-marker footer-icon" />
              <span>暫無</span>
            </li>
          </ul>
        </div>

        <!-- 匯款資訊 -->
        <div class="footer-col">
          <h3 class="footer-heading">
            <i class="pi pi-credit-card" /> 匯款資訊
          </h3>
          <div class="bank-card">
            <span class="bank-notice">下單後將於訂單郵件提供</span>
          </div>
        </div>

        <!-- 訂購須知 -->
        <div class="footer-col">
          <h3 class="footer-heading">
            <i class="pi pi-shield" /> 訂購須知
          </h3>
          <ul class="notice-list">
            <li>貴金屬商品價格隨國際盤價波動，以結帳當下系統顯示為準。</li>
            <li>完成下單後請於 48 小時內完成匯款，逾期將自動取消訂單。</li>
            <li>實體貴金屬商品依消保法<span class="red">不適用七天鑑賞期</span>，售出概不退換。</li>
            <li>款項確認後 2 個工作天內物流寄出。</li>
          </ul>
        </div>

        <!-- 官方 LINE -->
        <div class="footer-col">
          <h3 class="footer-heading">
            <i class="pi pi-comment" /> 官方 LINE
          </h3>
          <p class="footer-desc">
            加入官方 LINE 帳號，獲得最新國際銀價資訊與專屬優惠通知，並提供一對一客服諮詢。
          </p>
          <a class="line-btn" href="#">
            <i class="pi pi-comment" />
            加入好友 @silverlord
          </a>
        </div>

      </div>

      <div class="footer-bottom">
        &copy; {{ year }} SILVER LORD 貴金屬買賣平台．All rights reserved.
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
const year = new Date().getFullYear()
</script>

<style scoped>
.footer {
  background: #0f172a;
  color: #94a3b8;
  padding: 3.5rem 0 0;
  border-top: 1px solid #1e293b;
  margin-top: auto;
}

.footer-inner {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 1.5rem;
}

.footer-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 2.5rem;
  padding-bottom: 3rem;
}

.footer-heading {
  font-family: 'Playfair Display', Georgia, serif;
  color: #f1f5f9;
  font-size: 1.05rem;
  font-weight: 600;
  letter-spacing: 0.04em;
  margin-bottom: 1.25rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* 聯絡 */
.footer-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
  font-size: 0.85rem;
  line-height: 1.6;
}

.footer-list li {
  display: flex;
  align-items: flex-start;
  gap: 0.6rem;
}

.footer-icon {
  margin-top: 3px;
  flex-shrink: 0;
}

.muted {
  color: #475569;
  font-size: 0.8rem;
}

/* 匯款 */
.bank-card {
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 10px;
  padding: 1rem;
  font-size: 0.83rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.bank-row {
  display: flex;
  gap: 0.75rem;
}

.bank-label {
  color: #64748b;
  width: 4.5rem;
  flex-shrink: 0;
}

.bank-notice {
  color: #94a3b8;
  font-size: 0.85rem;
  line-height: 1.6;
}

.bank-account {
  font-family: 'Courier New', monospace;
  font-size: 1.1rem;
  font-weight: 700;
  color: #f1f5f9;
  letter-spacing: 0.08em;
  border-top: 1px solid #334155;
  padding-top: 0.6rem;
  margin-top: 0.25rem;
}

/* 須知 */
.notice-list {
  list-style: disc;
  padding-left: 1.2rem;
  margin: 0;
  font-family: 'EB Garamond', Georgia, serif;
  font-size: 0.95rem;
  line-height: 1.9;
  color: #64748b;
  display: flex;
  flex-direction: column;
  gap: 0.1rem;
}

.red {
  color: #f87171;
  font-weight: 600;
}

/* LINE */
.footer-desc {
  font-size: 0.85rem;
  line-height: 1.7;
  margin-bottom: 1.25rem;
}

.line-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: #06C755;
  color: white;
  font-weight: 700;
  font-size: 0.9rem;
  padding: 0.6rem 1.2rem;
  border-radius: 8px;
  transition: background 0.15s;
  cursor: pointer;
  text-decoration: none;
}

.line-btn:hover {
  background: #05a547;
}

/* Bottom bar */
.footer-bottom {
  border-top: 1px solid #1e293b;
  padding: 1.5rem 0;
  text-align: center;
  font-size: 0.8rem;
  color: #475569;
}
</style>
