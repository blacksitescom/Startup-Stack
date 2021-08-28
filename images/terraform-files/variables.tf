variable "namespace" {
  description = "The project namespace to use for unique resource naming"
  default     = "BLACKSITE-TEST"
  type        = string
}

variable "region" {
  description = "AWS region"
  default     = "us-west-2"
  type        = string
}

variable "ami_id" {
  description = "AWS AMI ID" 
  default     = "ami-078e778114d07e9bb"
  type        = string
}