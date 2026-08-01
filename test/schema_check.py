#!/usr/bin/env python3
"""Schema/code drift gate (make schema-check).

Validates the /calculate request/response contract in two steps:

  1. schemas/{request,response}_schema.json are each a valid JSON Schema (draft 2020-12),
  2. schemas/example_{request,response}.json each validate against their schema.

ignis stays dependency-free at runtime: request validation is hand-written in
internal/api/handler/calculation.go and surfaces.go, and this script is the
CI-time proof that those hand-written rules and the published schema describe
the same contract. A new override field without a schema entry (or vice
versa) has to be caught by eye when this script's checklist is updated, since
there is no automated Go-struct-to-schema diff — but a field present in one
and missing from the other's example will fail here.

Requires: python3 with jsonschema >= 4.18 (pip install jsonschema).
"""

import json
import pathlib
import sys

try:
    from jsonschema import Draft202012Validator
except ImportError as e:  # pragma: no cover
    sys.exit(f"schema-check needs the python 'jsonschema' package (>= 4.18): {e}")

ROOT = pathlib.Path(__file__).resolve().parent.parent
SCHEMAS = ROOT / "schemas"

PAIRS = [
    ("request", SCHEMAS / "request_schema.json", SCHEMAS / "example_request.json"),
    ("response", SCHEMAS / "response_schema.json", SCHEMAS / "example_response.json"),
]


def main() -> int:
    failed = False

    for name, schema_path, example_path in PAIRS:
        schema = json.loads(schema_path.read_text())
        Draft202012Validator.check_schema(schema)
        print(f"ok  meta-schema: {schema_path.relative_to(ROOT)}")

        example = json.loads(example_path.read_text())
        errors = sorted(
            Draft202012Validator(schema).iter_errors(example),
            key=lambda e: list(e.absolute_path),
        )
        if errors:
            failed = True
            print(f"FAIL {example_path.relative_to(ROOT)}")
            for err in errors:
                where = "/".join(map(str, err.absolute_path)) or "<root>"
                print(f"     at {where}: {err.message[:160]}")
        else:
            print(f"ok  {example_path.relative_to(ROOT)}")

    if failed:
        print("\nschema/example drift detected")
        return 1
    print(f"\nall {len(PAIRS)} example payload(s) match their schema")
    return 0


if __name__ == "__main__":
    sys.exit(main())
