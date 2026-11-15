# Formal verification

Consider formal verification when a consequential claim depends on many states or operation orders, such as a queue delivering each item once under retries. State the invariant and what the review needs to establish before choosing tooling.

| Method | Use it when | What the result supports |
| --- | --- | --- |
| Property-based testing, such as [Hypothesis](https://hypothesis.readthedocs.io/en/latest/) | The real implementation can be exercised with generated inputs or action sequences | A failing case is evidence for the tested implementation. Passing cases cover the explored examples. |
| Explicit enumeration | The states and transitions fit in a small, well-defined finite space | Exhaustiveness applies to the enumerated space if it covers every relevant reachable state and transition within the stated bounds. |
| Model checking, such as [TLC](https://lamport.azurewebsites.net/tla/tools.html) or [Alloy Analyzer](https://alloytools.org/tutorials/online/sidenote-format-check.html) | The property is naturally expressed over state transitions or relations and the model can capture the relevant behavior | A counterexample violates the modeled property. A completed check supports the property within the model's assumptions and analyzed scope. |

1. Define the initial states, transitions, invariant, and bounds. State environmental assumptions such as delivery order, retries, fairness, or transaction guarantees. Check that relevant states are reachable and that the invariant applies there. Where practical, use a known-bad variant to show the check can detect the violation.
2. Record the tool and version, command, model or property location, scope, and observed result. Mark timeouts, sampling limits, and unexplored states. Preserve the minimized counterexample and seed when applicable.
3. Map the counterexample to reachable implementation operations. Confirm the mapping with a targeted check or a complete source trace. Report a model-only result as a design finding only when the model or specification itself is under review. Otherwise retain the implementation claim as unverified until the mapping is established.
