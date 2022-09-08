# tools

The [`.clabot`](../.clabot) config file is the setup for the [cla-bot](https://colineberhardt.github.io/cla-bot/).
The contributors property contains a list of all the GitHub handles that are have signed the CLA, whether the user is an internal or external contributor.

Sourcegraph teammates can find the signed responses in [this spreadsheet](https://docs.google.com/spreadsheets/d/1_iBZh9PJi-05vTnlQ3GVeeRe8H3Wq1_FZ49aYrsHGLQ/edit#gid=1678726755), and edit the form [here](https://docs.google.com/forms/d/18Rd5caKejbk6ATTy-znjZofhKKeGS2AfxbsKoKFaXL8/edit).

## Syncing contributors

Credentials should be a GCP service account with at least read access to Google Forms resources.

The path to the key should be set to `GOOGLE_APPLICATION_CREDENTIALS`. For the GitHub Action, set the key to `CLABOT_CREDENTIALS` in the Secrets tab by encoding it in base64, e.g. `cat $GOOGLE_APPLICATION_CREDENTIALS | base64`.
The service account currently in use is [clabot](https://console.cloud.google.com/iam-admin/serviceaccounts/details/clabot@sourcegraph-ci.iam.gserviceaccount.com?project=sourcegraph-ci) - it must have the following access to the spreadsheets noted above:

- `https://www.googleapis.com/auth/forms.responses.readonly`
- `https://www.googleapis.com/auth/forms.body.readonly`

If access cannot be granted directly to the service account, impersonation must be used.
To set it up:

- Grant domain-wide delegation to the service account
- Set `GOOGLE_TARGET_SERVICE_ACCOUNT` to the address of the service account, e.g. `clabot@sourcegraph-ci.iam.gserviceaccount.com`
- Set `GOOGLE_IMPERSONATE_USER` to the address of the target account to impersonate, e.g. `robert@sourcegraph.com`

You can then manually sync contributors in an *additive* manner with the `sync` command:

```sh
go run ./tools/sync main.go
```

You can also manually run the sync workflow from the ["Actions" tab of the repository](https://github.com/sourcegraph/clabot-config/actions).

### Error reporting

Error reporting to Sentry can be enabled by setting `SENTRY_DSN`.
Currently this should be configured to go to the `clabot` project: [`sourcegraph/clabot`](https://sentry.io/organizations/sourcegraph/issues/?project=6732169&query=&statsPeriod=14d).

The logging is configured such that the `InstanceID` is `GITHUB_RUN_ID`.
