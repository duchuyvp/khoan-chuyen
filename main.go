package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed static/*
var staticFiles embed.FS

type analyzeRequest struct {
	Text string `json:"text"`
}

type pauseCard struct {
	ShouldPause bool     `json:"shouldPause"`
	Risk        string   `json:"risk"`
	Summary     string   `json:"summary"`
	Signals     []string `json:"signals"`
	NextSteps   []string `json:"nextSteps"`
	Disclaimer  string   `json:"disclaimer"`
}

type analyzer interface {
	Analyze(string) (pauseCard, error)
}

type demoAnalyzer struct{}

func (demoAnalyzer) Analyze(text string) (pauseCard, error) {
	lower := strings.ToLower(text)
	signals := make([]string, 0, 3)
	for needle, signal := range map[string]string{
		"chuyển tiền": "Yêu cầu chuyển tiền",
		"ngay":        "Tạo áp lực phải hành động ngay",
		"công an":     "Tự xưng là cơ quan có thẩm quyền",
		"an toàn":     "Hứa hẹn một “tài khoản an toàn”",
	} {
		if strings.Contains(lower, needle) {
			signals = append(signals, signal)
		}
	}
	pause := len(signals) > 0
	return pauseCard{
		ShouldPause: pause,
		Risk:        map[bool]string{true: "Cần tạm dừng", false: "Chưa thấy tín hiệu rõ"}[pause],
		Summary:     "Đừng chuyển tiền chỉ dựa trên tin nhắn hoặc cuộc gọi này.",
		Signals:     signals,
		NextSteps: []string{
			"Không bấm liên kết hoặc quét mã QR trong tin nhắn.",
			"Tự tìm số điện thoại chính thức của tổ chức và gọi xác minh.",
			"Nhờ một người thân kiểm tra trước khi chuyển tiền.",
		},
		Disclaimer: "Khoan Chuyển không kết luận ai là tội phạm. Kết quả chỉ giúp bạn tạm dừng và xác minh độc lập.",
	}, nil
}

func hasConcreteRiskSignal(text string) bool {
	lower := strings.ToLower(text)
	for _, phrase := range []string{
		"chuyển tiền", "chuyển ngay", "chuyển vốn", "tài khoản an toàn",
		"quét mã qr", "mã qr", "đóng phí", "nộp phí", "otp", "mật khẩu",
		"không được kể", "giữ bí mật", "lợi nhuận", "không có rủi ro",
	} {
		if strings.Contains(lower, phrase) {
			return true
		}
	}
	return false
}

type guardedAnalyzer struct {
	next analyzer
}

func (g guardedAnalyzer) Analyze(text string) (pauseCard, error) {
	card, err := g.next.Analyze(text)
	if err != nil {
		return pauseCard{}, err
	}
	if !hasConcreteRiskSignal(text) {
		card.ShouldPause = false
		card.Risk = "Chưa thấy tín hiệu rõ"
		card.Summary = "Nội dung chưa có yêu cầu chuyển tiền, mã QR, OTP, giữ bí mật hoặc lợi nhuận bất thường. Hãy tiếp tục dùng kênh chính thức nếu cần xác minh."
		card.Signals = []string{"Chưa quan sát thấy tín hiệu giao dịch rủi ro cụ thể trong nội dung được cung cấp."}
		card.NextSteps = []string{
			"Không cần hành động vội chỉ vì một tin nhắn.",
			"Nếu nội dung liên quan tài khoản hoặc hóa đơn, tự mở ứng dụng hoặc website chính thức.",
			"Không cung cấp mật khẩu, OTP hoặc thông tin tài khoản qua tin nhắn.",
		}
		card.Disclaimer = "Khoan Chuyển không thể xác nhận một tin nhắn là an toàn. Kết quả này chỉ cho biết chưa thấy tín hiệu giao dịch rủi ro cụ thể trong nội dung đã nhập."
	}
	return card, nil
}

func newHandler(a analyzer) http.Handler {
	mux := http.NewServeMux()
	staticRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /", http.FileServer(http.FS(staticRoot)))
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("POST /api/analyze", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "application/json") {
			http.Error(w, `{"error":"Content-Type phải là application/json"}`, http.StatusUnsupportedMediaType)
			return
		}
		var input analyzeRequest
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&input); err != nil || strings.TrimSpace(input.Text) == "" {
			http.Error(w, `{"error":"Vui lòng nhập nội dung cần kiểm tra"}`, http.StatusBadRequest)
			return
		}
		card, err := a.Analyze(input.Text)
		if err != nil {
			log.Printf("analysis failed: %v", err)
			http.Error(w, `{"error":"Không thể phân tích lúc này"}`, http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(card)
	})
	return securityHeaders(mux)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; script-src 'self' 'unsafe-inline'; connect-src 'self'; img-src 'self' data:; frame-ancestors 'none'; base-uri 'none'; form-action 'self'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

func selectAnalyzer(apiKey string) analyzer {
	if project := strings.TrimSpace(os.Getenv("GOOGLE_CLOUD_PROJECT")); project != "" {
		return guardedAnalyzer{next: geminiAnalyzer{
			Project: project,
			Region:  strings.TrimSpace(os.Getenv("GOOGLE_CLOUD_REGION")),
			Model:   strings.TrimSpace(os.Getenv("GEMINI_MODEL")),
		}}
	}
	if strings.TrimSpace(apiKey) == "" {
		return demoAnalyzer{}
	}
	return guardedAnalyzer{next: geminiAnalyzer{APIKey: apiKey}}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Khoan Chuyển listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, newHandler(selectAnalyzer(os.Getenv("GEMINI_API_KEY")))))
}
