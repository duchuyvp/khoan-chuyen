# Submission Package — Khoan Chuyển

Status values: **ready**, **needs owner**, **needs evidence**, **not started**.

| Requirement | Status | Artifact / next action |
|---|---|---|
| Public Cloud Run application | ready | https://khoan-chuyen-1011704041754.asia-southeast1.run.app |
| Product built with Google AI | ready | Vertex AI Gemini 2.5 Flash, structured output, runtime identity |
| Architecture visual | ready | `docs/architecture.svg` |
| Public source repository | ready | https://github.com/duchuyvp/khoan-chuyen |
| Google AI Studio artifact | ready | Public deployment: https://khoan-chuyen.ai.studio/ — signed-out risky and benign flows verified; AI Studio project: https://aistudio.google.com/apps/0aa23848-d343-4026-ab3b-0542300e8cc7?showPreview=true&showAssistant=false |
| Real-user evidence | needs evidence | Run `docs/validation-protocol.md`; record only consented anonymous results |
| Public YouTube demo | ready | https://youtu.be/sLm0sdW-KY8 — anonymous watch + oEmbed HTTP 200 verified; title and channel matched |
| Public social journey post | ready | LinkedIn skipped at entrant request; approved Facebook text published as Nguyễn Đức Huy with audience **Công khai** and verified in the authenticated feed |
| Completion form | needs owner | Submit only after all public URLs pass anonymous checks |
| Confirmation receipt | needs owner | Save official success page/email or submission ID |

## Judge-facing pitch

**Khoan Chuyển creates a deliberate pause before a suspicious transfer.** A user pastes or describes a message; Gemini extracts observable pressure and payment signals, states uncertainty instead of accusing anyone, and returns three independent verification steps. The result can be shared with a family member without sharing the original message.

## Scorecard mapping

- **Creativity:** intervention at the decision moment—not another generic scam-information chatbot.
- **Feasibility:** one mobile-first workflow, deployed publicly on Cloud Run.
- **Impact:** testable pause conversion, comprehension, and time-to-clarity metrics.
- **Google depth:** Cloud Run + Vertex AI Gemini structured output + keyless runtime identity.
- **Safety:** no storage, no secrets, no criminal attribution, independent-channel verification.

## Demo outline (target 90–120 seconds)

1. **Problem, 0–15s:** urgency isolates victims and compresses decision time.
2. **Workflow, 15–55s:** paste synthetic prize/QR scenario; show live Gemini result.
3. **Differentiator, 55–75s:** share only the safe summary with a family member—not the raw message.
4. **Architecture, 75–95s:** Cloud Run → keyless Vertex AI Gemini; no database and no message retention.
5. **Evidence, 95–115s:** show actual consented validation metrics; do not estimate or fabricate.
6. **Close:** “Scammers manufacture urgency. Khoan Chuyển manufactures a pause.”

## Final anonymous verification

Before submitting, open every URL in a signed-out browser and verify:

- app loads and completes live analysis;
- AI Studio link is public;
- YouTube video is public;
- approved Facebook journey post is public (LinkedIn intentionally skipped);
- repository/README is public if included;
- no API key, OTP, participant identity, or private project data appears.

## Google AI Studio publication verification

- App ID: `0aa23848-d343-4026-ab3b-0542300e8cc7`
- Public URL: https://khoan-chuyen.ai.studio/
- AI Studio project URL: https://aistudio.google.com/apps/0aa23848-d343-4026-ab3b-0542300e8cc7?showPreview=true&showAssistant=false
- Published status: **Ready**
- Public HTTP check: `200`; no sign-in required
- Risky synthetic check: rendered **CẢNH BÁO CAO**, three warning signals, uncertainty language, and exactly three independent verification steps
- Benign synthetic check: rendered **LƯU Ý XÁC MINH**, explicitly noted no financial request, and retained exactly three precautionary steps
- Cost control: Gemini API monthly spend cap set to `₫100,000` on project `khoan-chuyen-airiser-2026`
- Credential boundary: AI Studio created a platform-managed API key restricted to `generativelanguage.googleapis.com`; its value was not revealed or copied and was absent from delivered HTML
