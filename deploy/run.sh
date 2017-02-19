#docker run --name syntinel --detach -v /home/chris/Syntinel:/opt/Syntinel syntinel
docker run --name syntinel -it --rm -p 80:80 -v /home/chris/Syntinel:/opt/Syntinel -v /var/run/docker.sock:/var/run/docker.sock syntinel /bin/bash
