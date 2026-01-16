name: 🐰 CodeRabbit Review Issue
description: Report a CodeRabbit review feedback item
title: "[CodeRabbit] "
labels: ["code-review", "coderabbit"]
assignees: []

body:
  - type: markdown
    attributes:
      value: |
        # 🐰 CodeRabbit Review Item

        Use this template to document CodeRabbit review feedback and track follow-up actions.

  - type: dropdown
    id: severity
    attributes:
      label: Severity Level
      options:
        - "🔴 Critical - Security/Safety"
        - "🟠 High - Goroutine/Performance"
        - "🟡 Medium - Code Quality"
        - "🔵 Low - Style/Suggestion"
    validations:
      required: true

  - type: dropdown
    id: category
    attributes:
      label: Category
      options:
        - "Security"
        - "Performance"
        - "Concurrency"
        - "Error Handling"
        - "Code Quality"
        - "Best Practices"
        - "Documentation"
        - "Other"
    validations:
      required: true

  - type: textarea
    id: issue_description
    attributes:
      label: Issue Description
      description: "What does CodeRabbit's review suggest?"
      placeholder: "Describe the code issue CodeRabbit identified..."
    validations:
      required: true

  - type: textarea
    id: code_snippet
    attributes:
      label: Code Snippet
      description: "The problematic code"
      placeholder: |
        ```go
        // Your code here
        ```
      render: go
    validations:
      required: false

  - type: textarea
    id: suggested_fix
    attributes:
      label: Suggested Fix
      description: "How to fix the issue"
      placeholder: |
        ```go
        // Your fix here
        ```
      render: go
    validations:
      required: false

  - type: dropdown
    id: action_needed
    attributes:
      label: Action Needed
      options:
        - "Fix Required"
        - "Review & Discuss"
        - "Approved As-Is"
        - "Won't Fix"
        - "Deferred"
    validations:
      required: true

  - type: checkboxes
    id: task_items
    attributes:
      label: Task Items
      options:
        - label: "Code change implemented"
        - label: "Tests written/updated"
        - label: "Documentation updated"
        - label: "Peer reviewed"
        - label: "Merged to main"

  - type: textarea
    id: notes
    attributes:
      label: Additional Notes
      description: "Any additional context or discussion"
      placeholder: "Add any relevant notes..."
    validations:
      required: false

  - type: markdown
    attributes:
      value: |
        ---
        
        **Tip:** Link this issue to the PR that triggered the CodeRabbit review with `Closes #PR_NUMBER` in the PR description.
