package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rotemtam/ent-bank-io/ent/enttest"
	"github.com/stretchr/testify/require"
)

func TestPostUser(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	srv := newServer(client)

	w := httptest.NewRecorder()
	body := `{"name": "rotemtam", "email": "rotem@entgo.io", "balance": 100}`
	req, _ := http.NewRequest("POST", "/v1/user", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer 1234")
	srv.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
}

func TestUpdateBalance(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	user := client.User.Create().
		SetName("rotem").
		SetEmail("rotem@entgo.io").
		SetBalance(100).
		SaveX(context.Background())

	srv := newServer(client)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/v1/user/%d/balance", user.ID)
	body := `{"delta": -100}`
	req, _ := http.NewRequest("PATCH", url, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer 1234")
	srv.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
	var payload struct {
		Balance float64 `json:"balance"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &payload)
	require.NoError(t, err)
	require.EqualValues(t, 0, payload.Balance)
}

// TODO: Write the test for balanceAt.
func TestBalanceAt(t *testing.T) {

}
