# Registry-watcher

simple application that watch docker registry and notify to other service via webhook

## Feature

- get image tags per image in `config.yml`
- compare tags which is in `db` folder
- if some updates, it could notify with webhook
- polling images with cron job.