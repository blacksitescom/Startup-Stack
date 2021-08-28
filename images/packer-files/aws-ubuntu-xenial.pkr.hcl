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
  provisioner "shell" {
    inline = [
      "sudo apt-get clean",
      "sudo apt-get update",
      "sudo apt-get install -y apt-transport-https ca-certificates nfs-common apache2-utils",
      "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
      "curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -",
      "sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\"",
      "sudo apt-add-repository \"deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main\"",
      "sudo apt-get update",
      "sudo apt-get install -y docker-ce nomad consul vault ufw unzip",
      "sudo curl -L \"https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)\" -o /usr/local/bin/docker-compose",
      "sudo chmod +x /usr/local/bin/docker-compose",
      "wget https://releases.hashicorp.com/consul-template/0.25.1/consul-template_0.25.1_linux_amd64.zip",
      "sudo unzip consul-template_0.25.1_linux_amd64.zip",
      "sudo cp consul-template /usr/local/bin/consul-template",
      "sudo chmod +x /usr/local/bin/consul-template",
      "sudo mkdir -p /root/nomad/jobs",
      "sudo mkdir -p /root/docker-compose",
      "sleep 30",
      //TODO"sudo echo $(htpasswd -nb {{user `htaccess_user`}} {{user `htaccess_password`}} )"
    ]
  }
  provisioner "file" {
    source = "images/packer-files/consul/configure_consul.sh"
    //TODOdestination = "/root/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/consul/consul-server.service"
    //TODOdestination = "/etc/systemd/system/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/consul/consul-connect-enable.hcl"
    //TODOdestination = "/root/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/consul/consul-client.service"
    //TODOdestination = "/etc/systemd/system/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/nomad/nomad-server.hcl"
    //TODOdestination = "/root/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/nomad/nomad-client.hcl"
    //TODOdestination = "/root/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/nomad/configure_nomad.sh"
    //TODOdestination = "/root/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/nomad/nomad-client.service"
    //TODOdestination = "/etc/systemd/system/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/nomad/nomad-server.service"
    //TODOdestination = "/etc/systemd/system/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/nomad/jobs"
    //TODOdestination = "/root/nomad/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/vault/vault-config.hcl"
    //TODOdestination = "/root/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/vault/vault-server.service"
    //TODOdestination = "/etc/systemd/system/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/vault/enable_vault.sh"
    //TODOdestination = "/root/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/vault/init_vault.sh"
    //TODOdestination = "/root/"
    destination = "/tmp/"
  }
  provisioner "file" {
    source = "images/packer-files/docker-compose/"
    //TODOdestination = "/root/"
    destination = "/tmp/"
  }
  //TODOprovisioner "shell" {
    //TODOinline = [
      //TODO"chmod 600 /root/docker-compose/core/traefik-data/acme.json"
    //TODO]
  //TODO}
}