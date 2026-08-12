# Submission Package — Khoan Chuyển

Status values: **ready**, **needs owner**, **needs evidence**, **not started**.

| Requirement | Status | Artifact / next action |
|---|---|---|
| Public Cloud Run application | ready | https://khoan-chuyen-1011704041754.asia-southeast1.run.app |
| Product built with Google AI | ready | Vertex AI Gemini 2.5 Flash, structured output, runtime identity |
| Architecture visual | ready | `docs/architecture.svg` |
| Public source repository | needs owner | Create/select repository, then push reviewed commits |
| Google AI Studio share link | not started | Create competition-required AI Studio artifact and make it public |
| Real-user evidence | needs evidence | Run `docs/validation-protocol.md`; record only consented anonymous results |
| Public YouTube demo | not started | Record after evidence and final UI freeze |
| LinkedIn journey post | needs owner | Draft after demo URL exists; owner publishes from own account |
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
- LinkedIn post is public;
- repository/README is public if included;
- no API key, OTP, participant identity, or private project data appears.
