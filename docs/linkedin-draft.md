# LinkedIn draft — owner approval required before publication

Scammers don’t just steal information—they manufacture urgency.

That is the idea behind **Khoan Chuyển**, my project for **#AIRiserVietnam #BuildwithGoogleAI**.

A user pastes or describes a suspicious Vietnamese message. Gemini does not declare anyone a criminal or promise that a message is safe. Instead, it creates a structured “pause card” with:

- observable pressure and payment signals;
- an uncertainty-aware risk summary;
- three independent verification steps;
- a way to share the safe summary with a family member without sharing the original message.

The app runs publicly on Google Cloud Run and uses Gemini 2.5 Flash on Vertex AI through a dedicated runtime identity—no model API key in the browser or container. The service has no application database and does not persist submitted messages.

One useful lesson came from testing: the first synthetic regression correctly paused 3 risky scenarios but overreacted to 2 benign controls. I added a conservative guardrail and reran the same public test set; the result was 5/5. This is reproducible technical testing, not a claim of real-user impact. Consented user validation remains the next evidence milestone.

Live app: https://khoan-chuyen-1011704041754.asia-southeast1.run.app

Source: https://github.com/duchuyvp/khoan-chuyen

Demo video: [PUBLIC_YOUTUBE_URL]

Scammers manufacture urgency. Khoan Chuyển manufactures a pause.

#Gemini #VertexAI #CloudRun #AI #CyberSafety #ScamPrevention
