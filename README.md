# Eltodo lunch bot

## Prerequisites
* Docker
* Set up env variables. All are required:
    - `WEBHOOK_MAIN_URL` - webhook to slack for posting lunch menues
    - `WEBHOOK_DEBUG_URL` - webhook to slack to posting errors
    - `CRON` (default: `0 10 * * 1-5`) - in default it will send every
    workday at 10:00 AM
    - `URL_BK` - url to Božská komedie menu
    - `URL_DC` - url to Di Carlo Lhotka menu
    - `URL_NK` - url to Na Kamýku menu
    - `URL_PP` - url to Pizzerie Pepino menu
    - `PATH_PDFTOTEXT` (default: `/usr/bin/pdftotext`) - path to
    `pdftotext` executable.
    - `TIMEZONE` (default: `Europe/Prague`)
    - `FOOTER` (default: `Neručíme za věrohodnost údajů. Vždy zkontrolujte oficiální nabídku.`)
    - `BOT_NAME` (default: `obědbot`)

## Getting started
1. Download or build image
2. Run it: `docker run -d --name eltodo-lunch-bot IMAGE_NAME`. And don't
forget to set up environment variables.
3. You can start requestbin from docker-compose for testing webhooks.