FROM python:latest

RUN apt update \
    && apt install -y \
      nodejs \
      npm \
    && apt clean
RUN npm install -g n
RUN n stable
RUN apt purge -y nodejs npm
RUN npm install -g aws-cdk
RUN pip3 install awscli aws-cdk.core --upgrade

ENV AWS_ACCESS_KEY_ID ""
ENV AWS_SECRET_ACCESS_KEY ""
ENV AWS_DEFAULT_REGION "ap-northeast-1"
ENV BOTTOKEN ""
ENV VERIFICATIONTOKEN ""
ENV OAUTHTOKEN ""

COPY . /workdir
WORKDIR /workdir
RUN npm i
CMD ["bash"]
