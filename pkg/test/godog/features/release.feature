Feature: bump release version
  In order to reflect the source code changes in the version number
  As a release manager
  I need to bump the correct semantic version level

  Scenario: Bump major level
    Given there is a directory with a git repository
    And there is a commit "go live" with the tag 2.0.0
    And there is a commit "BREAKING CHANGE: everything is new"

    When I run semver-bumper

    Then I see the version 3.0.0

  Scenario: Minor major level
    Given there is a directory with a git repository
    And there is a commit "go live" with the tag 2.0.0
    And there is a commit "feat: more functions"

    When I run semver-bumper

    Then I see the version 2.1.0

  Scenario: Patch major level
    Given there is a directory with a git repository
    And there is a commit "go live" with the tag 2.0.0
    And there is a commit "fix: bug"

    When I run semver-bumper

    Then I see the version 2.0.1
