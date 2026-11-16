# Selecting test cases

Derive cases before enumerating implementation branches. Use the techniques that fit the model and risk; none is mandatory merely because a function has inputs. Record the items selected and the important items deferred. Coverage of selected items does not imply full coverage of a technique.

| Technique | Use when | What to select |
| --- | --- | --- |
| Equivalence partitioning | Inputs, outputs, or conditions fall into groups with the same specified behavior | Choose a representative from each relevant group, including invalid groups; use separate cases for invalid inputs when one error could mask another |
| Boundary value analysis | An ordered domain has thresholds or adjacent partitions | Check the boundary and its neighbor across the partition; for risk such as an off-by-one defect, also check the neighbor on the other side using the domain's actual step, precision, or calendar |
| Decision table | Several conditions jointly change an action, error, or precedence | Map conditions to outcomes, discard impossible combinations, and merge irrelevant conditions only when the rule guarantees they do not affect the outcome; cover selected feasible rules and make any risk-based omissions explicit |
| State transitions | Results depend on the prior state or on an event sequence | Select relevant starting states, events, guards, and observable new states; check invalid transitions and forbidden effects when specified, and multi-step sequences when single transitions miss a risk |
| Pairwise or higher-strength combinations | Several factors interact and the feasible cross-product is too expensive | Model valid combinations and constraints; use pairwise when two-factor interactions are a reasonable risk model, and explicitly add high-risk higher-order combinations or increase strength; derive expected results separately |
| Property-based testing | An independent invariant holds over a broad class of inputs | Generate meaningful valid and invalid data, include boundary examples, retain shrinkable and reproducible failures, and keep explicit examples for defects the property would not catch |

Example: A discount applies to a cart of at least 10 items only for members, unless a manual price override is set. Test the inclusive threshold with member carts of 9 and 10. Test a nonmember with 10 and a member with 10 plus an override. If the override can also affect nonmembers or smaller carts differently under the specification, include that rule. Larger carts and different item IDs need no separate tests unless they guard another risk.

For a property, explain why it follows from the contract and which wrong implementations could satisfy it. For example, round-tripping a serializer and deserializer can miss two components that agree on the wrong wire format. Add an independently established wire example when interoperability matters.

Check the consequences of a rejected operation when the specification defines them. For example, if a rejected upload must leave no stored file, assert the rejection and verify that storage is unchanged. A test that checks only the error can pass even if the upload wrote data before returning it.
