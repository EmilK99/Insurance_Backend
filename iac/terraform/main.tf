  terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.31.0"
    }
  }

  backend "s3" {
    bucket = "flightapp-backend-terraform-state"
    key    = "terraform/state"
    region = "eu-central-1"
  }
}

  provider "aws" {
    profile = "iwan7cuw_flightapp"
    region  = "eu-central-1"
  }

  module "ecr-repository-flightapp_backend" {
    source          = "./modules/ecr-repository"
    repository_name = "flightapp_backend"
  }

  module "flightapp-backend-ec2-key-pair" {
    source          = "./modules/ec2-key-pair"
    key_pair_name = "flightapp_backend"
    ssh_public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC/W/FDHbVKXcbGXmreFWEFzJvSWPdiX7fYdYGriFzMT+Ejt5T1mbRNwoSF5FzIQf//uM1wSBWUGYaOkzrNG/5YB1aeTTS+EkJSHU/mVhON3+JnHSWG63r+2dDzZ9FFosaT9+2uHgsBv+Xm4E4FNdEyrnbHPXIL2K2O3Xm1eMeHJM3qCsKHaVMB2KWbMv4f3FSsCXKPFDPbyePyQdD6SZu2PzxQpHCSU9EY7LHq9Kh63IW+5WU60WbOdXUR+ewI4VrK8loqEhj2GTdVjkVUaQO94WgunWAQyDKTsb++3dzeXrFrebNGRzXUoj88h1nutL9hLCyVDSw0hPFWC8bONJmpIhpRPt/tUYzo9HnQhx4k3BYjjQHXPhM69kIed3kZAkCOEtkDWugUagamjv1eYEPUjlo5rwB1LSxqCoEGf02hfsR9DLZAueH5JrbJlsdx9cqO++5tq8sBTlhkpnHW0FmyDHN2lv+yUCpckyXQhphFbuc4N5jGaGKKzK4z2/OW0cE="
  }

  module "flightapp-backend-vpc" {
    source          = "./modules/vpc"
    flightapp_vpc_cidr_block = "10.0.0.0/16"
    flightapp_vpc_tag_name = "Flightapp vpc"
  }

  module "flightapp-backend-subnet" {
    source          = "./modules/vpc/modules/subnets"
    flightapp_subnet_cidr_block = "10.0.1.0/24"
    flightapp_subnet_tag_name = "Flightapp backend subnet"
    vpc_id = module.flightapp-backend-vpc.vpc_id

  }

  module "flightapp-backend-sg" {
    source          = "./modules/security-groups"
    vpc_id = module.flightapp-backend-vpc.vpc_id
    internal_cidr_blocks = [module.flightapp-backend-subnet.cidr_block]
    allowed_external_sg_addrss = ["94.180.117.169/32", "31.131.22.140/32",  "167.71.12.184/32", "176.59.148.217/32"]
    allowed_ssh_external_sg_addrss = ["94.180.117.169/32", "31.131.22.140/32", "167.71.12.184/32"]
  }

  module "flightapp-ec2-network-interface" {
    source          = "./modules/ec-2/modules/network_interface"
    subnet_id = module.flightapp-backend-subnet.subnet_id
    security_groups = [module.flightapp-backend-sg.security_group_id]
  }

  module "flightapp-ec2-isntance" {
    source          = "./modules/ec-2"
    instance_type   = "t3.micro"
    network_interface_id = module.flightapp-ec2-network-interface.interface_id
    tags_instance_name = "flightapp_backend"
    security_groups = [module.flightapp-backend-sg.security_group_id]
    #root block device
    volume_type     = "gp3"
    volume_size     = "30"
    key_name = module.flightapp-backend-ec2-key-pair.key_name
    iam_instance_profile = module.flightapp-backend-instance-profile.profile_name

  }

  module "flightapp-internet-gateway" {
    source          = "./modules/vpc/modules/internet-gateway"
    vpc_id = module.flightapp-backend-vpc.vpc_id
  }

  module "flightapp-route-table" {
    source          = "./modules/vpc/modules/routing-table"
    vpc_id = module.flightapp-backend-vpc.vpc_id
    gateway_id = module.flightapp-internet-gateway.gateway_id
  }

  module "flightapp-main-route-table-association" {
    source          = "./modules/vpc/modules/routing-table/module/route-table-association"
    vpc_id = module.flightapp-backend-vpc.vpc_id
    route_table_id = module.flightapp-route-table.routing_table_id
  }

  module "flightapp-backend-ec2-role" {
    source          = "./modules/roles"
  }

  module "flightapp-ecr-policy" {
    source          = "./modules/roles/modules/policy"
  }

  module "flightapp-ecr-policy-attachement" {
    source          = "./modules/roles/modules/role-policy-attachement"
    policy_arn      = module.flightapp-ecr-policy.policy_arn
    role_name       = module.flightapp-backend-ec2-role.role_name
  }

  module "flightapp-backend-instance-profile" {
    source          = "./modules/roles/modules/instance-profile"
    role_name       = module.flightapp-backend-ec2-role.role_name
  }

