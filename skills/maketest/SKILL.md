---
name: maketest
description: >-
  Design and write software tests from requirements, agreed examples, or contracts. Use when asked to add tests, catch a bug with a regression test, pin legacy behavior before a refactor, plan test cases, or develop with TDD or BDD. Choose cases by risk and use the fastest layer that can observe the behavior. テストを追加・作成する、テストケースを設計する、回帰テストを書く、TDDで実装する。
license: MIT
compatibility: Uses the project's test runner for executable checks. Test design can be done without a runner.
metadata:
  author: otoneko.
  language: en
---

# maketest

Choose tests that would detect plausible defects in behavior the project promises. Keep the feedback loop fast and leave evidence of what the tests exercise.

## 1. Establish the oracle and scope

Read the user request, relevant requirements, project instructions, nearby tests, and the runner or CI configuration. An **oracle** is the basis for deciding whether an outcome is correct: a requirement, agreed example, published contract, domain invariant, or independently justified property. Use code to find the public seam, fixtures, and affected paths. Decide expectations from the oracle, even if that exposes a bug in the code or an existing test.

Use the table to choose a workflow and combine rows where useful. For example, TDD controls the implementation sequence while a decision table supplies the test cases. In a short plan or working note, identify the oracle, the unit of behavior, and the stopping rule. Say which risks the chosen seam cannot observe. Keep this lightweight for a single test.

| Workflow | Oracle | Test unit | Stop when |
| --- | --- | --- | --- |
| TDD for requested implementation | Requirement or agreed example | One observable behavior per cycle | The behavior list for the requested scope is covered or gaps are explicit |
| BDD / ATDD | Agreed acceptance rules and examples | One scenario with a distinct user outcome | Agreed outcomes have evidence at a suitable layer |
| Specification-based design | Specification and input or state model | One meaningful partition, boundary, rule, or transition | Selected technique items and material risks are covered or recorded as gaps |
| Contract | Published interface or actual consumer requirements | One interaction or compatibility requirement | Used interactions and relevant error variants have verification evidence |
| Domain / DDD | Business rules and invariants | One rule across a domain operation or sequence | Relevant invariants and their important violation paths are exercised |
| Regression for a known defect | Issue, agreed corrected outcome, or other independent oracle | A minimal case that distinguishes buggy and correct behavior | The defect is detected on a known bad version when available, and related cases add distinct value |
| Characterization before refactoring | Observed public behavior, explicitly labeled | One existing behavior whose change would be costly | Risky behaviors in the refactor's scope are pinned; uncertain behavior stays explicit |

Resolve contradictions that change an assertion. If intent is missing, ask a focused question and continue with the cases that have a source. Mark the unresolved case as a gap. Accept multiple outcomes only when the oracle permits them and the test can check the conditions for each. Do not silently promote a passing observation to intended behavior. Use characterization only when preserving behavior is part of the request; avoid pinning suspected defects as correct.

## 2. Choose cases and the cheapest useful layer

Start with the changed behavior and what would break for users or dependents. Map each relevant requirement to existing or proposed coverage. Consider valid and invalid inputs, meaningful boundaries, error handling, state changes, and side effects; keep a case when it guards a distinct outcome, risk, or plausible defect. For high impact behavior, explicitly account for missing coverage. Read [test design techniques](references/test-design.md) for branching rules, ordered inputs, state, or combinations. Read [workflow examples](references/workflows.md) for acceptance, contract, domain, or characterization work.

Reuse coverage that already catches the same defect. Prefer the smallest test that can see the actual failure: pure code for decision logic, integration for wiring or a storage guarantee, and end to end for critical journeys that depend on the deployed pieces working together. Combine equivalent setup with readable, named parameter rows. Each row still runs and must earn its cost. Avoid asserting framework behavior unless the project's configuration or wrapper changes its guarantees.

For a large test change or CI performance work, agree a useful runtime target and read [suite cost and reliability](references/ci.md). Use the [test plan template](assets/test-plan-template.md) when the user asks for a plan or there are several risks and layers to track.

## 3. Write an observable check and establish failure evidence

- **Test-first work:** Run the focused test before implementing the behavior. Confirm the intended test ran and failed for the predicted behavioral reason. A compile error, missing fixture, or runner failure does not establish that the assertion can detect the defect. When implementation is in scope, make the smallest coherent change, rerun to green, then refactor and rerun the affected tests.
- **Known defect:** Prefer the same test on a known bad baseline and on the candidate fix. Keep inputs and runner comparable. If a tests-only request exposes a bug, report the failing test and the requirement it contradicts; keep production fixes within the user's scope.
- **Already working behavior:** A first-run pass is normal. Check discovery and assertions, then use prior defect evidence or a focused negative control where useful. For example, confirm that the test rejects a version with the relevant limit shifted by one. If this is unsafe or impractical, report exactly what evidence is missing. See [failure checks](references/failure-checks.md) for isolation, mutation, and evidence labels.

Assert meaningful observable results: the value, error, persisted state, event, or boundary effect the contract promises. Prefer real deterministic collaborators. Use a fake or stub when the real dependency is expensive, nondeterministic, or needed to inject a specific failure. Verify calls only when the interaction itself matters, such as emitting exactly one receipt for repeated submissions. Avoid asserting private calls, using the same production calculation for the expected result, accepting unreviewed snapshots, or weakening expected outcomes to get green.

## 4. Run and report

For each new behavior assertion, name a plausible incorrect outcome and check mentally that the assertion distinguishes it. Strengthen weak assertions or remove redundant cases. Make fixtures independent, control time or randomness when relevant, wait for observable conditions with a timeout, and avoid retrying away assertion failures.

Run the focused tests, then the relevant repository checks for the change. Confirm the intended cases ran as normal tests. Check for filters, skips, expected-failure markers, or caught errors that hide the outcome. Follow existing policy for intentional skips or expected failures and report them as gaps in passing coverage. Repeat or shuffle a scoped set when there is meaningful timing, concurrency, or shared-state risk and the runner supports it. Use existing CI requirements and performance budgets. Measure when claiming a speed improvement; use comparable before and after runs.

Before handoff, compare the finished tests against the requirements in scope. For each material requirement, find the assertion that would fail if it were broken, or record the gap. A test name, covered line, or successful interaction does not establish that the outcome was checked.

Report the covered behavior and oracle, the exact checks and results, how any failure was demonstrated, and any material gaps. Distinguish observed failure from a mental check, and a drafted test from a test that ran. If tests cannot run, name the blocker and give the command to verify them.
