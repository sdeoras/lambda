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

## deploy
From within the folder run following command:
```bash
gcloud functions deploy router --region=us-central1 --trigger-http --runtime=go111 --set-env-vars=JWT_SECRET_KEY=${JWT_SECRET_KEY} --entry-point=Route
```
The deployment process may take several minutes. Once deployed you can list and describe the function
for more details
```bash
gcloud functions list
gcloud functions describe router
```