# Evidence that a test can catch a defect

Choose the strongest evidence practical for the risk:

| Evidence | What to establish |
| --- | --- |
| Observed red | The intended test ran on the missing behavior or known bad version and failed at the expected assertion |
| Observed green | The intended test ran on the candidate and passed with the same expectation; report this separately from observed red |
| Isolated negative control | A plausible one-change defect violates the intended condition and the test rejects it; the unchanged candidate passes |
| Mental check | Name the defect and the assertion that should catch it; label it as unexecuted |

Prefer a known bad baseline or existing mutation tooling to editing the active worktree. If a temporary edit is necessary, isolate it in a disposable copy or worktree that includes the candidate and its tests. Do not overwrite or reset the user's changes. Keep the perturbation limited to the hypothesized defect; a syntax error, disabled test, or unconditional throw does not show that the actual behavior is checked. Confirm the final run uses the unperturbed candidate.

If the repository already has mutation tooling, try a bounded run over the relevant code when it adds evidence. Choose meaningful changes, such as inverting a guard or dropping a required field. Inspect survivors for equivalent behavior, unavailable coverage, irrelevant mutants, and actual assertion gaps. A killed mutant is evidence about that defect model. Timeouts, infrastructure errors, and skipped cases must not be claimed as observed assertion failures. Report the scope, tool result, and any limits. Avoid universal score cutoffs or adding a mutation dependency solely for routine test authoring.
