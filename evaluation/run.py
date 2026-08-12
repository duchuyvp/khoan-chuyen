#!/usr/bin/env python3
import argparse
import datetime
import json
import time
import urllib.request
from pathlib import Path


def analyze(base_url, text):
    body = json.dumps({"text": text}, ensure_ascii=False).encode()
    request = urllib.request.Request(
        base_url.rstrip("/") + "/api/analyze",
        data=body,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    started = time.monotonic()
    with urllib.request.urlopen(request, timeout=90) as response:
        payload = json.load(response)
    return payload, round((time.monotonic() - started) * 1000)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--base-url", default="https://khoan-chuyen-1011704041754.asia-southeast1.run.app")
    parser.add_argument("--cases", default=str(Path(__file__).with_name("cases.json")))
    parser.add_argument("--output", default=str(Path(__file__).with_name("latest-results.json")))
    args = parser.parse_args()

    cases = json.loads(Path(args.cases).read_text())
    results = []
    for case in cases:
        payload, latency_ms = analyze(args.base_url, case["text"])
        results.append({
            "id": case["id"],
            "kind": case["kind"],
            "expected_pause": case["expected_pause"],
            "actual_pause": payload["shouldPause"],
            "correct": payload["shouldPause"] == case["expected_pause"],
            "risk": payload["risk"],
            "summary": payload["summary"],
            "signal_count": len(payload["signals"]),
            "step_count": len(payload["nextSteps"]),
            "latency_ms": latency_ms,
        })

    document = {
        "generated_at": datetime.datetime.now(datetime.timezone.utc).isoformat(),
        "base_url": args.base_url,
        "method": "Synthetic safety regression; not real-user impact evidence.",
        "passed": sum(item["correct"] for item in results),
        "total": len(results),
        "results": results,
    }
    Path(args.output).write_text(json.dumps(document, ensure_ascii=False, indent=2) + "\n")
    print(json.dumps({"passed": document["passed"], "total": document["total"], "output": args.output}))
    raise SystemExit(0 if document["passed"] == document["total"] else 1)


if __name__ == "__main__":
    main()
