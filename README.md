```
git clone https://github.com/gordonIanJ/blacksite.git
cd blacksite
docker build ./ --tag=blacksite 
docker run -it --entrypoint="/bin/bash" blacksite
blacksite help
blacksite images help
export AWS_ACCESS_KEY_ID=<your AWS access key>
export AWS_SECRET_ACCESS_KEY=<your AWS secret access key>
blacksite images add <a name to give the AMI to be created>
```