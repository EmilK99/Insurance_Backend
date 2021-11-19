resource "aws_route_table" "this" {
  vpc_id = var.vpc_id


  #  route {
  #    cidr_block = var.cidr_block
  #    gateway_id = "local"
  #  }

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = var.gateway_id


  }

}