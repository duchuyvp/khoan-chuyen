package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAnalyzeCreatesPauseCard(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", strings.NewReader(`{"text":"Công an yêu cầu chuyển tiền ngay vào tài khoản an toàn để điều tra"}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	newHandler(demoAnalyzer{}).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", res.Code, res.Body.String())
	}
	var card pauseCard
	if err := json.NewDecoder(res.Body).Decode(&card); err != nil {
		t.Fatal(err)
	}
	if !card.ShouldPause {
		t.Fatal("expected transfer to be paused")
	}
	if len(card.Signals) == 0 || len(card.NextSteps) == 0 {
		t.Fatalf("incomplete card: %#v", card)
	}
}

func TestGeminiAnalyzerParsesStructuredPauseCard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "test-key" {
			t.Fatalf("missing API key")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"{\"shouldPause\":true,\"risk\":\"Cần tạm dừng\",\"summary\":\"Có dấu hiệu mạo danh.\",\"signals\":[\"Tạo áp lực\"],\"nextSteps\":[\"Gọi số chính thức\"],\"disclaimer\":\"Hãy xác minh độc lập.\"}"}]}}]}`))
	}))
	defer server.Close()

	card, err := (geminiAnalyzer{APIKey: "test-key", BaseURL: server.URL}).Analyze("chuyển tiền ngay")
	if err != nil {
		t.Fatal(err)
	}
	if !card.ShouldPause || card.Risk != "Cần tạm dừng" || len(card.Signals) != 1 {
		t.Fatalf("unexpected card: %#v", card)
	}
}

func TestHomePageExplainsPauseBeforeTransfer(t *testing.T) {
	res := httptest.NewRecorder()
	newHandler(demoAnalyzer{}).ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d", res.Code)
	}
	body := res.Body.String()
	for _, want := range []string{"Khoan Chuyển", "Tạm dừng", "/api/analyze"} {
		if !strings.Contains(body, want) {
			t.Fatalf("home page missing %q", want)
		}
	}
}

func TestAnalyzerSelectionUsesGeminiWhenKeyExists(t *testing.T) {
	t.Setenv("GOOGLE_CLOUD_PROJECT", "")
	if _, ok := selectAnalyzer("secret").(geminiAnalyzer); !ok {
		t.Fatal("expected Gemini analyzer")
	}
	if _, ok := selectAnalyzer("").(demoAnalyzer); !ok {
		t.Fatal("expected demo analyzer without key")
	}
}

func TestAnalyzerSelectionPrefersVertexRuntimeIdentity(t *testing.T) {
	t.Setenv("GOOGLE_CLOUD_PROJECT", "project-id")
	t.Setenv("GOOGLE_CLOUD_REGION", "asia-southeast1")
	a, ok := selectAnalyzer("developer-key").(geminiAnalyzer)
	if !ok || a.Project != "project-id" || a.APIKey != "" {
		t.Fatalf("expected keyless Vertex analyzer, got %#v", a)
	}
}
