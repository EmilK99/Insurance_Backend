resource "aws_security_group" "this" {
  name        = "fligtapp_backend"
  description = "Allow communication with backend"
  vpc_id      = var.vpc_id

  ingress {
    description      = "backend from vpc"
    from_port        = 8080
    to_port          = 8080
    protocol         = "tcp"
    cidr_blocks      = var.internal_cidr_blocks
  }

  ingress {
    description      = "backend from allowed external ip addrs"
    from_port        = 8080
    to_port          = 8080
    protocol         = "tcp"
    cidr_blocks      = var.allowed_external_sg_addrss
  }

  ingress {
    description      = "ssh from external addrs"
    from_port        = 22
    to_port          = 22
    protocol         = "tcp"
    cidr_blocks      = var.allowed_ssh_external_sg_addrss
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "Flightapp backend sg"
  }
}