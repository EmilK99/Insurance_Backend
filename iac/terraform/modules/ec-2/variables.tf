variable "security_groups" {
  type    = list
}

variable "instance_type" {
  type    = string
}

variable "network_interface_id" {
  type    = string
}

variable "tags_instance_name" {
  type    = string
}

variable "volume_type" {
  type    = string
}

variable "volume_size" {
  type    = string
}

variable "key_name" {
  type    = string
}

variable "iam_instance_profile" {
  type    = string
}