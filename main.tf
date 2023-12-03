
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

variable "aws_region" {
  type    = string
  default = "sa-east-1" # São Paulo
}

variable "ubuntu-ami" {
  type    = string
  default = "ami-0b6c2d49148000cd5" # Ubuntu 22.04
}

provider "aws" {
  region = var.aws_region
}

resource "aws_vpc" "load_test_vpc" {
  cidr_block = "10.0.0.0/16"
  enable_dns_hostnames = true
  tags = {
    Name = "load-test-vpc"
  }
}

resource "aws_subnet" "load_test_subnet" {
  vpc_id     = aws_vpc.load_test_vpc.id
  cidr_block = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone = "sa-east-1a"
  tags = {
    Name = "load-test-subnet"
  }
}

# Additional subnet (in a different AZ, e.g., sa-east-1b)
resource "aws_subnet" "load_test_subnet_2" {
  vpc_id     = aws_vpc.load_test_vpc.id
  cidr_block = "10.0.2.0/24"  # Ensure a different CIDR block than the first subnet
  map_public_ip_on_launch = true
  availability_zone = "sa-east-1b"  # Different AZ
  tags = {
    Name = "load-test-subnet-2"
  }
}

resource "aws_internet_gateway" "load_test_gateway" {
  vpc_id = aws_vpc.load_test_vpc.id
  tags = {
    Name = "load-test-gateway"
  }
}

resource "aws_route_table" "load_test_route_table" {
  vpc_id = aws_vpc.load_test_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.load_test_gateway.id
  }
  tags = {
    Name = "load-test-route-table"
  }
}

resource "aws_route_table_association" "load_test_route_table_assoc" {
  subnet_id      = aws_subnet.load_test_subnet.id
  route_table_id = aws_route_table.load_test_route_table.id
}

resource "aws_security_group" "load_test_security_group" {
  name        = "load-test-security-group"
  description = "Allow HTTP and SSH for load testing"
  vpc_id      = aws_vpc.load_test_vpc.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "load-test-security-group"
  }
}

resource "aws_instance" "load_test_api" {
  ami           = var.ubuntu-ami
  instance_type = "t2.micro"
  vpc_security_group_ids = [aws_security_group.load_test_security_group.id]
  subnet_id     = aws_subnet.load_test_subnet.id
  key_name      = "ec2-node-key"

  tags = {
    Name = "load-test-api"
  }

  user_data = file("./ec2-startup.sh")
}

resource "aws_db_instance" "load_test_postgres" {
  allocated_storage    = 20
  storage_type         = "gp2"
  engine               = "postgres"
  engine_version       = "15.4"
  instance_class       = "db.t3.micro"
  db_name              = "load_test_db" 
  username             = "tester"
  password             = "StressTestPassword"
  parameter_group_name = "default.postgres15"
  db_subnet_group_name = aws_db_subnet_group.load_test_subnet_group.name
  skip_final_snapshot  = true
  publicly_accessible  = true

  vpc_security_group_ids = [aws_security_group.load_test_db_security_group.id]

  tags = {
    Name = "load-test-postgres"
  }
}

resource "aws_db_subnet_group" "load_test_subnet_group" {
  name       = "load-test-subnet-group"
  subnet_ids = [aws_subnet.load_test_subnet.id, aws_subnet.load_test_subnet_2.id]

  tags = {
    Name = "load-test-db-subnet-group"
  }
}

resource "aws_security_group" "load_test_db_security_group" {
  name        = "load-test-db-security-group"
  description = "Allow PostgreSQL traffic"
  vpc_id      = aws_vpc.load_test_vpc.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] 
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "load-test-db-security-group"
  }
}

resource "null_resource" "db_setup" {
  depends_on = [aws_db_instance.load_test_postgres, aws_security_group.load_test_db_security_group]

  provisioner "local-exec" {
    command = "sleep 60 && psql -h ${aws_db_instance.load_test_postgres.address} -U ${aws_db_instance.load_test_postgres.username} -d ${aws_db_instance.load_test_postgres.db_name} -f ./create-db.sql"
    environment = {
      PGPASSWORD = aws_db_instance.load_test_postgres.password
    }
  }
}

output "ec2_public_ip" {
  value = aws_instance.load_test_api.public_ip
  description = "The public IP address of the EC2 server instance."
}

