Feature: get help
  In order to know what I can do
  As a developer
  I need to be able to see the command help

  Scenario: Display help with long option
    When I run semver-bumper --help
    Then I see the help page
    And the exit code is 1

  Scenario: Display help with short option
    When I run semver-bumper -h
    Then I see the help page
    And the exit code is 1
