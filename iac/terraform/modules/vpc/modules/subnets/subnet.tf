data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_subnet" "this" {
  vpc_id                  = var.vpc_id
  cidr_block              = var.flightapp_subnet_cidr_block
  map_public_ip_on_launch = true
  availability_zone       = data.aws_availability_zones.available.names[0]

  tags = {
    Name = var.flightapp_subnet_tag_name
  }
}

