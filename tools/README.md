# tools

The [`.clabot`](.clabot) config file is the setup for the [cla-bot](https://colineberhardt.github.io/cla-bot/).
The contributors property contains a list of all the GitHub handles that are have signed the CLA, whether the user is an internal or external contributor.

Sourcegraph teammates can find the signed responses in [this spreadsheet](https://docs.google.com/spreadsheets/d/1_iBZh9PJi-05vTnlQ3GVeeRe8H3Wq1_FZ49aYrsHGLQ/edit#gid=1678726755), and edit the form [here](https://docs.google.com/forms/d/18Rd5caKejbk6ATTy-znjZofhKKeGS2AfxbsKoKFaXL8/edit).

## Syncing contributors

Credentials should be a GCP service account with at least read access to Google Forms resources.

The path to the key should be set to `GOOGLE_APPLICATION_CREDENTIALS`. For the GitHub Action, set the key to `CLABOT_CREDENTIALS` in the Secrets tab by encoding it in base64, e.g. `cat $GOOGLE_APPLICATION_CREDENTIALS | base64`.

You can manually sync contributors in an *additive* manner with the `sync` command:

```sh
go run ./cmd/sync main.go
```
