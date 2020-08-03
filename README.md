# NyanSync

[![Join the chat at https://gitter.im/state-of-the-art/NyanSync](https://badges.gitter.im/state-of-the-art/NyanSync.svg)](https://gitter.im/state-of-the-art/NyanSync?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Sharing platform to face your files in the web.

Originally based on [Syncthing](https://syncthing.net/) but could be used to share any source
through your server or redirect to original sources. Will protect your files by token, user/pass,
tls etc.

* **NyanSync** - system to sync & share your files, combining NyanShare and Syncthing
* **NyanShare** - application to share your data with anyone

*WARNING:* development still in progress.

## Usage
TODO

## TODO

* Dynamically update sources list and any data from backend
* Refresh JWT token for active user
* Login page shoud close all modal windows
* Implement application settings modal
* Remove admin password file on first login
* First user login - ask to change password

## Build NyanShare

Install `imagemagick` (convert used to generate png files)

Run `./build.sh` from repo or from clean workspace

### Steps

1. Run the go container
   ```
   host$ docker run --rm -it -v ${PWD}:/go/src/github.com/state-of-the-art/NyanSync:ro golang:1.14
   ```
2. Create tmp directory and run build inside:
   ```
   docker# mkdir tmp; cd tmp
   docker# /go/src/github.com/state-of-the-art/NyanSync/build.sh
   ```
3. Run the NyanShare (possible with providing where the gui is placed for development):
   ```
   docker# NyanShare &
   ```
4. Check the generated admin password to login:
   ```
   docker# cat /root/.config/NyanShare/admin_init_pass.txt
   ```

## Deploy NyanSync to GCP

You can relatively easy deploy NyanSync on Google Cloud Platform to get low cost and secured private
file sharing system. That will require some knowledge about how GCP is working, but overall it's not
so hard to do, following the next steps:

### Prerequesties

* Created GCP project with your full access to console
* Useful DNS name to assign to static IP address and allow GCP to create HTTPS certificate

### Steps

1. Go to `Cloud Source Repositories` and mirror the NyanSync repository
2. Go to `Cloud Build`, create the triggers based on NyanSync repository and trigger them:
    * `nyansync-gcsfuse-master`:
        * Branch: `^master$`
        * Included: `components/Dockerfile.gcsfuse`
        * Directory: `components`
        * Dockerfile: `Dockerfile.gcsfuse`
        * Image name: `gcr.io/%%PROJECT_NAME%%/nyansync-gcsfuse:latest`
    * `nyansync-encfs-master`:
        * Branch: `^master$`
        * Included: `components/Dockerfile.encfs`
        * Directory: `components`
        * Dockerfile: `Dockerfile.encfs`
        * Image name: `gcr.io/%%PROJECT_NAME%%/nyansync-encfs:latest`
    * `nyansync-syncthing-master`:
        * Branch: `^master$`
        * Included: `components/Dockerfile.syncthing`
        * Directory: `components`
        * Dockerfile: `Dockerfile.syncthing`
        * Image name: `gcr.io/%%PROJECT_NAME%%/nyansync-syncthing:latest`
3. Create GCP project service accounts:
    * `nyansync-service-account` - will be used to access buckets, don't need to be assigned to
    roles or API KEY generated 
    * `instances-controller` - for controller instance to make sure nyansync will work well
4. Assign role `Compute Instance Admin` to `instances-controller` service account
5. Create buckets to store the data and configs:
    * `%%PROJECT_NAME%%-nyansync-data` - use your project name here
        * multi-region
        * standard
        * uniform
        * Google-managed key
    * `%%PROJECT_NAME%%-nyansync-init` - use your project name here
        * multi-region
        * standard
        * uniform
        * Google-managed key
6. Assign access to the buckets:
    * `%%PROJECT_NAME%%-nyansync-data`:
        * Remove viewers from the permissions list
        * Add member: `nyansync-service-account@%%PROJECT_NAME%%.iam.gserviceaccount.com`:
        `Storage Legacy Bucket Owner`, `Storage Legacy Object Owner`
    * `%%PROJECT_NAME%%-nyansync-init`:
        * Remove viewers from the permissions list
        * Add member: `nyansync-service-account@%%PROJECT_NAME%%.iam.gserviceaccount.com`:
        `Storage Legacy Object Reader`
    * `artifacts.%%PROJECT_NAME%%.appspot.com` - stores built docker container images
        * Add member: `nyansync-service-account@%%PROJECT_NAME%%.iam.gserviceaccount.com`:
        `Storage Object Reader`
7. Generate config files for encfs:
    * Open `Cloud Shell`
    * Create empty files to store configs: `touch ~/fs.data ~/fs.conf`
    * Run docker: `docker run --rm -it -v ~/fs.data:/fs.data -v ~/fs.conf:/fs.conf alpine:3`
    * Install encfs inside: `apk add encfs`
    * Generate random password: `dd if=/dev/urandom | tr -dc _A-Z-a-z-0-9- | head -c32 > /fs.data`
    * Run encfs to generate a config: `yes | encfs -f /tmp/encfs /tmp/encfs_dec --extpass "cat /fs.data"`
    * Copy the generated config file and exit docker: `cat /tmp/encfs/.encfs6.xml > /fs.conf; exit`
    * Copy configs to the init bucket: `gsutil cp fs.* gs://%%PROJECT_NAME%%-nyansync-init/`
    * It's a good idea to duplicate the config & password in your password storage
7. Create instance template with required params:
    * Config: `N1`, `n1-standard-2` (recommended)
    * Image: `Container-Optimized OS 80` + `10GB` disk
    * Service account: `nyansync-service-account`
    * Firewall: allow HTTPS traffic
8. Create HTTPS load balancer

user-data:

## OpenSource

This is an experimental project - main goal is to test State Of The Art philosophy on practice.

We would like to see a number of independent developers working on the same project issues
for the real money (attached to the ticket) or just for fun. So let's see how this will work.

### License

Repository and it's content is covered by `Apache v2.0` - so anyone can use it without any concerns.

If you will have some time - it will be great to see your changes merged to the original repository -
but it's your choise, no pressure.
