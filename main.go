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

func newHandler(a analyzer) http.Handler {
	mux := http.NewServeMux()
	staticRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /", http.FileServer(http.FS(staticRoot)))
	mux.HandleFunc("POST /api/analyze", func(w http.ResponseWriter, r *http.Request) {
		var input analyzeRequest
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&input); err != nil || strings.TrimSpace(input.Text) == "" {
			http.Error(w, `{"error":"Vui lòng nhập nội dung cần kiểm tra"}`, http.StatusBadRequest)
			return
		}
		card, err := a.Analyze(input.Text)
		if err != nil {
			http.Error(w, `{"error":"Không thể phân tích lúc này"}`, http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(card)
	})
	return mux
}

func selectAnalyzer(apiKey string) analyzer {
	if strings.TrimSpace(apiKey) == "" {
		return demoAnalyzer{}
	}
	return geminiAnalyzer{APIKey: apiKey}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Khoan Chuyển listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, newHandler(selectAnalyzer(os.Getenv("GEMINI_API_KEY")))))
}
