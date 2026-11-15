# Review checklist

Pick the viewpoints relevant to the change. Use the repository's supported inputs, callers, policies, and domain rules when deciding what constitutes a defect.

| Viewpoint | What to inspect | Evidence to seek |
| --- | --- | --- |
| Correctness and contracts | Boundary values, empty and missing fields, rounding, pagination, error paths, public APIs, and changed defaults | Assertions against the affected entry point, or a trace from a supported input through the changed path to a concrete mismatch |
| Time | Time zones, daylight-saving transitions, leap days, period boundaries, expiration, and clock source | Tests with a controlled clock and timestamps on both sides of the boundary |
| State and concurrency | Transaction boundaries, retries, idempotency, partial failures, cancellation, and competing updates | Controlled interleavings or fault injection in isolated fixtures, and invariants checked before and after each transition |
| Authorization and privacy | Ownership and tenant scoping, permission checks on changed paths, trusted identity inputs, and data in responses, logs, or error messages | Source trace of the identity and policy checks, schema or data-flow evidence, and existing defensive checks where appropriate |
| Persistence and compatibility | Existing records, mixed versions during rollout, ordering of schema and code changes, caches, and replicas | Representative old-format fixtures, backward-compatibility checks, and the documented rollout or rollback path |
| Domain and design | Required business rules, integration contracts, failure recovery, and ownership of state | Applicable requirement, reachable path, and consequence, with a focused invariant check when possible |
| Downstream effects and cost | Shared callers, configuration and packaging, work that grows with input size, and resource cleanup | Caller tests, the supported build or package check, query counts, or a bounded measurement with fixed inputs |
| Test quality | Whether the assertion reaches the changed code and would detect the effect under review | A targeted assertion against the expected result, checked on the relevant variants when feasible |

For style or standards reviews, cite the applicable trusted rule and the concrete violation. For behavior or specification reviews, cite the expected outcome and how the implementation diverges.
