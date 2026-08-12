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
	if _, ok := selectAnalyzer("secret").(guardedAnalyzer); !ok {
		t.Fatal("expected Gemini analyzer")
	}
	if _, ok := selectAnalyzer("").(demoAnalyzer); !ok {
		t.Fatal("expected demo analyzer without key")
	}
}

func TestAnalyzerSelectionPrefersVertexRuntimeIdentity(t *testing.T) {
	t.Setenv("GOOGLE_CLOUD_PROJECT", "project-id")
	t.Setenv("GOOGLE_CLOUD_REGION", "asia-southeast1")
	guard, ok := selectAnalyzer("developer-key").(guardedAnalyzer)
	a, innerOK := guard.next.(geminiAnalyzer)
	if !ok || !innerOK || a.Project != "project-id" || a.APIKey != "" {
		t.Fatalf("expected keyless Vertex analyzer, got %#v", a)
	}
}

func TestReadyEndpointReportsReadyWithoutCallingAnalyzer(t *testing.T) {
	res := httptest.NewRecorder()
	newHandler(demoAnalyzer{}).ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/readyz", nil))

	if res.Code != http.StatusOK || strings.TrimSpace(res.Body.String()) != "ok" {
		t.Fatalf("status = %d, body = %q", res.Code, res.Body.String())
	}
}

func TestResponsesIncludeBrowserSecurityHeaders(t *testing.T) {
	res := httptest.NewRecorder()
	newHandler(demoAnalyzer{}).ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))

	for header, want := range map[string]string{
		"Content-Security-Policy": "default-src 'self'",
		"Referrer-Policy":         "no-referrer",
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "DENY",
	} {
		if got := res.Header().Get(header); !strings.Contains(got, want) {
			t.Errorf("%s = %q, want it to contain %q", header, got, want)
		}
	}
}

func TestAnalyzeRejectsNonJSONRequests(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", strings.NewReader(`{"text":"test"}`))
	req.Header.Set("Content-Type", "text/plain")
	res := httptest.NewRecorder()
	newHandler(demoAnalyzer{}).ServeHTTP(res, req)

	if res.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnsupportedMediaType)
	}
}

func TestHomePageOffersPrivateResultSharing(t *testing.T) {
	res := httptest.NewRecorder()
	newHandler(demoAnalyzer{}).ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	body := res.Body.String()
	for _, want := range []string{"id=\"share\"", "navigator.share", "Không gửi nội dung gốc"} {
		if !strings.Contains(body, want) {
			t.Fatalf("home page missing %q", want)
		}
	}
}

func TestClearEverydayMessagesDoNotTriggerPause(t *testing.T) {
	cases := []string{
		"Mẹ nhắc con chiều nay mua giúp rau và gọi lại khi tan làm nhé.",
		"Hóa đơn tháng này sắp đến hạn. Bạn hãy tự mở ứng dụng của nhà cung cấp hoặc gọi số trên hợp đồng để kiểm tra.",
	}
	for _, text := range cases {
		if hasConcreteRiskSignal(text) {
			t.Errorf("unexpected risk signal in %q", text)
		}
	}
}

func TestConcreteTransferRiskSignalsAreDetected(t *testing.T) {
	cases := []string{
		"Quét mã QR và đóng phí trong 10 phút.",
		"Chuyển ngay 45 triệu vào tài khoản an toàn và không được kể cho ai.",
		"Cam kết lợi nhuận 30% mỗi tháng, không có rủi ro.",
		"Đọc mã OTP để xác minh tài khoản.",
	}
	for _, text := range cases {
		if !hasConcreteRiskSignal(text) {
			t.Errorf("missing risk signal in %q", text)
		}
	}
}
