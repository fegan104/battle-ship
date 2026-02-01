# Deploying Battleship Server to Google Cloud Run

This guide details the steps to deploy your Go server to Google Cloud Run.

## Prerequisites

1.  **Google Cloud Project**: You need a GCP project with billing enabled.
2.  **gcloud CLI**: Install the Google Cloud SDK: [Installation Guide](https://cloud.google.com/sdk/docs/install).

## Steps

### 1. Initialize gcloud

If you haven't already, log in and set your project:

```bash
gcloud auth login
gcloud config set project [YOUR_PROJECT_ID]
```

### 2. Enable Required Services

Enable Cloud Run and Artifact Registry (or Container Registry):

```bash
gcloud services enable run.googleapis.com artifactregistry.googleapis.com cloudbuild.googleapis.com
```

### 3. Build and Push the Image

You can use **Cloud Build** to build and push the image without needing Docker installed locally.

First, creating a repository (if you don't have one):
```bash
gcloud artifacts repositories create battleship-repo \
    --repository-format=docker \
    --location=us-central1 \
    
Then build and submit the image:
```bash
gcloud builds submit --tag us-central1-docker.pkg.dev/<PROJECT_ID>/battleship-repo/server:latest .

Deploy the image:

```bash
gcloud run deploy battleship-server \
  --image us-central1-docker.pkg.dev/<PROJECT_ID>/battleship-repo/server:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --port 8080
```

-   `--allow-unauthenticated`: Makes the server publicly accessible. Remove this if you want to restrict access.
-   `--port 8080`: Matches the port exposed in your Dockerfile and application.

### 5. Access Your Server

After deployment, the command will output a URL (e.g., `https://battleship-server-xyz.a.run.app`).

You can connect your WebSocket client to:
`wss://battleship-server-xyz.a.run.app/ws`
