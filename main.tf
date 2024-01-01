# https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_VPC.Scenarios.html
# https://medium.com/strategio/using-terraform-to-create-aws-vpc-ec2-and-rds-instances-c7f3aa416133
provider "aws" {
  region = "us-west-2" # Change this to your desired region
}

# NETWORKS

resource "aws_vpc" "house_cup_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}

# public subnet a
resource "aws_subnet" "subnet_a_public" {
  vpc_id            = aws_vpc.house_cup_vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-west-2a"

  tags = {
    Name = "House Cup Subnet A - Public"
  }
}

# private subnet a
resource "aws_subnet" "subnet_a_private" {
  vpc_id            = aws_vpc.house_cup_vpc.id
  cidr_block        = "10.0.3.0/24"
  availability_zone = "us-west-2a"

  tags = {
    Name = "House Cup Subnet A - Private"
  }
}

# private subnet b
resource "aws_subnet" "subnet_b_private" {
  vpc_id            = aws_vpc.house_cup_vpc.id
  cidr_block        = "10.0.4.0/24"
  availability_zone = "us-west-2b"

  tags = {
    Name = "House Cup Subnet B - Private"
  }
}

# db subnet group (needed for rds to enable multi az)
resource "aws_db_subnet_group" "db_subnet_group" {
  subnet_ids = [aws_subnet.subnet_a_private.id, aws_subnet.subnet_b_private.id]


  tags = {
    Name = "House Cup DB Subnet Group"
  }
}


# Internet Gateway
resource "aws_internet_gateway" "house_cup_igw" {
  vpc_id = aws_vpc.house_cup_vpc.id
}

# Public route table and association

# Route table for the public subnet
resource "aws_route_table" "public_route_table" {
  vpc_id = aws_vpc.house_cup_vpc.id

  # Route to the Internet Gateway
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.house_cup_igw.id
  }

  tags = {
    Name = "public-route-table"
  }
}

# Associate public subnet group with the public route table
resource "aws_route_table_association" "subnet_public_a_association" {
  subnet_id      = aws_subnet.subnet_a_public.id # Use one of the public subnets associated with the RDS subnet group
  route_table_id = aws_route_table.public_route_table.id
}

# Route table for private subnet
resource "aws_route_table" "private_subnet_route_table" {
  vpc_id = aws_vpc.house_cup_vpc.id

  tags = {
    Name = "Private Subnet Route Table"
  }
}

resource "aws_route_table_association" "private_subnet_association_a" {
  subnet_id      = aws_subnet.subnet_a_private.id
  route_table_id = aws_route_table.private_subnet_route_table.id
}

resource "aws_route_table_association" "private_subnet_association_b" {
  subnet_id      = aws_subnet.subnet_b_private.id
  route_table_id = aws_route_table.private_subnet_route_table.id
}

# SERVERS

# Security group for Golang server
resource "aws_security_group" "golang_sg" {
  name        = "golang_sg"
  description = "Security Group for Golang Server"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Adjust the CIDR block based on your security requirements
  }

  ingress {
    from_port   = 9090
    to_port     = 9090
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Inbound rule for HTTP
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow HTTPS (port 443) traffic
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Inbound rule for ICMP (ping)
  ingress {
    from_port   = -1
    to_port     = -1
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # You can customize egress rules as needed
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# EC2 instance for Golang server
resource "aws_instance" "golang_server" {
  ami                         = "ami-0e186acd30b9cf6a7" # Specify the AMI ID for your desired EC2 instance
  instance_type               = "t4g.small"
  key_name                    = "HouseCupKeyPair" # Replace with the name of your key pair
  associate_public_ip_address = true

  # Associate the golang_sg security group
  vpc_security_group_ids = [aws_security_group.golang_sg.id] # Ensure the instance gets a public IP

  user_data = <<-EOF
                     #!/bin/bash
                     sudo yum update -y
                     sudo yum install docker -y
                     sudo service docker start
                     sudo docker run -d -p 9090:9090 ngambino0192/house-cup-api:latest -e DB_HOST=${aws_db_instance.mysql.endpoint} ngambino0192/house-cup-api:latest
                     EOF

  tags = {
    Name = "golang-server"
  }
}

resource "aws_security_group" "bastion_sg" {
  name        = "bastion_sg"
  description = "Security Group for Bastion Server"
  vpc_id      = aws_vpc.house_cup_vpc.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Adjust the CIDR block based on your security requirements
  }

  # Inbound rule for HTTP
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# EC2 instance for bastion host
resource "aws_instance" "bastion_host" {
  ami                         = "ami-0de43e61758b7158c"
  instance_type               = "t2.micro"
  key_name                    = "HouseCupKeyPair"
  associate_public_ip_address = true

  # Specify the correct subnet ID for public subnet A or B
  subnet_id = aws_subnet.subnet_a_public.id

  # Associate a security group allowing SSH access from your local machine
  vpc_security_group_ids = [aws_security_group.bastion_sg.id]

  tags = {
    Name = "bastion-host"
  }
}


# DATABASES

resource "aws_security_group" "rds_sg" {
  vpc_id = aws_vpc.house_cup_vpc.id

  # Add any additional ingress/egress rules as needed
  ingress {
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
    security_groups = [aws_security_group.bastion_sg.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Create RDS instance in the private subnet
resource "aws_db_instance" "mysql" {
  allocated_storage = 20
  engine            = "mysql"
  engine_version    = "5.7"
  instance_class    = "db.t2.micro"
  identifier        = "house-cup-db"
  username          = "admin"
  password          = "password"

  db_subnet_group_name   = aws_db_subnet_group.db_subnet_group.id
  vpc_security_group_ids = [aws_security_group.rds_sg.id]

  skip_final_snapshot = true
  #   publicly_accessible = false

  tags = {
    Name = "house-cup-db-tag"
  }
}


# OUTPUTS

output "golang_server_public_dns" {
  value = aws_instance.golang_server.public_dns
}

output "db_instance_endpoint" {
  value = aws_db_instance.mysql.endpoint
}

output "aws_instance_bastion_public_dns" {
  value = aws_instance.bastion_host.public_ip
}

output "db_subnet_group_id" {
  value = aws_db_subnet_group.db_subnet_group.id
}

output "db_subnet_group_vpc_id" {
  value = aws_db_subnet_group.db_subnet_group.vpc_id
}

