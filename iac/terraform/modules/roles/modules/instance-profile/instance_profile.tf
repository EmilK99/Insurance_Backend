resource "aws_iam_instance_profile" "flightapp_backend_iam_instance_profile" {
  name = "flightapp_backend_iam_instance_profile"
  role = var.role_name
}
