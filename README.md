# RaspiBot

A simple backend to control a robot using sbc motor shield from https://shop.sb-components.co.uk/.

The robot uses motor 1 to change the direction of the camera (looking up and down, no turning). Motor 2 and 3 control the wheels on one side of the bot to allow it to drive forward (both forward), backward (both reverse) or turn (one forward, the other motor on reverse). The Robot istelf is built off old Lego Mindstorms(tm) motors and lego technics, using a large Anker Powerbank which keeps the bot running (without streaming the camera to the net) for more than 96 hours. When streaming via Jitsi Meet (https://meet.jit.si), the CPU load goes up to 100% and the camera needs additional power. The motors are powered by a 9v block or (while developing) alternatively by a power supply providing 9v. As I am using a Raspberry 3 with onboard WLAN and Bluetooth, the bot  has no cables restricting its area of movement.

The main objective for creating this bot was to allow a remote player to participate in our RPG sessions and provide direct control to the camera instead of voice commands to the players moving the cam.

## License
Copyright © 2017 Ott-Consult UG (haftungsbeschränkt), Jörn Ott <raspibot@ott-consult.de>
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## API

Currently, the following API endpoints under /api/v1 are exposed:
* POST /drive/forward - makes the motors 2 and 3  drive forward
* POST /drive/reverse - makes the motors 2 and 3 drive backwards
* POST /drive/turnleft - makes motor 2 drive backwards and motor 3 drive forward
* POST /drive/turnright - makes motor 2 drive forward and motor 3 drive backwards
* POST /camera/up - makes motor 1 turn the camera upwards
* POST /camera/down - makes motor 2 turn the camera down to look at the "ground"
* POST /stop or GET /stop - stops all motors

All API endpoints except the /stop ones require an additional set of parameters. Currently, only one parameter is defined: "duration". This is the duration for the given command in seconds. If the duration is below 0, the motor is turned on forever (or if the robot falls off the table and breaks or someone calls the /stop endpoint).

## Authentication
The RaspiBot requires a http basic auth to restrict usage. You need to provide a JSON file "users.json" with the allowed users and a bcrypt hashed password.

## Security
In my setup, the RaspiBot uses a Raspberry Pi 3 with builtin WLAN and bluetooth. It uses OpenVPN to connect to my network. SSL is offloaded by an Nginx reverse proxy using the following configuration:
```
upstream robot {
    server RASPIBOT_IP:80;
    server 127.0.0.1:40080 backup;
}

server {
    listen                    443 ssl;
    server_name               RASPIBOT_FQDN;
    ssl_certificate           /etc/letsencrypt/live/RASPIBOT_FQDN/fullchain.pem;
    ssl_certificate_key       /etc/letsencrypt/live/RASPIBOT_FQDN/privkey.pem;
    ssl_client_certificate    /etc/letsencrypt/live/RASPIBOT_FQDN/fullchain.pem;
    ssl_session_timeout       5m;
    ssl_protocols             TLSv1.2;
    ssl_ciphers               EDH+aRSA+AES:EECDH+aRSA+AES:!SSLv3;
    ssl_prefer_server_ciphers on;
    ssl_session_cache         shared:SSL:10m;
    ssl_dhparam               /etc/nginx/dhparam.pem; 
    add_header                X-XSS-Protection "1; mode=block";
    add_header                Strict-Transport-Security "max-age=31536000; includeSubDomains"; 
    ssl_stapling              on;  
    ssl_stapling_verify       on;
    resolver                  8.8.8.8 8.8.4.4 valid=300s;
    resolver_timeout          5s;
    location ~ .well-known/acme-challenge/ {
        root /var/www/letsencrypt;
        default_type text/plain;
    }
    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_pass http://robot;
        client_max_body_size 1G;
    }
}

server {
    listen       80;       
    server_name  RASPIBOT_FQDN;       
    add_header   X-XSS-Protection "1; mode=block";       
    add_header   Strict-Transport-Security "max-age=31536000; includeSubDomains";         
    return       301 https://$host$request_uri;
}

server {
    listen       127.0.0.1:40080;
    server_name  RASPIBOT_FQDN;       
    root         /srv/www/RASPIBOT_FQDN;
    index        index.html;
    add_header   cache-control "max-age=0";
    add_header   cache-control "no-cache"; 
    add_header   expires       "0";
    add_header   expires       "Tue, 01 Jan 1980 1:00:00 GMT";
    add_header   pragma        "no-cache";
}

```

## Static http filesystem
To bundle the frontend with the application, I am using Statik (https://github.com/rakyll/statik). After changing any file in the static folder, you need to run statik to update the statik.go file.

## ToDo
* Use Web Sockets for a more direct response
* Create an Android App for it
* Add RPG support tools, e.g. campaigns with Ini tracker, loot list etc
* Get higher quality out of the meeting
* Document the bot (Schematics)

## Credits
In no particular order:
* SB Components (https://shop.sb-components.co.uk/) for their motor shield and the python API
* JBD for Statik (https://github.com/rakyll)
* Julien Schmidt (https://github.com/julienschmidt/) for httprouter
* Nathan Osman (https://github.com/nathan-osman/) for the GPIO library
* The Oxygen team (http://www.oxygen-icons.org/) for the Oxygen icons used