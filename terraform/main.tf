data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "web" {
  count         = 6
  ami           = data.aws_ami.ubuntu.id
  instance_type = "c5.xlarge"

  root_block_device {
    volume_type = "gp3"
    volume_size = 100
  }
  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_instance" "backend" {
  count         = 2
  ami           = "ami-05188fcabea1c2e9f"
  instance_type = "c5.large"

  root_block_device {
    volume_type = "gp2"
    volume_size = 150
  }
  ebs_block_device {
    device_name = "/dev/sdb"
    volume_size = 70
    volume_type = "gp3"
  }
  ebs_block_device {
    device_name = "/dev/sdc"
    volume_size = 140
    volume_type = "gp3"
  }
  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_ebs_volume" "main" {
  count             = 2
  availability_zone = "us-west-2a"
  size              = 40

  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_ebs_volume" "extra" {
  count             = 3
  availability_zone = "us-west-2a"
  size              = 100
  type              = "gp3"

  tags = {
    Name = "HelloWorld"
  }
}
