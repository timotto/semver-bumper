Feature: add release tags to a git repository
  In order to trust the tool
  As a user
  I need to know it works

  Scenario: multiple release bumps
    Given there is a directory with a git repository
    And there is a commit "initial commit"

    When I run semver-bumper -t v -0 1.2.3

    Then I see the version 1.2.3
    And I tag the git with v1.2.3


    Given there is a commit "fix: bug"
    When I run semver-bumper -t v -0 1.2.3
    Then I see the version 1.2.4
    And I tag the git with v1.2.4


    Given there is a commit "fix: bug"
    And there is a commit "feat: feature"
    When I run semver-bumper -t v -0 1.2.3
    Then I see the version 1.3.0
    And I tag the git with v1.3.0


    Given there is a commit "feat: feature 2"
    And there is a commit "BREAKING CHANGE: all new"
    When I run semver-bumper -t v -0 1.2.3
    Then I see the version 2.0.0
    And I tag the git with v2.0.0

  Scenario: multiple release bumps with prereleases
    Given there is a directory with a git repository
    And there is a commit "initial commit"

    When I run semver-bumper -t v

    Then I see the version 1.0.0
    And I tag the git with v1.0.0


    Given there is a commit "fix: bug"
    When I run semver-bumper -t v --pre rc
    Then I see the version 1.0.1-rc.1
    And I tag the git with v1.0.1-rc.1


    Given there is a commit "fix: bug"
    When I run semver-bumper -t v --pre rc
    Then I see the version 1.0.1-rc.2
    And I tag the git with v1.0.1-rc.2


    Given there is a commit "feat: feature"
    When I run semver-bumper -t v --pre rc
    Then I see the version 1.1.0-rc.1
    And I tag the git with v1.1.0-rc.1

    When I run semver-bumper -t v
    Then I see the version 1.1.0
    And I tag the git with v1.1.0

