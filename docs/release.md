# Releasing a Policy

1. Update `io.kubewarden.policy.version` of the policy's `metadata.yml` version to match the version to be released.
2. Make a PR in GitHub and add the label `TRIGGER-RELEASE`.
   1. If you don't want to add the changes to the changelog, add the label `release/skip-changelog`.
3. Get an approval and merge the change.

## Mirroring the Changes to our Internal GitLab

1. Ensure that you add the GitLab repository as an origin.
2. Mirror the branches to GitLab with `git push --mirror --no-tags`.
    1. It's a good idea to have the main branch up-to-date before we add the tags. This is because we need to have an updated `.gitlab-ci.yml` to trigger release pipelines.
3. Mirror the tags to GitLab with `git push --mirror`.
