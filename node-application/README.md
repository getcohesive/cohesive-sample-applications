Steps:

1. Generate SSH keys using ssh-keygen
2. Add public key to GitHub
3. Print the private key using the command `awk '{printf "%s\\n", $0}' ~/.ssh/id_rsa`
4. Copy the printed key to Cohesive as build arg within double quotes

Build:
`docker build --build-arg ssh_key="$(cat ~/.ssh/id_rsa)" --build-arg redis_host=localhost --build-arg redis_port=6379 --build-arg redis_db=1 -t nodeapp .`
