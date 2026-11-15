## Review scope

- Target and baseline — <reviewed snapshot, PR or local diff, and comparison point>
- Intent and viewpoints — <relevant requirements and areas checked>
- Status — <complete, incomplete, or no changes, with the reason if incomplete>

## Findings

1. **[<Critical, Important, or Minor>] <actionable defect>** · `<path>:<line>`
   - Trigger and impact — <reachable conditions and concrete consequence>
   - Expected behavior — <requirement, caller contract, or invariant>
   - Evidence — <execution or static trace, with the decisive result>

Use “No confirmed findings” when none meet the evidence threshold.

## Dropped candidates

<Count by reason, such as refuted, out of scope, or duplicate. Add details when they affect the review.>

## Unverified

<Material claim, why verification is blocked, and the specific check that would settle it, if known.>

## Verification summary

<For each executed check, give the command or procedure, snapshot, and relevant result. For static evidence, give the source locations and contract. Include limits that affect the conclusion.>
