package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"silverlord/backend/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"xorm.io/xorm"
	"xorm.io/xorm/core"
)

func productDB(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	engine, err := xorm.NewEngineWithDB("postgres", "postgres://test:test@localhost/test?sslmode=disable", core.FromDB(db))
	if err != nil {
		t.Fatal(err)
	}
	previous := database.DB
	database.DB = engine
	t.Cleanup(func() {
		database.DB = previous
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
		db.Close()
	})
	return mock
}

func productRequest(t *testing.T, handler echo.HandlerFunc, method, body string, id string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, "/api/admin/products", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	if id != "" {
		c.SetParamNames("id")
		c.SetParamValues(id)
	}
	if err := handler(c); err != nil {
		t.Fatal(err)
	}
	return rec
}

func expectProductList(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectExec("LOCK TABLE product IN SHARE ROW EXCLUSIVE MODE").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT "id" FROM "product"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10).AddRow(20).AddRow(30))
}

func TestProductOrderValidation(t *testing.T) {
	for _, body := range []string{`{}`, `{"product_ids":[10,10]}`, `{"product_ids":[0]}`, `{"product_ids":[-1]}`, `{"product_ids":["10"]}`} {
		t.Run(body, func(t *testing.T) {
			productDB(t) // Invalid requests must not start a transaction or write.
			rec := productRequest(t, AdminReorderProducts, http.MethodPut, body, "")
			if rec.Code != 400 || !strings.Contains(rec.Body.String(), string(ErrProductInvalidOrder)) {
				t.Fatalf("%d %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestProductOrderRejectsStaleList(t *testing.T) {
	for _, body := range []string{`{"product_ids":[20,10]}`, `{"product_ids":[20,10,99]}`} {
		t.Run(body, func(t *testing.T) {
			mock := productDB(t)
			expectProductList(mock)
			mock.ExpectRollback()
			rec := productRequest(t, AdminReorderProducts, http.MethodPut, body, "")
			if rec.Code != 409 || !strings.Contains(rec.Body.String(), string(ErrProductOrderChanged)) {
				t.Fatalf("%d %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestProductOrderCommitsAllPositions(t *testing.T) {
	mock := productDB(t)
	expectProductList(mock)
	for i, id := range []int64{30, 10, 20} {
		mock.ExpectExec(`UPDATE "product" SET "sort_order" = \$1 WHERE "id"=\$2`).WithArgs(int64(i+1), id).WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mock.ExpectCommit()
	rec := productRequest(t, AdminReorderProducts, http.MethodPut, `{"product_ids":[30,10,20]}`, "")
	if rec.Code != 204 {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
}

func TestProductOrderRollsBackPartialFailure(t *testing.T) {
	mock := productDB(t)
	expectProductList(mock)
	mock.ExpectExec(`UPDATE "product" SET "sort_order" = \$1 WHERE "id"=\$2`).WithArgs(int64(1), int64(30)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`UPDATE "product" SET "sort_order" = \$1 WHERE "id"=\$2`).WithArgs(int64(2), int64(10)).WillReturnError(errors.New("write failed"))
	mock.ExpectRollback()
	rec := productRequest(t, AdminReorderProducts, http.MethodPut, `{"product_ids":[30,10,20]}`, "")
	if rec.Code != 500 {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
}

func TestProductListUsesStableOrder(t *testing.T) {
	mock := productDB(t)
	mock.ExpectQuery(`SELECT .* FROM "product" ORDER BY sort_order ASC, id ASC`).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "sort_order"}).AddRow(10, "First", 0).AddRow(20, "Second", 0))
	mock.ExpectQuery(`SELECT .* FROM "system_config"`).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	rec := productRequest(t, AdminGetProducts, http.MethodGet, "", "")
	if rec.Code != 200 || strings.Index(rec.Body.String(), "First") > strings.Index(rec.Body.String(), "Second") {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
}

func TestProductEditPreservesSavedOrder(t *testing.T) {
	mock := productDB(t)
	mock.ExpectQuery(`SELECT .* FROM "product"`).WithArgs(int64(20)).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "weight", "sort_order"}).AddRow(20, "Before", 1, 7))
	// Match the complete UPDATE: ordinary edits must omit sort_order even if sent by a stale form.
	mock.ExpectExec(`^UPDATE "product" SET "name" = \$1, "description" = \$2, "image_url" = \$3, "image_url2" = \$4, "image_url3" = \$5, "weight" = \$6, "labor_cost" = \$7, "premium_percent" = \$8, "is_active" = \$9, "product_type" = \$10, "stock_qty" = \$11, "expected_arrival" = \$12 WHERE "id"=\$13$`).WillReturnResult(sqlmock.NewResult(0, 1))
	rec := productRequest(t, AdminUpdateProduct, http.MethodPut, `{"name":"After","weight":1,"sort_order":99}`, "20")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"sort_order":7`) {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
}

func TestNewProductAppendsToOrder(t *testing.T) {
	mock := productDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("LOCK TABLE product IN SHARE ROW EXCLUSIVE MODE").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT "sort_order" FROM "product" ORDER BY "sort_order" DESC LIMIT 1`).WillReturnRows(sqlmock.NewRows([]string{"sort_order"}).AddRow(8))
	mock.ExpectQuery(`INSERT INTO "product" .* RETURNING "id"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(40))
	mock.ExpectCommit()
	rec := productRequest(t, AdminCreateProduct, http.MethodPost, `{"name":"New","weight":1,"sort_order":-10}`, "")
	if rec.Code != 201 || !strings.Contains(rec.Body.String(), `"sort_order":9`) {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
}
