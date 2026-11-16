# Workflow details

## Acceptance: BDD / ATDD

Map each story to rules, examples, and open questions. An example should give a concrete initial state, an action, and a meaningful outcome. Confirm unresolved outcomes before treating them as acceptance criteria. Use the project's scenario format if it has one. Prepare incidental state through fixtures or an API; drive the UI only where the journey or UI is what the user cares about.

## Interface contract

Separate schema conformance from consumer-driven verification. For conformance, check the real provider at the boundary against the published contract, including relevant error responses. For a consumer-driven contract, exercise the real consumer client against the contract double, then verify those interactions against the real provider in controlled provider states. A green consumer run alone leaves provider compatibility unverified.

Cover what the parties depend on, such as a mandatory field, nullable value, changed status, or message shape. Use matchers where exact values are unrelated to compatibility. Keep provider business rules in functional tests. Replace a broad scenario only if the remaining checks preserve evidence for its risks; contract tests alone may miss deployment, routing, credentials, or a full user journey.

## Domain rules and invariants

Drive a domain operation through its public interface and use real domain objects. Assert the observable outcome and preservation of the invariant, including the unchanged state or absence of an event after a rejected operation when the contract requires it. Test persistence, transactions, or concurrency at the boundary that provides the guarantee.

## Known defect and characterization

A bug report stating that processing the same order twice must create only one invoice gives a corrected outcome. Process the order twice and assert that exactly one invoice exists through the public interface or a test database. Assert an error on the second request only if the interface specifies one; returning the original invoice could also satisfy the rule.

A refactor request might instead need to preserve a legacy export whose header order is undocumented. Pin the observed public header order only if compatibility is intended. Label it as observed behavior and keep the fixture small enough to review. If the output appears wrong or includes incidental ordering, unstable fields, or secrets, resolve what must be preserved before committing a snapshot.
