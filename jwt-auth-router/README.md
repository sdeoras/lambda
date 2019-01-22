# jwt auth on GCF

## about

## usage
You will need two environment variables
```bash
export GCLOUD_PROJECT_NAME="<your GCP project name>"
export JWT_SECRET_KEY="<your secret key>"
```

This is assuming that the cloud function HTTP endpoint URI has
following format:
`"https://us-central1-" + os.Getenv("GCLOUD_PROJECT_NAME")+ ".cloudfunctions.net/router"`

## zip
`zip.sh` will zip code required to run on GCF
