workflow "Test package" {
  on = "push"
  resolves = ["Test"]
}

action "Lint" {
  uses = "stefanprodan/gh-actions-demo@dd27a05696158ce614f97d9e19043c29002c3b7a"
  args = "fmt"
}

action "Test" {
  needs = ["Lint"]
  uses = "stefanprodan/gh-actions-demo@dd27a05696158ce614f97d9e19043c29002c3b7a"
  args = "test"
}
