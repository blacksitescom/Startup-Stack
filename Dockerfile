#FROM golang:1.16-alpine
FROM arm64v8/golang:1.16

ENV PKR_VAR_ami_name=test
ENV PKR_VAR_instance_type=t2.micro
ENV PKR_VAR_region=us-west-2
ENV PKR_VAR_subnet_id=subnet-002fa0d8863ebd706
ENV PKR_VAR_security_group_id=sg-0c5426731ae8feadc

ENV PACKER_VERSION=1.7.4
ENV PACKER_SHA256SUM=3660064a56a174a6da5c37ee6b36107098c6b37e35cc84feb2f7f7519081b1b0

ENV TERRAFORM_VERSION=1.0.4
ENV TERRAFORM_SHA256SUM=5c0be4d52de72143e2cd78e417ee2dd582ce229d73784fd19444445fa6e1335e

#RUN apk update && apk add bash git wget openssl
RUN apt-get update -y && apt-get install -y bash git wget openssl unzip

#ADD https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_linux_amd64.zip ./
ADD https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_linux_arm64.zip ./
ADD https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_SHA256SUMS ./

#ADD https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip ./
ADD https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_arm64.zip ./
ADD https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_SHA256SUMS ./

#RUN sed -i '/.*linux_amd64.zip/!d' packer_${PACKER_VERSION}_SHA256SUMS
RUN sed -i '/.*linux_arm64.zip/!d' packer_${PACKER_VERSION}_SHA256SUMS
#RUN sha256sum -cs packer_${PACKER_VERSION}_SHA256SUMS
#RUN unzip packer_${PACKER_VERSION}_linux_amd64.zip -d /bin
RUN unzip packer_${PACKER_VERSION}_linux_arm64.zip -d /bin
#RUN rm -f packer_${PACKER_VERSION}_linux_amd64.zip
RUN rm -f packer_${PACKER_VERSION}_linux_arm64.zip

#RUN sed -i '/.*linux_amd64.zip/!d' terraform_${TERRAFORM_VERSION}_SHA256SUMS
RUN sed -i '/.*linux_arm64.zip/!d' terraform_${TERRAFORM_VERSION}_SHA256SUMS
#RUN sha256sum -cs terraform_${TERRAFORM_VERSION}_SHA256SUMS
#RUN unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /bin
RUN unzip terraform_${TERRAFORM_VERSION}_linux_arm64.zip -d /bin
#RUN rm -f terraform_${TERRAFORM_VERSION}_linux_amd64.zip
RUN rm -f terraform_${TERRAFORM_VERSION}_linux_arm64.zip

WORKDIR /app
COPY main.go ./
COPY go.mod  ./
COPY go.sum ./
COPY cmd ./cmd
COPY execHelpers ./execHelpers
COPY images ./images
WORKDIR /app/execHelpers
RUN go build
WORKDIR /app/images
RUN go build
WORKDIR /app
RUN go install github.com/gordonianj/blacksite
ENV PACKER_PLUGIN_PATH=/app/images/packer-files/plugins 
ENV PACKER_LOG=1
