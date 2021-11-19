output "cidr_block" {
  value = aws_subnet.this.cidr_block
}

output "subnet_id" {
  value = aws_subnet.this.id
}
