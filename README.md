# Traefik Basic Authentication Manager

***!!! Still work in progress !!!***

Simple API server for managing Traefik Basic Authentication based on OIDC logins.

Some systems allow only basic authentication and cannot be intergrated with modern authentication protocols. This application attempts to solve this problem by introducing temporary login credentials that expire on their own after some time and the user has to renew them again based on authentication through modern authentication protocols (OIDC).

Main features:

- Creates new Traefik basic authentication credentials. Credentials are deleted automatically after specified time. User can create new one.

- Activates Traefik basic auth logins for specified time based on authentication from OIDC (sessions). User needs to periodically renew these sessions.

## Prerequisities

This application is designed to work with:

- Traefik v2 as reverse proxy - providing TLS and basic authentication (<https://github.com/traefik/traefik>).
- Traefik middleware Traefik Forward Auth – provides OIDC authentication to Traefik (<https://github.com/thomseddon/traefik-forward-auth/>).

## How it works

Based on the OIDC authentication, the user creates new basic authentication credentials that will be valid for the period of time specified in the app configuration (`validity - credential`). These credentials are also inserted into the Traefik basic authentication configuration. They are only present in the Traefik basic authentication configuration for the time period specified in the appn configuration (`validity - session`).

If the user has valid login credentials, they can reactivate them. When the login credentials expire, the user must create new ones.

Basic authentication format:

- Usernames – are created based on OIDC login (without domain)
- Password – automatically generated and hashed (bcrypt)

## How to run using Docker

1. Run and configure [Traefik](https://github.com/traefik/traefik) and [Traefik Forward Auth](https://github.com/thomseddon/traefik-forward-auth/)

    *TODO: add examples*

2. Configure new group with guid `11111` for accessing Traefik basic auth file and DB folder.

    ```bash
    groupadd -g 11111 tbam
    ```

3. Create new folder for DB files

    ```bash
    mkdir nutsdb
    chmod 775 nutsdb
    chown root:tbam nutsdb
    ```

4. Ensure TBAM can access to Traefik basic auth file (see example file: `basic-auth_example.yml`)

    ```bash
    chmod 660 basic-auth.yml
    chown root:tbam basic-auth.yml
    ```

5. Create configuration (see: `config_example.yml`)

6. Run Docker

    ```bash
    docker run -d -v /config.yml:/tbam/config.yml -v ./nutsdb:/tbam/nutsdb -v ./basic-auth.yml:/tbam/basic-auth.yml -p 8080:8080 jirm/tbam:latest
    ```

7. Main app page is available on port `:8080`
