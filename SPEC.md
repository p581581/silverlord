# 白銀即時報價買賣平台 — 系統規格書

> 最後更新：2026-04-07

---

## 1. 系統概述

白銀（XAG）即時報價電商平台，整合 Kitco 銀價 API、動態定價、購物車、結帳、匯款回報、庫存管理與後台訂單管理，支援 Telegram 多人通知與 SMTP 訂單確認信。

**技術棧**

| 層級 | 技術 |
|------|------|
| Backend | Go + Echo v4 + xorm + PostgreSQL |
| Frontend | Vue 3 + Vite + PrimeVue v4 + Pinia |
| 部署 | 單一二進位（`-tags prod` 嵌入 Vue dist）|
| 多語言 | vue-i18n（繁中 / English / 日本語）|

---

## 2. 目錄結構

```
silverLord/
├── backend/
│   ├── main.go
│   ├── frontend_dev.go      # build tag: !prod（no-op）
│   ├── frontend_prod.go     # build tag: prod（embed web/）
│   ├── build-linux.bat      # npm build → xcopy → go build -tags prod
│   ├── config/config.go
│   ├── database/db.go
│   ├── cache/price_cache.go
│   ├── middleware/jwt.go
│   ├── models/
│   │   ├── product.go
│   │   ├── order.go
│   │   ├── system_config.go
│   │   └── spot_price.go
│   ├── handlers/
│   │   ├── product.go
│   │   ├── order.go
│   │   ├── admin.go
│   │   ├── upload_helper.go
│   │   └── errors.go
│   └── services/
│       ├── cron.go
│       ├── pricing.go
│       ├── email.go
│       └── telegram.go
└── frontend/
    └── src/
        ├── views/client/
        │   ├── ProductList.vue
        │   ├── ProductDetail.vue
        │   ├── Cart.vue
        │   ├── Checkout.vue
        │   ├── Payment.vue
        │   └── OrderQuery.vue
        ├── views/admin/
        │   ├── Login.vue
        │   ├── Layout.vue
        │   ├── Dashboard.vue
        │   ├── Products.vue
        │   ├── Orders.vue
        │   ├── Config.vue
        │   └── Shipping.vue
        ├── components/
        │   ├── shared/Navbar.vue
        │   └── client/ProductCard.vue
        ├── stores/
        │   ├── auth.ts
        │   ├── cart.ts
        │   └── price.ts
        └── i18n/locales/
            ├── zh-TW.ts
            ├── en.ts
            └── ja.ts
```

---

## 3. 環境變數（`backend/.env`）

| 變數 | 說明 | 預設值 |
|------|------|--------|
| `DB_HOST` | PostgreSQL 主機 | `localhost` |
| `DB_PORT` | PostgreSQL 埠 | `5432` |
| `DB_USER` | 資料庫帳號 | `postgres` |
| `DB_PASSWORD` | 資料庫密碼 | `postgres` |
| `DB_NAME` | 資料庫名稱 | `silverlord` |
| `JWT_SECRET` | JWT 簽章金鑰 | `change-me-in-production` |
| `ADMIN_PASSWORD` | 後台登入密碼 | `admin123` |
| `PORT` | 伺服器埠 | `8080` |
| `SMTP_HOST` | SMTP 主機 | — |
| `SMTP_PORT` | SMTP 埠 | `587` |
| `SMTP_USER` | SMTP 帳號 | — |
| `SMTP_PASSWORD` | SMTP 密碼 | — |
| `SMTP_FROM` | 寄件地址 | — |

前端（`frontend/.env.production`）：
- `VITE_API_BASE_URL=`（留空，prod 同源；dev 設 `http://localhost:8080`）

---

## 4. 資料模型

### 4.1 Product（商品）

| 欄位 | 型別 | 說明 |
|------|------|------|
| `id` | int64 PK autoincr | |
| `name` | string notnull | 商品名稱 |
| `description` | text | 商品描述 |
| `image_url` | string | 圖片路徑（`/uploads/...`）|
| `weight` | float64 notnull | 重量（troy oz）|
| `labor_cost` | int default 0 | 利潤（TWD）|
| `premium_percent` | float64 default 5.0 | 溢價百分比（%）|
| `is_active` | bool default true | 是否上架 |
| `product_type` | string default 'in_stock' | `in_stock`（庫存品）/ `pre_order`（預購）|
| `stock_qty` | int default 0 | 庫存數量（僅 in_stock 有效）|
| `expected_arrival` | string default '' | 預計到貨日 YYYY-MM-DD（僅 pre_order 有效）|

**定價公式：**
```
單價(TWD) = Round(XAG/USD × weight × USD/TWD × (1 + premium_percent/100)) + labor_cost
```
- 若 `manual_ex_rate > 0`：使用固定匯率，否則用即時匯率
- `manual_offset` 加在 XAG/USD 之前做微調

### 4.2 Order（訂單）

| 欄位 | 型別 | 說明 |
|------|------|------|
| `id` | int64 PK autoincr | |
| `order_no` | string unique notnull | `SN{YYYYMMDD}{4位流水}` |
| `total_locked_price` | int notnull | 鎖定總額（含運費，TWD）|
| `email` | string | 買家 Email（必填）|
| `name` | string | 收件人姓名（必填）|
| `phone` | string | 收件人電話（必填）|
| `ship_method` | string | `mail` / `711` / `pickup` |
| `ship_address` | string | 郵寄地址（mail 時必填）|
| `ship_store` | string | 7-11 門市名稱（711 時必填）|
| `ship_fee` | int default 0 | 實際運費（TWD）|
| `remit_info` | string | 匯款末五碼 |
| `status` | int default 0 | 訂單狀態（見下表）|
| `created_at` | time | 建立時間（xorm 自動）|

**訂單狀態：**

| 值 | 名稱 | 說明 |
|----|------|------|
| `0` | 待匯款 | 初始狀態，等待付款 |
| `1` | 已通知 | 用戶回報匯款後 |
| `2` | 已確認/出貨 | 管理員確認入帳後 |
| `3` | 不成立 | 管理員標記取消 |

**狀態轉換規則：**

| 來源 | 目標 | 操作者 | 庫存處理 |
|------|------|--------|---------|
| 0 | 1 | 用戶（回報匯款）| 無 |
| 1 | 0 | 管理員（退回）| 無 |
| 1 | 2 | 管理員（確認入帳）| 無 |
| 0 | 3 | 管理員（不成立）| 退庫存 |
| 2 | 3 | 管理員（不成立）| 退庫存 |
| 1 | 3 | **不允許** | 需先退回至 0 |
| 刪除 | — | 管理員 | 若 status ≠ 3 則退庫存 |

### 4.3 OrderItem（訂單明細）

| 欄位 | 型別 | 說明 |
|------|------|------|
| `id` | int64 PK autoincr | |
| `order_id` | int64 index | 關聯訂單 |
| `product_id` | int64 | 商品 ID |
| `product_name` | string | 下單時快照 |
| `weight` | float64 | 下單時快照（troy oz）|
| `quantity` | int default 1 | 數量 |
| `locked_unit_price` | int | 下單時鎖定單價（TWD）|

### 4.4 SystemConfig（系統設定，唯一一筆 id=1）

| 欄位 | 型別 | 說明 |
|------|------|------|
| `id` | int PK | 固定為 1 |
| `manual_offset` | float64 default 0 | 價格修正（加在 XAG/USD 上）|
| `manual_ex_rate` | float64 default 0 | 固定匯率（0 = 使用即時）|
| `tg_bot_token` | string | Telegram Bot Token |
| `tg_chat_ids` | string | 逗號分隔 Chat ID（支援多組）|
| `ship_fee_mail` | int default 0 | 郵寄運費（TWD）|
| `ship_fee_711` | int default 0 | 7-11 店到店運費（TWD）|
| `ship_fee_pickup` | int default 0 | 自取運費（TWD）|

### 4.5 SpotPrice（報價紀錄）

| 欄位 | 型別 | 說明 |
|------|------|------|
| `id` | int64 PK autoincr | |
| `xag_usd` | float64 | 白銀現貨美金價（每盎司）|
| `usd_twd` | float64 | 美金台幣匯率 |
| `created_at` | time index | 紀錄時間 |

---

## 5. API 路由

### 5.1 公開 API（`/api/v1`）

| Method | 路徑 | 說明 |
|--------|------|------|
| POST | `/admin/login` | 後台登入，回傳 JWT（24hr）|
| GET | `/products` | 商品列表（僅 is_active=true，含即時報價）|
| GET | `/products/:id` | 單一商品詳情（含即時報價）|
| GET | `/spot-price` | 即時銀價快照 |
| GET | `/ship-fees` | 運費設定（mail / 711）|
| POST | `/checkout` | 建立訂單（庫存品扣庫存）|
| POST | `/report-payment` | 用戶回報匯款（status 0→1）|
| POST | `/orders/:order_no/resend-email` | 重發訂單確認信 |
| GET | `/orders/:order_no` | 查詢訂單狀態與明細 |

### 5.2 管理員 API（`/api/admin`，需 JWT Bearer Token）

**商品管理**

| Method | 路徑 | 說明 |
|--------|------|------|
| GET | `/products` | 取得所有商品（含下架）|
| POST | `/products` | 新增商品 |
| PUT | `/products/:id` | 更新商品所有欄位 |
| PATCH | `/products/:id/toggle` | 切換 is_active |
| PATCH | `/products/:id/stock` | 直接設定庫存數量 |
| DELETE | `/products/:id` | 刪除商品 |

**訂單管理**

| Method | 路徑 | 說明 |
|--------|------|------|
| GET | `/orders` | 訂單列表（?q=&status= 篩選）|
| GET | `/orders/:id` | 訂單詳情含商品明細 |
| PATCH | `/orders/:id/status` | 更新訂單狀態（含庫存退回邏輯）|
| DELETE | `/orders/:id` | 刪除訂單（含庫存退回）|

**系統設定**

| Method | 路徑 | 說明 |
|--------|------|------|
| GET | `/config` | 取得系統設定 |
| PUT | `/config` | 更新系統設定 |
| POST | `/config/test-tg` | 測試 Telegram 通知 |
| POST | `/upload` | 上傳圖片（multipart，回傳 `/uploads/xxx`）|

---

## 6. 前台功能

### 6.1 Navbar
- Logo「🥈 Silver Lord」為首頁連結（`/`）
- 語言切換（繁中 / English / 日本語）
- 購物車 icon 顯示數量 Badge

### 6.2 商品列表（`/`）
- 格狀排列 ProductCard
- 每張卡片左上角商品性質標籤（綠色「庫存品」/ 藍色「預購」）
- 庫存品庫存 = 0：顯示「庫存不足」，禁用加入購物車
- 預購：顯示預計到貨日
- 圖片右上角浮動購物車按鈕
- 圖片與名稱可點擊進入商品詳情（`/products/:id`）
- 即時報價倒數計時 + 進度條（60 秒刷新，價格變動閃爍動畫）

### 6.3 商品詳情（`/products/:id`）
- 左上角「← 返回」按鈕（`router.back()`）
- 圖片左上角商品性質標籤
- 顯示庫存數量（庫存品）/ 預計到貨日（預購）
- 即時報價（每 60 秒自動更新）
- 加入購物車（庫存 = 0 時禁用）

### 6.4 購物車（`/cart`）
- 商品名稱為超連結（→ 商品詳情）
- 數量調整（InputNumber，1~99）
- 結帳報價 10 分鐘倒數（過期後禁用結帳按鈕）
- 過期警告 Message

### 6.5 結帳（`/checkout`）
- 必填欄位：Email、收件人姓名、電話
- 出貨方式卡片選擇：郵寄（填地址）/ 7-11 店到店（填門市名稱）
- 運費自動帶入（`GET /api/v1/ship-fees`）
- 顯示商品小計 + 運費 + 訂單總計
- 10 分鐘倒數警告
- 下單成功：顯示訂單編號與鎖定總價，提示 30 分鐘內匯款，提供前往匯款回報頁按鈕
- 訂單編號存入 localStorage（`my_orders`，最多 20 筆）

### 6.6 匯款回報（`/payment`）
- 輸入訂單編號 + 匯款末五碼
- 觸發 Telegram 通知給所有管理員 Chat ID
- 訂單 status 0 → 1

### 6.7 訂單查詢（`/order-query`）
- 輸入訂單編號查詢：狀態、商品明細、運費、Email
- 可重發訂單確認信（POST `/orders/:order_no/resend-email`）

---

## 7. 後台功能

### 7.1 登入（`/admin/login`）
- 密碼登入，取得 JWT，存於 localStorage
- JWT 24 小時有效，過期自動導回登入頁

### 7.2 Dashboard（`/admin/dashboard`）
- 即時銀價（XAG/USD、XAG/TWD）
- 最後更新時間

### 7.3 商品管理（`/admin/products`）

**列表欄位：** 圖片、名稱、商品性質（Tag）、庫存數量、重量、溢價、利潤、即時價格、上架狀態

**操作：**
- 新增/編輯 Dialog：
  - 商品性質選擇（庫存品 / 預購）
  - 庫存品：庫存數量 InputNumber
  - 預購：預計到貨日 DatePicker（YYYY-MM-DD）
  - 重量、利潤、溢價、圖片上傳、上架 Checkbox
- 庫存調整 Modal（庫存品）：直接設定庫存數
- 切換上下架
- 刪除（ConfirmDialog 確認）

### 7.4 訂單管理（`/admin/orders`）

**篩選：** 訂單編號關鍵字、狀態下拉（0/1/2/3）

**列表欄位：** 訂單編號、鎖定金額、Email、匯款末五碼、狀態 Tag、建立時間

**狀態 Tag 顏色：**
- 0 待匯款：橙色
- 1 已通知：藍色
- 2 已確認：綠色
- 3 不成立：紅色

**列表樣式：**
- status=1：列底色淡黃（`#fff8e1`）
- status=3：列底色淡紅（`#fdf2f2`）、半透明

**操作按鈕：**

| 按鈕 | 顯示條件 | 動作 |
|------|---------|------|
| 確認入帳 | status=1 | status → 2 |
| 退回 | status=1 | status → 0（不退庫存）|
| 不成立 | status≠1 且 status≠3 | status → 3（退庫存）|
| 刪除 | 所有 | 刪除訂單（status≠3 則退庫存）|

**注意：** status=1（已通知/已轉帳）不可直接設為不成立，需先退回至 0。

**點擊列：** 展開詳情 Dialog（姓名、電話、Email、出貨方式/地址/門市、匯款末五碼、商品明細）

### 7.5 系統設定（`/admin/config`）
- 銀價修正（manual_offset，加在 XAG/USD 上）
- 固定匯率（manual_ex_rate，0 = 使用即時）
- Telegram Bot Token
- Telegram Chat ID 動態清單（新增/刪除，支援多組）
- 測試 Telegram 通知按鈕

### 7.6 運費設定（`/admin/shipping`）
- 郵寄運費（TWD）
- 7-11 店到店運費（TWD）

---

## 8. 庫存邏輯

| 事件 | 庫存品處理 | 預購處理 |
|------|-----------|---------|
| 結帳下單 | 扣除數量（不足拒絕下單）| 不處理 |
| 訂單 → 不成立（status 3）| 退回庫存 | 不處理 |
| 刪除訂單（status ≠ 3）| 退回庫存 | 不處理 |
| 刪除訂單（status = 3）| 不退（已退過）| 不處理 |
| 管理員庫存調整 | 直接設定新數值 | — |

---

## 9. 通知機制

### Email（SMTP）
- 觸發：結帳成功後非同步寄送訂單確認信
- 可由用戶（訂單查詢頁）或管理員手動重發

### Telegram
- 觸發：
  1. 新訂單建立（含商品明細與合計）
  2. 用戶回報匯款（含末五碼與金額）
- 支援多組 Chat ID（逗號分隔，`SendTelegramMulti` 逐一發送）
- 使用 Bot API `sendMessage`（HTML parse mode）

---

## 10. 報價機制

- Cron Job 每 60 秒從 Kitco/Yahoo Finance 抓取 XAG/USD 與 USD/TWD
- 最新報價存於記憶體 Cache（`cache.Get()`），不重複查詢 DB
- 前端 `priceStore` 每 60 秒呼叫 `/api/v1/products` 刷新報價
- 倒數計時顯示距下次更新秒數（ProgressBar）
- 報價變動時商品卡片價格有閃爍動畫

---

## 11. 建置與部署

### 開發模式
```bash
# Backend（port 8080）
cd backend && go run .

# Frontend dev server（port 5173）
cd frontend && npm run dev
```

### 生產建置（Linux amd64，執行 `backend/build-linux.bat`）
1. `npm run build`（frontend/dist/）
2. `xcopy` dist → `backend/web/`
3. `go build -tags prod -o silverlord`（嵌入 web/）

### 生產啟動
```bash
./silverlord   # 同時提供 API 與靜態前端，port 8080
```

SPA fallback 使用 `e.RouteNotFound("/*")`，不攔截 `/api/*` 路由。

---

## 12. 前端路由

| 路徑 | 頁面 | 說明 |
|------|------|------|
| `/` | ProductList | 商品列表 |
| `/products/:id` | ProductDetail | 商品詳情 |
| `/cart` | Cart | 購物車 |
| `/checkout` | Checkout | 結帳 |
| `/payment` | Payment | 匯款回報 |
| `/order-query` | OrderQuery | 訂單查詢 |
| `/admin/login` | Login | 後台登入 |
| `/admin/dashboard` | Dashboard | 後台首頁 |
| `/admin/products` | Products | 商品管理 |
| `/admin/orders` | Orders | 訂單管理 |
| `/admin/config` | Config | 系統設定 |
| `/admin/shipping` | Shipping | 運費設定 |
