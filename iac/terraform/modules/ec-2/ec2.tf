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
  user_data = <<-EOF
          #!/bin/bash
          sudo apt-get update
          sudo apt-get install -y docker.io docker-compose
          sudo groupadd docker
          sudo usermod -aG docker ubuntu
          EOF

  root_block_device {
    volume_type     = var.volume_type
    volume_size     = var.volume_size
  }

  network_interface {
    network_interface_id = var.network_interface_id
    device_index         = 0
  }

#  user_data = << EOF
#          #! /bin/bash
#          sudo apt-get update
#          sudo apt-get install -y docker docker-compose
#          sudo groupadd docker
#          sudo usermod -aG docker $USER
#  EOF

  tags = {
    Name = var.tags_instance_name
  }
}

