packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "ami_name" {
  type = string
  default = "null" 
}

variable "instance_type" {
  type = string
  default = "t2.micro"
}

variable "region" {
  type = string
  default = "us-west-2"
}

variable "subnet_id" {
  type = string
  default = null
}

variable "security_group_id" {
  type = string
  default = null
}

source "amazon-ebs" "ubuntu" {
  ami_name          = var.ami_name 
  instance_type     = var.instance_type 
  region            = var.region 
  subnet_id         = var.subnet_id 
  security_group_id = var.security_group_id 
  source_ami_filter {
    filters = {
      name                = "ubuntu/images/*ubuntu-xenial-16.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
  tags = {
    BuiltBy = "Blacksite"
  }
}

build {
  name    = "blacksite"
  sources = [
    "source.amazon-ebs.ubuntu"
  ]
}