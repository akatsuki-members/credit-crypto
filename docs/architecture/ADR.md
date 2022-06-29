# Architecture Decision Record

> Documenting architectural decisions helps a project succeed by helping current and future contributors understand the reasons for doing things a certain way.

## How to use?

* Please follow this template to add decision records.

```yaml
Number/Date: A unique increasing number and date that usually follows the ADR-nnnn pattern to help sort them from old to new
Title: Indicates the content
Context (Why): Describes the current situation and why you made this decision or thought it necessary—some variations explicitly break out an "alternatives covered" section to ensure all considerations get recorded
Decision (What/How): Describes the what and how of the feature
Status: Describes the status; note that ADRs can be superseded later by newer ADRs
Consequences: Describes the effect of the decision, listing positive and negative aspects
```

## Decisions

* Testing

1. [29/06/2022] We are going to use [testify](https://github.com/stretchr/testify) as the testing Library for Go.

    - Why? 
        * Sometimes `testing` native library doesn't provide out-of-the-box functionalities to make complex comparissons and mocking.
        * Comparing complex structs and slices demands lots of time.
        * Mocking structs and functionalities could be cumbersome.
    - What/How?
        * We are going to try using `testing` native library as possible and `testify` only for complex and time consuming tasks.
        * Testify supports create mock objects.
        * Testify is a mature Go library that supports complex object comparissons. 
        * The project is working at this date on a new version.
        * We are not going to use suite capabilities since debugging them is so complex.
    - Status
        * Pending for approval
    - Consequences
        * [+] Using testify in some escenarios like mocking and comparing could reduce lots of time.
        * [-] Coding assertions where you need to compare slices, maps and structs using native library could be cumbersome.
        * [-] Project will have a dependency with `testify`.

## Resources

* Why you should be using architecture decision records to document your project
    - https://www.redhat.com/architect/architecture-decision-records

* DOCUMENTING ARCHITECTURE DECISIONS
    - https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions