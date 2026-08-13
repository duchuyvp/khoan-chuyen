# Facebook journey post — approval required before publication

Scammers don’t just steal information—they manufacture urgency.

Đó là ý tưởng phía sau **Khoan Chuyển**, dự án mình xây dựng cho **#AIRiserVietnam #BuildwithGoogleAI**.

Người dùng dán hoặc mô tả một tin nhắn đáng ngờ bằng tiếng Việt. Gemini không kết luận ai là tội phạm và cũng không hứa rằng một tin nhắn là an toàn. Thay vào đó, ứng dụng tạo một “thẻ tạm dừng” gồm:

- các dấu hiệu gây áp lực và yêu cầu thanh toán có thể quan sát được;
- phần tóm tắt rủi ro có nêu rõ sự không chắc chắn;
- ba bước xác minh qua kênh độc lập;
- cách chia sẻ phần cảnh báo an toàn với người thân mà không chia sẻ nội dung gốc.

Ứng dụng chạy công khai trên Google Cloud Run và dùng Gemini 2.5 Flash trên Vertex AI thông qua runtime identity riêng—không có model API key trong trình duyệt hay container. Dịch vụ không có cơ sở dữ liệu ứng dụng và không lưu nội dung người dùng gửi.

Một bài học hữu ích từ quá trình kiểm thử: lần chạy synthetic regression đầu tiên tạm dừng đúng 3 tình huống rủi ro nhưng phản ứng quá mức với 2 tình huống lành tính. Sau khi thêm guardrail bảo thủ và chạy lại cùng bộ kiểm thử công khai, kết quả là 5/5. Đây là kiểm thử kỹ thuật có thể tái lập, không phải tuyên bố về tác động trên người dùng thật.

🌐 Ứng dụng: https://khoan-chuyen-1011704041754.asia-southeast1.run.app

💻 Mã nguồn: https://github.com/duchuyvp/khoan-chuyen

🎥 Demo: https://youtu.be/sLm0sdW-KY8

Scammers manufacture urgency. Khoan Chuyển manufactures a pause.

#Gemini #VertexAI #CloudRun #AI #CyberSafety #ScamPrevention
