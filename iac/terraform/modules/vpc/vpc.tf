resource "aws_vpc" "this" {
  cidr_block       = var.flightapp_vpc_cidr_block
  instance_tenancy = "default"

  tags = {
    Name = var.flightapp_vpc_tag_name
  }
}

