# Khoan Chuyển!

A Vietnamese, family-oriented scam-interruption web app for AI Riser Vietnam 2026.

> Scammers manufacture urgency. Khoan Chuyển manufactures a pause.

## Current vertical slice

- Mobile-first Vietnamese PWA-style web UI
- Text analysis with a safe deterministic demo mode
- Gemini 3.5 Flash structured-output integration when `GEMINI_API_KEY` is set
- Explicit uncertainty and non-accusation boundaries
- Standard-library Go server, embedded static frontend, Cloud Run-ready container

## Run

```bash
go test ./...
go run .
```

Open <http://localhost:8080>. Without an API key, the app runs in clearly limited demo mode. For Gemini-backed analysis:

```bash
GEMINI_API_KEY='...' go run .
```

Never commit the key. For Cloud Run, store it in Secret Manager and expose it as the `GEMINI_API_KEY` environment variable.

## Safety boundaries

- The app does not label people, phone numbers, or accounts as criminals.
- It does not promise that a message is safe or fraudulent.
- It recommends independently finding official contact channels.
- Users are warned not to submit passwords, OTPs, or full account numbers.
- The current MVP does not store submitted messages.

## Tests

```bash
go test ./...
go vet ./...
```
