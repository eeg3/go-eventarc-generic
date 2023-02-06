# Eventarc Generic Handler

This is a generic handler for [Eventarc](https://cloud.google.com/eventarc/docs/) events passed to a [Cloud Run](https://cloud.google.com/run/docs) service. It's based on 
[GoogleCloudPlatform/golang-samples](https://github.com/GoogleCloudPlatform/golang-samples/tree/main/eventarc/generic) with a few modifications:

- Handles the request body differently.
- Swaps from Encoder to Marshal.
- [Logrus](https://github.com/sirupsen/logrus) for log output.
- Includes `cloudbuild.yaml` for build and deploy with [Cloud Build](https://cloud.google.com/build/docs).
- Adjusted testing.

## Example Usage

The following example will deploy the service with an Eventarc trigger on a GCS bucket based on new objects being created.

1. Ensure project ID is configured: `gcloud config set project <PROJECT_ID>`
2. Enable APIs:
```
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable eventarc.googleapis.com
```
3. Grant Cloud Build permission to deploy to Cloud Run:
```
export PROJECT_ID=$(gcloud config get-value project)
export PROJECT_NUMBER=$(gcloud projects list --filter="$(gcloud config get-value project)" --format="value(PROJECT_NUMBER)")
gcloud projects add-iam-policy-binding $PROJECT_ID \
        --member="serviceAccount:$PROJECT_NUMBER@cloudbuild.gserviceaccount.com" \
        --role="roles/run.admin"
gcloud projects add-iam-policy-binding $PROJECT_ID \
        --member="serviceAccount:$PROJECT_NUMBER@cloudbuild.gserviceaccount.com" \
        --role="roles/iam.serviceAccountUser"
```
4. Submit Cloud Build job to build, push, and deploy: `gcloud builds submit .`
5. Grant the default compute service account Eventarc permissions:
```
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --role roles/eventarc.eventReceiver \
  --member serviceAccount:$PROJECT_NUMBER-compute@developer.gserviceaccount.com
```
6. Create a test bucket: 
```
export BUCKET_NAME=eventarc-gcs-$PROJECT_ID
gsutil mb gs://$BUCKET_NAME
```
7. Grant the Cloud Storage service account Pub/Sub permissions:
```
SERVICE_ACCOUNT_STORAGE=$(gsutil kms serviceaccount -p $PROJECT_ID)

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member serviceAccount:$SERVICE_ACCOUNT_STORAGE \
    --role roles/pubsub.publisher
```
8. Grant the Service Account Token Creator role to the Pub/Sub service account:
```
gcloud projects add-iam-policy-binding $(gcloud config get-value project) \
    --member="serviceAccount:service-${PROJECT_NUMBER}@gcp-sa-pubsub.iam.gserviceaccount.com" \
    --role='roles/iam.serviceAccountTokenCreator'
```
9. Create the Eventarc trigger:
```
gcloud eventarc triggers create gcs-create \
--location=us \
--service-account=$PROJECT_NUMBER-compute@developer.gserviceaccount.com \
--destination-run-service=go-eventarc-generic \
--destination-run-region=us-central1 \
--destination-run-path="/" \
--event-filters="bucket=$BUCKET_NAME" \
--event-filters="type=google.cloud.storage.object.v1.finalized"
```
10. Copy a sample file into GCS:
```
touch test_file.txt
gsutil cp test_file.txt gs://$BUCKET_NAME/
```
11. Check service logs: `gcloud beta run services logs read go-eventarc-generic --region=us-central1`

## Licensing

Code in this repository is licensed under the Apache 2.0. See [LICENSE](LICENSE).
