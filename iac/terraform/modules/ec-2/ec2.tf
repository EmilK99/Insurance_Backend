data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "this" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = var.instance_type
#  security_groups = var.security_groups
  key_name = var.key_name
# change of userdata forces replacement
  user_data = <<-EOF
          #!/bin/bash
          sudo apt-get update
          sudo apt-get install -y docker.io docker-compose awscli
          sudo grteroupadd docker
          sudo usermod -aG docker ubuntu
          runuser -l  ubuntu -c 'aws ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin 706235040724.dkr.ecr.eu-central-1.amazonaws.com/flightapp_backend'
          EOF
  iam_instance_profile = var.iam_instance_profile
  root_block_device {
    volume_type     = var.volume_type
    volume_size     = var.volume_size
  }

  network_interface {
    network_interface_id = var.network_interface_id
    device_index         = 0
  }

  provisioner "file" {
    source      = "../../test.docker-compose.yml"
    destination = "/home/ubuntu/docker-compose.yml"
  }
  provisioner "local-exec" {
    command = "runuser -l  ubuntu -c 'docker-compose up -d -f /home/ubuntu/docker-compose.yml'"
  }
#  user_data = << EOF
#          #! /bin/bash
#          sudo apt-get update
#          sudo apt-get install -y docker docker-compose awscli
#          sudo groupadd docker
#          sudo usermod -aG docker $USER
#  EOF

  tags = {
    Name = var.tags_instance_name
  }
}

