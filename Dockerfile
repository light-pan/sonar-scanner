FROM  openjdk:8-jdk-stretch

WORKDIR  /data 
RUN mkdir -p /root/.ssh/
RUN  echo "Host gitlab.oneitfarm.com\n\tHostname 180.96.7.13\n\tPort 29622" >> /root/.ssh/config
ADD id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan -p 29622 gitlab.oneitfarm.com,180.96.7.13 >> /root/.ssh/known_hosts
ADD sonar-scanner ./sonar-scanner
ADD sonar_linux_amd64 ./sonar

CMD ./sonar
