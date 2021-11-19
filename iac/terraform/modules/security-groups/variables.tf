variable "vpc_id" {
  type    = string
}

variable "internal_cidr_blocks" {
  type    = list
}

variable "allowed_external_sg_addrss" {
  type    = list
}

variable "allowed_ssh_external_sg_addrss" {
  type    = list
}

